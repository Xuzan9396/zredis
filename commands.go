package zredis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

// RedisCommander 统一的Redis命令接口
type RedisCommander interface {
	// 核心命令方法
	Cmd(cmdStr string, keysAndArgs ...interface{}) (interface{}, error)
	LuaScript(script string, key string, args ...interface{}) (interface{}, error)
	
	// 基础命令
	Get(key string) (interface{}, error)
	Set(key string, val interface{}) (interface{}, error)
	SetEx(key string, val interface{}, timeExpire int64) (interface{}, error)
	SetNx(key string, val interface{}) (interface{}, error)
	SetNxEx(key string, val interface{}, timeExpire int) (interface{}, error)
	Del(key string) (interface{}, error)
	Exists(key string) (interface{}, error)
	Expire(key string, timeInt int) error
	ExpireAt(key string, timestampInt int64) (interface{}, error)
	Keys(pre_key string) (interface{}, error)
	
	// 计数命令  
	IncrBy(key string) (interface{}, error)
	IncrbyVal(key string, val interface{}) (interface{}, error)
	IncrbyFloat(key string, val interface{}) (interface{}, error)
	DecrByNum(key string, num interface{}) (interface{}, error)
	
	// Hash命令
	Hset(key string, field, val interface{}) (interface{}, error)
	Hget(key string, field string) (interface{}, error)
	HgetAll(key string) (interface{}, error)
	Hdel(key string, field interface{}) (interface{}, error)
	Hexists(key string, field interface{}) bool
	HIncrby(key string, field string, val interface{}) (interface{}, error)
	HMget(key string, fields []interface{}) (interface{}, error)
	
	// Set命令
	SAdd(key string, val interface{}) error
	SRem(key string, val interface{}) (interface{}, error)
	SCard(key string) (interface{}, error)
	SIsMember(key string, val interface{}) (bool, error)
	SMembers(key string) (interface{}, error)
	
	// ZSet命令
	ZAdd(key string, score, val interface{}) error
	ZAddBool(key string, score, val interface{}) (int, error)
	ZRem(key string, val interface{}) (interface{}, error)
	ZCard(key string) (interface{}, error)
	ZScore(key string, val interface{}) (interface{}, error)
	ZRange(key string, start, end int, withScore bool) (interface{}, error)
	ZRevRange(key string, start, end int, withScore bool) (interface{}, error)
	ZRangeByScore(key string, start, end interface{}, withScore bool) (interface{}, error)
	ZRevRank(key string, val interface{}) (interface{}, error)
	ZIncrBy(key string, offset interface{}, val interface{}) (interface{}, error)
	ZIncrByExpire(key string, offset interface{}, val interface{}, expireInt int) (interface{}, error)
	
	// List命令
	LPush(key string, val interface{}) (interface{}, error)
	RPop(key string) (interface{}, error)
	BRPop(key string, timeout int) (interface{}, error)
	LLen(key string) (interface{}, error)
	
	// Bit命令
	SetBit(key string, offset, val interface{}) (interface{}, error)
	GetBit(key string, offset interface{}) (interface{}, error)
	BitCount(key string) (interface{}, error)
	
	// 模式删除
	DelPattern(patternKey string) error
}

// 统一的Redis命令实现
type redisCommands struct {
	executor func(cmdStr string, keysAndArgs ...interface{}) (interface{}, error)
	luaExecutor func(script string, key string, args ...interface{}) (interface{}, error)
}

func NewRedisCommands(executor func(string, ...interface{}) (interface{}, error), luaExecutor func(string, string, ...interface{}) (interface{}, error)) RedisCommander {
	return &redisCommands{
		executor: executor,
		luaExecutor: luaExecutor,
	}
}

func (r *redisCommands) Cmd(cmdStr string, keysAndArgs ...interface{}) (interface{}, error) {
	return r.executor(cmdStr, keysAndArgs...)
}

func (r *redisCommands) LuaScript(script string, key string, args ...interface{}) (interface{}, error) {
	return r.luaExecutor(script, key, args...)
}

func (r *redisCommands) Get(key string) (interface{}, error) {
	return r.executor("GET", key)
}

func (r *redisCommands) Set(key string, val interface{}) (interface{}, error) {
	return r.executor("SET", key, val)
}

func (r *redisCommands) SetEx(key string, val interface{}, timeExpire int64) (interface{}, error) {
	return r.executor("SETEX", key, timeExpire, val)
}

func (r *redisCommands) SetNx(key string, val interface{}) (interface{}, error) {
	return r.executor("SETNX", key, val)
}

func (r *redisCommands) SetNxEx(key string, val interface{}, timeExpire int) (interface{}, error) {
	return r.executor("SET", key, val, "EX", timeExpire, "NX")
}

func (r *redisCommands) Del(key string) (interface{}, error) {
	return r.executor("DEL", key)
}

func (r *redisCommands) Exists(key string) (interface{}, error) {
	return r.executor("EXISTS", key)
}

func (r *redisCommands) Expire(key string, timeInt int) error {
	_, err := r.executor("EXPIRE", key, timeInt)
	return err
}

func (r *redisCommands) ExpireAt(key string, timestampInt int64) (interface{}, error) {
	return r.executor("EXPIREAT", key, timestampInt)
}

func (r *redisCommands) Keys(pre_key string) (interface{}, error) {
	return r.executor("KEYS", pre_key)
}

func (r *redisCommands) IncrBy(key string) (interface{}, error) {
	return r.executor("INCR", key)
}

func (r *redisCommands) IncrbyVal(key string, val interface{}) (interface{}, error) {
	return r.executor("INCRBY", key, val)
}

func (r *redisCommands) IncrbyFloat(key string, val interface{}) (interface{}, error) {
	return r.executor("INCRBYFLOAT", key, val)
}

func (r *redisCommands) DecrByNum(key string, num interface{}) (interface{}, error) {
	return r.executor("DECRBY", key, num)
}

func (r *redisCommands) Hset(key string, field, val interface{}) (interface{}, error) {
	return r.executor("HSET", key, field, val)
}

func (r *redisCommands) Hget(key string, field string) (interface{}, error) {
	return r.executor("HGET", key, field)
}

func (r *redisCommands) HgetAll(key string) (interface{}, error) {
	return r.executor("HGETALL", key)
}

func (r *redisCommands) Hdel(key string, field interface{}) (interface{}, error) {
	return r.executor("HDEL", key, field)
}

func (r *redisCommands) Hexists(key string, field interface{}) bool {
	res, _ := redis.Int(r.executor("HEXISTS", key, field))
	return res == 1
}

func (r *redisCommands) HIncrby(key string, field string, val interface{}) (interface{}, error) {
	return r.executor("HINCRBY", key, field, val)
}

func (r *redisCommands) HMget(key string, fields []interface{}) (interface{}, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	args = append(args, fields...)
	return r.executor("HMGET", args...)
}

func (r *redisCommands) SAdd(key string, val interface{}) error {
	_, err := r.executor("SADD", key, val)
	return err
}

func (r *redisCommands) SRem(key string, val interface{}) (interface{}, error) {
	return r.executor("SREM", key, val)
}

func (r *redisCommands) SCard(key string) (interface{}, error) {
	return r.executor("SCARD", key)
}

func (r *redisCommands) SIsMember(key string, val interface{}) (bool, error) {
	existsInt64, err := redis.Int64(r.executor("SISMEMBER", key, val))
	if err != nil {
		return false, err
	}
	return existsInt64 == 1, nil
}

func (r *redisCommands) SMembers(key string) (interface{}, error) {
	return r.executor("SMEMBERS", key)
}

func (r *redisCommands) ZAdd(key string, score, val interface{}) error {
	_, err := r.executor("ZADD", key, score, val)
	return err
}

func (r *redisCommands) ZAddBool(key string, score, val interface{}) (int, error) {
	return redis.Int(r.executor("ZADD", key, score, val))
}

func (r *redisCommands) ZRem(key string, val interface{}) (interface{}, error) {
	return r.executor("ZREM", key, val)
}

func (r *redisCommands) ZCard(key string) (interface{}, error) {
	return r.executor("ZCARD", key)
}

func (r *redisCommands) ZScore(key string, val interface{}) (interface{}, error) {
	return r.executor("ZSCORE", key, val)
}

func (r *redisCommands) ZRange(key string, start, end int, withScore bool) (interface{}, error) {
	args := []interface{}{key, start, end}
	if withScore {
		args = append(args, "WITHSCORES")
	}
	return r.executor("ZRANGE", args...)
}

func (r *redisCommands) ZRevRange(key string, start, end int, withScore bool) (interface{}, error) {
	args := []interface{}{key, start, end}
	if withScore {
		args = append(args, "WITHSCORES")
	}
	return r.executor("ZREVRANGE", args...)
}

func (r *redisCommands) ZRangeByScore(key string, start, end interface{}, withScore bool) (interface{}, error) {
	args := []interface{}{key, start, end}
	if withScore {
		args = append(args, "WITHSCORES")
	}
	return r.executor("ZRANGEBYSCORE", args...)
}

func (r *redisCommands) ZRevRank(key string, val interface{}) (interface{}, error) {
	return r.executor("ZREVRANK", key, val)
}

func (r *redisCommands) ZIncrBy(key string, offset interface{}, val interface{}) (interface{}, error) {
	return r.executor("ZINCRBY", key, offset, val)
}

func (r *redisCommands) ZIncrByExpire(key string, offset interface{}, val interface{}, expireInt int) (interface{}, error) {
	res, err := r.executor("ZINCRBY", key, offset, val)
	if err != nil {
		return nil, err
	}
	
	// 设置过期时间，如果失败则记录错误但不影响主要操作的结果
	_, expireErr := r.executor("EXPIRE", key, expireInt)
	if expireErr != nil {
		// 这里可以根据业务需求决定是否返回错误
		// 目前保持向后兼容，只记录日志但不返回错误
		// 可以考虑使用日志系统记录这个错误
	}
	
	return res, err
}

func (r *redisCommands) LPush(key string, val interface{}) (interface{}, error) {
	return r.executor("LPUSH", key, val)
}

func (r *redisCommands) RPop(key string) (interface{}, error) {
	return r.executor("RPOP", key)
}

func (r *redisCommands) BRPop(key string, timeout int) (interface{}, error) {
	return r.executor("BRPOP", key, timeout)
}

func (r *redisCommands) LLen(key string) (interface{}, error) {
	return r.executor("LLEN", key)
}

func (r *redisCommands) SetBit(key string, offset, val interface{}) (interface{}, error) {
	return r.executor("SETBIT", key, offset, val)
}

func (r *redisCommands) GetBit(key string, offset interface{}) (interface{}, error) {
	return r.executor("GETBIT", key, offset)
}

func (r *redisCommands) BitCount(key string) (interface{}, error) {
	return r.executor("BITCOUNT", key)
}

func (r *redisCommands) DelPattern(patternKey string) error {
	cursor := "0"
	for {
		values, err := redis.Values(r.executor("SCAN", cursor, "MATCH", patternKey))
		if err != nil {
			return err
		}
		
		cursor, _ = redis.String(values[0], nil)
		keys, _ := redis.Strings(values[1], nil)
		
		for _, key := range keys {
			_, err = r.executor("DEL", key)
			if err != nil {
				return err
			}
		}
		
		if cursor == "0" {
			break
		}
	}
	return nil
}

// 缓存相关类型和函数
type FuncType func() ([]byte, error)
type FuncTypeInt func() ([]byte, int64, error)

// CallBackMsgpackCacheWithCommander 带缓存的回调函数，使用指定的命令实例
func CallBackMsgpackCacheWithCommander(commander RedisCommander, redisKey string, funcs FuncType, opt ...int64) ([]byte, error) {
	existsInt, _ := redis.Int(commander.Exists(redisKey))
	if existsInt == 0 {
		res, err := funcs()
		if err != nil {
			return nil, err
		}
		if res == nil {
			return nil, errors.New("data is null")
		}
		var timeOut int64 = 86400
		if len(opt) > 0 {
			timeOut = opt[0]
		}
		commander.SetEx(redisKey, res, timeOut)
		return res, nil
	} else {
		bytesRes, err := redis.Bytes(commander.Get(redisKey))
		if err != nil {
			return nil, err
		}
		return bytesRes, nil
	}
}

// CallBackMsgpackCacheInWithCommander 带内部时间控制的缓存回调函数，使用指定的命令实例
func CallBackMsgpackCacheInWithCommander(commander RedisCommander, redisKey string, funcs FuncTypeInt) ([]byte, error) {
	existsInt, _ := redis.Int(commander.Exists(redisKey))
	if existsInt == 0 {
		res, timeOut, err := funcs()
		if err != nil {
			return nil, err
		}
		if res == nil {
			return nil, errors.New("data is null")
		}
		
		_, err = commander.Set(redisKey, res)
		if err != nil {
			return nil, err
		}
		if timeOut > 0 {
			commander.ExpireAt(redisKey, timeOut)
		}
		return res, nil
	} else {
		bytesRes, err := redis.Bytes(commander.Get(redisKey))
		if err != nil {
			return nil, err
		}
		return bytesRes, nil
	}
}