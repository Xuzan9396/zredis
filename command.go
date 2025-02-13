package zredis

import (
	"github.com/garyburd/redigo/redis"
)

// -------------------------  公众函数  -----------------------

func (c *RedisPool) CommonGet(key string, databases ...string) (interface{}, error) {
	return c.CommonCmd("GET", key)

}

func (c *RedisPool) CommonExists(key string, databases ...string) (interface{}, error) {
	return c.CommonCmd("EXISTS", key)

}

func (c *RedisPool) CommonSetNx(key string, val interface{}, databases ...string) (interface{}, error) {
	return c.CommonCmd("SETNX", key, val)
}

func (c *RedisPool) CommonSetNxEx(key string, val interface{}, timeExpire int) (interface{}, error) {
	// 返回OK 和 nil
	return c.CommonCmd("SET", key, val, "EX", timeExpire, "NX")
}

func (c *RedisPool) CommonSet(key string, val interface{}) (interface{}, error) {
	return c.CommonCmd("SET", key, val)

}

func (c *RedisPool) CommonSetEx(key string, val interface{}, timeExpire int64) (interface{}, error) {
	return c.CommonCmd("SETEX", key, timeExpire, val)
}

func (c *RedisPool) CommonSISMEMBER(key string, val interface{}) (bools bool, err error) {

	var (
		existsInt64 int64
	)

	existsInt64, err = redis.Int64(c.CommonCmd("SISMEMBER", key, val))

	if err != nil {
		return
	}

	if existsInt64 == 1 {
		bools = true
		return
	}

	return
}

func (c *RedisPool) CommonZADD(key string, score, val interface{}) (err error) {

	_, err = c.CommonCmd("ZADD", key, score, val)

	if err != nil {
		return
	}

	return
}

func (c *RedisPool) CommonZADDBool(key string, score, val interface{}) (boolInt int, err error) {

	boolInt, err = redis.Int(c.CommonCmd("ZADD", key, score, val))

	if err != nil {
		return
	}

	return
}

func (c *RedisPool) CommonSADD(key string, val interface{}) (err error) {

	_, err = c.CommonCmd("SADD", key, val)

	if err != nil {
		return
	}

	return
}

func (c *RedisPool) CommonEXPIRE(key string, timeInt int) (err error) {

	_, err = c.CommonCmd("EXPIRE", key, timeInt)

	if err != nil {
		return
	}

	return
}

// 倒序获取榜单
func (c *RedisPool) CommonZrevrank(key string, val interface{}) (interface{}, error) {
	return c.CommonCmd("ZREVRANK", key, val)
}

func (c *RedisPool) CommonZCARD(key string) (interface{}, error) {
	//avgArr := dealRedis(constvar.ACTIVE_520_HEARTBEAT_VAL_RANK, avgs...)
	return c.CommonCmd("ZCARD", key)

}

func (c *RedisPool) CommonSCARD(key string) (interface{}, error) {
	//avgArr := dealRedis(constvar.ACTIVE_520_HEARTBEAT_VAL_RANK, avgs...)
	return c.CommonCmd("SCARD", key)

}

func (c *RedisPool) CommonHset(key string, key_one, val interface{}) (interface{}, error) {
	return c.CommonCmd("HSET", key, key_one, val)

}
func (c *RedisPool) CommonZRange(key string, start, end int, withScore bool) (interface{}, error) {
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return c.CommonCmd("ZRANGE", arr...)

}

// 分数排序
func (c *RedisPool) CommonZRevRange(key string, start, end int, withScore bool) (interface{}, error) {
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return c.CommonCmd("ZREVRANGE", arr...)

}

func (c *RedisPool) CommonZrem(key string, val interface{}) (interface{}, error) {

	return c.CommonCmd("ZREM", key, val)
}

func (c *RedisPool) CommonSrem(key string, val interface{}) (interface{}, error) {

	return c.CommonCmd("SREM", key, val)
}

func (c *RedisPool) CommonZRANGEBYSCORE(key string, start, end interface{}, withScore bool) (interface{}, error) {

	//ZRANGEBYSCORE RED_RANK:UNION_UNBIND:20201202 -inf 4 WITHSCORES
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return c.CommonCmd("ZRANGEBYSCORE", arr...)

}

func (c *RedisPool) CommonZscore(redisKey string, val interface{}) (interface{}, error) {

	return c.CommonCmd("ZSCORE", redisKey, val)
}

func (c *RedisPool) CommonSMEMBERS(key string) (interface{}, error) {
	return c.CommonCmd("SMEMBERS", key)

}

func (c *RedisPool) CommonSetBit(key string, offset, val interface{}) (interface{}, error) {

	return c.CommonCmd("SETBIT", key, offset, val)

}

func (c *RedisPool) CommonGetBit(key string, offset interface{}) (interface{}, error) {

	return c.CommonCmd("GETBIT", key, offset)

}

func (c *RedisPool) CommonBitCount(key string) (interface{}, error) {
	return c.CommonCmd("BITCOUNT", key)
}

func (c *RedisPool) CommonDel(key string) (interface{}, error) {
	return c.CommonCmd("del", key)
}

func (c *RedisPool) CommonZIncrBy(key string, offest interface{}, val interface{}) (interface{}, error) {

	return c.CommonCmd("ZINCRBY", key, offest, val)

}

func (c *RedisPool) CommonZIncrByExpire(key string, offest interface{}, val interface{}, expireInt int) (interface{}, error) {
	res, err := c.CommonCmd("ZINCRBY", key, offest, val)
	c.CommonCmd("EXPIRE", key, expireInt)
	return res, err

}

func (c *RedisPool) CommonIncrBy(key string) (interface{}, error) {
	return c.CommonCmd("INCR", key)
}

func (c *RedisPool) CommonKeys(pre_key string) (interface{}, error) {
	return c.CommonCmd("KEYS", pre_key)

}

func (c *RedisPool) CommonHget(pre_key string, val string) (interface{}, error) {
	return c.CommonCmd("HGET", pre_key, val)
}

func (c *RedisPool) CommonDECRBYByNum(key string, num interface{}) (interface{}, error) {
	return c.CommonCmd("DECRBY", key, num)

}

// hdel
func (c *RedisPool) CommonHdel(pre_key string, val interface{}) (interface{}, error) {
	return c.CommonCmd("HDEL", pre_key, val)
}

// HEXISTS
func (c *RedisPool) CommonHexists(pre_key string, val interface{}) bool {

	res, _ := redis.Int(c.CommonCmd("HEXISTS", pre_key, val))
	if res == 1 {
		return true
	}

	return false

}

func (c *RedisPool) CommonHIncrby(key string, field string, val any) (interface{}, error) {
	return c.CommonCmd("HINCRBY", key, field, val)

}

func (c *RedisPool) CommonIncrby(key string, val any) (interface{}, error) {
	return c.CommonCmd("INCRBY", key, val)

}

// floatIncrby
func (c *RedisPool) CommonIncrbyFloat(key string, val any) (interface{}, error) {

	return c.CommonCmd("INCRBYFLOAT", key, val)

}

func (c *RedisPool) CommonHMget(key string, field []interface{}) (interface{}, error) {
	fields := make([]interface{}, 0)
	fields = append(fields, key)
	fields = append(fields, field...)
	return c.CommonCmd("HMGET", fields...)
}

func (c *RedisPool) CommonHgetAll(key string) (interface{}, error) {
	return c.CommonCmd("HGETALL", key)

}

func (c *RedisPool) CommonLPush(key string, val interface{}) (interface{}, error) {
	return c.CommonCmd("LPUSH", key, val)
}

func (c *RedisPool) CommonRPop(key string) (interface{}, error) {
	return c.CommonCmd("RPOP", key)
}

// 该命令会阻塞
func (c *RedisPool) CommonBRpop(key string, timeout int) (interface{}, error) {
	return c.CommonCmd("BRPOP", key, timeout)
}

func (c *RedisPool) CommonLLen(key string) (interface{}, error) {
	return c.CommonCmd("LLEN", key)
}
