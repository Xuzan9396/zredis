package sredis

import (
	"crypto/tls"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

type RedisPool struct {
	redis_pool  *redis.Pool
	maxActive   int
	maxIdle     int
	idleTime    time.Duration
	redisOption []redis.DialOption
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
	optionDefalt := []redis.DialOption{
		redis.DialConnectTimeout(time.Duration(5) * time.Second),
		redis.DialReadTimeout(time.Duration(10) * time.Second),
		redis.DialWriteTimeout(time.Duration(10) * time.Second),
	}
	if len(redisPool.redisOption) > 0 {
		optionDefalt = append(optionDefalt, redisPool.redisOption...)
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
				optionDefalt...,
			)
			if err != nil {
				//zlog.F().Error("Redis Redis 连接错误", err)
				log.Println(err)
				return nil, err
			}
			//验证redis 是否有密码
			if auth != "" {
				if _, err := c.Do("AUTH", auth); err != nil {

					//zlog.F().Error("Connect to redis AUTH error", err)
					log.Println("Connect to redis AUTH error", err)
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

	c := pool.Get()
	if c.Err() != nil {
		//zlog.F().Fatalf("conn:%s,err:%v", conn, c.Err())
		//log.Fatalf("conn:%s,err:%v", conn, c.Err())
		log.Println("conn:%s,err:%v", conn, c.Err())
		pool.Close() // 关闭连接池避免资源泄露
		return nil
	}
	c.Close()

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

// WithRedisOption 设置redis连接选项
func WithRedisOption(opts ...redis.DialOption) Redis_func {
	return func(r *RedisPool) {
		r.redisOption = opts
	}
}

// WithRedisTLS 设置redis连接是否使用TLS
// 注意: 生产环境中应该提供正确的TLS配置，而不是跳过验证
func WithRedisTLS() Redis_func {
	return func(r *RedisPool) {
		r.redisOption = append(r.redisOption, redis.DialUseTLS(true), redis.DialTLSConfig(&tls.Config{
			MinVersion: tls.VersionTLS11, // 使用TLS 1.2或更高版本
			// InsecureSkipVerify: true, // 生产环境中应该移除此选项
		}))
	}
}

// WithRedisTLSConfig 允许用户提供自定义的TLS配置
func WithRedisTLSConfig(tlsConfig *tls.Config) Redis_func {
	return func(r *RedisPool) {
		if tlsConfig == nil {
			// 提供默认的安全TLS配置
			tlsConfig = &tls.Config{
				MinVersion: tls.VersionTLS11,
			}
		}
		r.redisOption = append(r.redisOption, redis.DialUseTLS(true), redis.DialTLSConfig(tlsConfig))
	}
}

func (this *RedisPool) CommonCmd(cmdStr string, keysAndArgs ...interface{}) (reply interface{}, err error) {
	if this == nil {
		//zlog.F().Errorf("RedisPool 实例为空")
		return nil, fmt.Errorf("RedisPool instance is nil")
	}

	if this.redis_pool == nil {
		//zlog.F().Errorf("Redis 连接池为空")
		return nil, fmt.Errorf("redis pool is nil")
	}

	c := this.redis_pool.Get()
	if c == nil {
		//zlog.F().Errorf("获取 Redis 连接失败: 连接为空")
		return nil, fmt.Errorf("redis connection is nil")
	}
	if c.Err() != nil {
		//zlog.F().Errorf("Redis DoCommonCmd: %v", c.Err())
		return nil, c.Err()
	}
	defer c.Close()
	res, err := c.Do(cmdStr, keysAndArgs...)
	if err == nil {
		return res, err
	}
	return nil, err
}

func (this *RedisPool) CommonLuaScript(script string, key string, args ...interface{}) (reply interface{}, err error) {
	if this == nil {
		//zlog.F().Errorf("RedisPool 实例为空")
		return nil, fmt.Errorf("RedisPool instance is nil")
	}

	if this.redis_pool == nil {
		//zlog.F().Errorf("Redis 连接池为空")
		return nil, fmt.Errorf("redis pool is nil")
	}

	c := this.redis_pool.Get()
	if c == nil {
		//zlog.F().Errorf("获取 Redis 连接失败: 连接为空")
		return nil, fmt.Errorf("redis connection is nil")
	}
	if c.Err() != nil {
		//zlog.F().Errorf("LuaCommonCmd get redis error: %v", c.Err())
		return nil, c.Err()
	}
	defer c.Close()

	lua := redis.NewScript(1, script)
	return lua.Do(c, append([]interface{}{key}, args...)...)
}
