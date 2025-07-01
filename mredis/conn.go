package mredis

import (
	"fmt"
	"github.com/Xuzan9396/zlog"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

type RedisPool struct {
	redis_pool *redis.Pool
	maxActive  int
	maxIdle    int
	idleTime   time.Duration
}

type Redis_func func(*RedisPool)

type RedisPoolManager struct {
	// 管理多个 Redis 连接池
	pools map[string]*RedisPool
	mu    sync.RWMutex // 读写锁，保证并发安全
}

var redisManager *RedisPoolManager

func init() {
	redisManager = &RedisPoolManager{
		pools: make(map[string]*RedisPool),
	}
}

// Conn 创建或获取指定名称的 Redis 连接池
func Conn(name, conn, auth string, dbnum int, opts ...Redis_func) {
	// 如果已经存在该连接池，直接返回
	redisManager.mu.RLock()
	if _, exists := redisManager.pools[name]; exists {
		redisManager.mu.RUnlock()
		return
	}
	redisManager.mu.RUnlock()

	// 创建新的连接池
	redisPool := &RedisPool{
		maxActive: 100,
		maxIdle:   50,
		idleTime:  300 * time.Second,
	}

	for _, opt := range opts {
		opt(redisPool)
	}

	pool := &redis.Pool{
		MaxActive:   redisPool.maxActive, // 最大活跃
		MaxIdle:     redisPool.maxIdle,   // 最大空闲
		IdleTimeout: redisPool.idleTime,  // 空闲连接超时
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
				zlog.F().Error("Redis 连接错误", err)
				return nil, err
			}
			//验证redis 是否有密码
			if auth != "" {
				if _, err := c.Do("AUTH", auth); err != nil {
					c.Close()
					zlog.F().Fatalf("Connect to redis AUTH error: %v", err)
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

	// 存储连接池
	redisManager.mu.Lock()
	redisManager.pools[name] = &RedisPool{
		redis_pool: pool,
		maxActive:  redisPool.maxActive,
		maxIdle:    redisPool.maxIdle,
		idleTime:   redisPool.idleTime,
	}
	redisManager.mu.Unlock()
}

// 获取指定名称的 Redis 连接池
func getPool(name string) (*redis.Pool, error) {
	redisManager.mu.RLock()
	defer redisManager.mu.RUnlock()

	if pool, exists := redisManager.pools[name]; exists {
		return pool.redis_pool, nil
	}
	return nil, fmt.Errorf("RedisPool not found: %s", name)
}

// CommonCmd 执行通用的 Redis 命令
func CommonCmd(name, cmdStr string, keysAndArgs ...interface{}) (reply interface{}, err error) {
	pool, err := getPool(name)
	if err != nil {
		zlog.F().Errorf("获取 Redis 连接池失败: %v", err)
		return nil, err
	}
	
	if pool == nil {
		zlog.F().Errorf("Redis 连接池为空: %s", name)
		return nil, fmt.Errorf("redis pool is nil for name: %s", name)
	}

	c := pool.Get()
	if c == nil {
		zlog.F().Errorf("获取 Redis 连接失败: 连接为空")
		return nil, fmt.Errorf("redis connection is nil")
	}
	if c.Err() != nil {
		zlog.F().Errorf("Redis DoCommonCmd: %v", c.Err())
		return nil, c.Err()
	}
	defer c.Close()

	res, err := c.Do(cmdStr, keysAndArgs...)
	if err == nil {
		return res, err
	}
	return nil, err
}

// CommonLuaScript 执行 Lua 脚本命令
func CommonLuaScript(name, script string, key string, args ...interface{}) (reply interface{}, err error) {
	pool, err := getPool(name)
	if err != nil {
		zlog.F().Errorf("获取 Redis 连接池失败: %v", err)
		return nil, err
	}
	
	if pool == nil {
		zlog.F().Errorf("Redis 连接池为空: %s", name)
		return nil, fmt.Errorf("redis pool is nil for name: %s", name)
	}

	c := pool.Get()
	if c == nil {
		zlog.F().Errorf("获取 Redis 连接失败: 连接为空")
		return nil, fmt.Errorf("redis connection is nil")
	}
	if c.Err() != nil {
		zlog.F().Errorf("LuaCommonCmd get redis error: %v", c.Err())
		return nil, c.Err()
	}
	defer c.Close()

	lua := redis.NewScript(1, script)
	return lua.Do(c, append([]interface{}{key}, args...)...)
}

// WithMaxActive 设置最大活跃连接数
func WithMaxActive(maxActive int) Redis_func {
	return func(r *RedisPool) {
		r.maxActive = maxActive
	}
}

// WithMaxIdle 设置最大空闲连接数
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
