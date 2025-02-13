package sredis

import (
	"github.com/Xuzan9396/zlog"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisPool struct {
	redis_pool *redis.Pool

	maxActive int
	maxIdle   int
	idleTime  time.Duration
}
type Redis_func func(*RedisPool)

func Conn(conn, auth string, dbnum int, opts ...Redis_func) *RedisPool {

	redisPool := &RedisPool{
		maxActive: 100,
		maxIdle:   50,
		idleTime:  300 * time.Second,
	}

	for _, opt := range opts {
		opt(redisPool)
	}

	pool := &redis.Pool{
		MaxActive:   redisPool.maxActive, // 最大活跃 假设应用在高并发场景下，最大并发请求数为 1000，那么可以将 MaxActive 设置为 2000-3000。
		MaxIdle:     redisPool.maxIdle,   // 最大空闲 ,一般启动时候保持的链接数
		IdleTimeout: redisPool.idleTime,  // 空闲连接超时 超过这个时间会关闭空闲链接
		Wait:        false,               // true: 如果达到最大连接数，等待空闲连接 false: 直接返回错误
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				conn,
				redis.DialConnectTimeout(time.Duration(5)*time.Second),
				redis.DialReadTimeout(time.Duration(10)*time.Second),
				redis.DialWriteTimeout(time.Duration(10)*time.Second),
			)
			if err != nil {
				zlog.F().Error("Redis Redis 连接错误", err)
				return nil, err
			}
			//验证redis 是否有密码
			if auth != "" {
				if _, err := c.Do("AUTH", auth); err != nil {

					zlog.F().Error("Connect to redis AUTH error", err)
					c.Close()
					return nil, err
				}
			}
			c.Do("select", dbnum)

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	redisPool.redis_pool = pool

	return redisPool
}

// WithMaxActive 设置最大活跃
func WithMaxActive(maxActive int) Redis_func {
	return func(r *RedisPool) {
		r.maxActive = maxActive
	}
}

// WithMaxIdle 设置最大空闲
func WithMaxIdle(maxIdle int) Redis_func {
	return func(r *RedisPool) {
		r.maxIdle = maxIdle
	}
}

// WithIdleTime 设置空闲连接超时时间
func WithIdleTime(idleTime time.Duration) Redis_func {
	return func(r *RedisPool) {
		r.idleTime = idleTime
	}
}

func (this *RedisPool) CommonCmd(cmdStr string, keysAndArgs ...interface{}) (reply interface{}, err error) {
	c := this.redis_pool.Get()
	if c.Err() != nil {
		zlog.F().Errorf("Redis DoCommonCmd: %v", c.Err())
		return
	}
	defer c.Close()
	res, err := c.Do(cmdStr, keysAndArgs...)
	if err == nil {
		return res, err
	}
	return nil, err
}

func (this *RedisPool) CommonLuaScript(script string, key string, args ...interface{}) (reply interface{}, err error) {
	c := this.redis_pool.Get()
	if c.Err() != nil {
		zlog.F().Errorf("LuaCommonCmd get redis error: %v", c.Err())
		return
	}
	defer c.Close()

	lua := redis.NewScript(1, script)
	return lua.Do(c, append([]interface{}{key}, args...)...)
}
