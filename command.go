package zredis

import (
	"github.com/garyburd/redigo/redis"
)

// -------------------------  公众函数  -----------------------
func CommonGet(key string, databases ...string) (interface{}, error) {
	return CommonCmd("GET", key)

}

func CommonExists(key string, databases ...string) (interface{}, error) {
	return CommonCmd("EXISTS", key)

}

func CommonSetNx(key string, val interface{}, databases ...string) (interface{}, error) {
	return CommonCmd("SETNX", key, val)
}

func CommonSetNxEx(key string, val interface{}, timeExpire int) (interface{}, error) {
	// 返回OK 和 nil
	return CommonCmd("SET", key, val, "EX", timeExpire, "NX")
}

func CommonSet(key string, val interface{}) (interface{}, error) {
	return CommonCmd("SET", key, val)

}

func CommonSetEx(key string, val interface{}, timeExpire int64) (interface{}, error) {
	return CommonCmd("SETEX", key, timeExpire, val)
}

func CommonSISMEMBER(key string, val interface{}) (bools bool, err error) {

	var (
		existsInt64 int64
	)

	existsInt64, err = redis.Int64(CommonCmd("SISMEMBER", key, val))

	if err != nil {
		return
	}

	if existsInt64 == 1 {
		bools = true
		return
	}

	return
}

func CommonZADD(key string, score, val interface{}) (err error) {

	_, err = CommonCmd("ZADD", key, score, val)

	if err != nil {
		return
	}

	return
}

func CommonZADDBool(key string, score, val interface{}) (boolInt int, err error) {

	boolInt, err = redis.Int(CommonCmd("ZADD", key, score, val))

	if err != nil {
		return
	}

	return
}

func CommonSADD(key string, val interface{}) (err error) {

	_, err = CommonCmd("SADD", key, val)

	if err != nil {
		return
	}

	return
}

func CommonEXPIRE(key string, timeInt int) (err error) {

	_, err = CommonCmd("EXPIRE", key, timeInt)

	if err != nil {
		return
	}

	return
}

// 倒序获取榜单
func CommonZrevrank(key string, val interface{}) (interface{}, error) {
	return CommonCmd("ZREVRANK", key, val)
}

func CommonZCARD(key string) (interface{}, error) {
	//avgArr := dealRedis(constvar.ACTIVE_520_HEARTBEAT_VAL_RANK, avgs...)
	return CommonCmd("ZCARD", key)

}

func CommonSCARD(key string) (interface{}, error) {
	//avgArr := dealRedis(constvar.ACTIVE_520_HEARTBEAT_VAL_RANK, avgs...)
	return CommonCmd("SCARD", key)

}

func CommonHset(key string, key_one, val interface{}) (interface{}, error) {
	return CommonCmd("HSET", key, key_one, val)

}
func CommonZRange(key string, start, end int, withScore bool) (interface{}, error) {
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return CommonCmd("ZRANGE", arr...)

}

// 分数排序
func CommonZRevRange(key string, start, end int, withScore bool) (interface{}, error) {
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return CommonCmd("ZREVRANGE", arr...)

}

func CommonZrem(key string, val interface{}) (interface{}, error) {

	return CommonCmd("ZREM", key, val)
}

func CommonSrem(key string, val interface{}) (interface{}, error) {

	return CommonCmd("SREM", key, val)
}

func CommonZRANGEBYSCORE(key string, start, end interface{}, withScore bool) (interface{}, error) {

	//ZRANGEBYSCORE RED_RANK:UNION_UNBIND:20201202 -inf 4 WITHSCORES
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return CommonCmd("ZRANGEBYSCORE", arr...)

}

func CommonZscore(redisKey string, val interface{}) (interface{}, error) {

	return CommonCmd("ZSCORE", redisKey, val)
}

func CommonSMEMBERS(key string) (interface{}, error) {
	return CommonCmd("SMEMBERS", key)

}

func CommonSetBit(key string, offset, val interface{}) (interface{}, error) {

	return CommonCmd("SETBIT", key, offset, val)

}

func CommonGetBit(key string, offset interface{}) (interface{}, error) {

	return CommonCmd("GETBIT", key, offset)

}

func CommonBitCount(key string) (interface{}, error) {
	return CommonCmd("BITCOUNT", key)
}

func CommonDel(key string) (interface{}, error) {
	return CommonCmd("del", key)
}

func CommonZIncrBy(key string, offest interface{}, val interface{}) (interface{}, error) {

	return CommonCmd("ZINCRBY", key, offest, val)

}

func CommonZIncrByExpire(key string, offest interface{}, val interface{}, expireInt int) (interface{}, error) {
	res, err := CommonCmd("ZINCRBY", key, offest, val)
	CommonCmd("EXPIRE", key, expireInt)
	return res, err

}

func CommonIncrBy(key string) (interface{}, error) {
	return CommonCmd("INCR", key)
}

func CommonKeys(pre_key string) (interface{}, error) {
	return CommonCmd("KEYS", pre_key)

}

func CommonHget(pre_key string, val string) (interface{}, error) {
	return CommonCmd("HGET", pre_key, val)
}

func CommonDECRBYByNum(key string, num interface{}) (interface{}, error) {
	return CommonCmd("DECRBY", key, num)

}

// hdel
func CommonHdel(pre_key string, val interface{}) (interface{}, error) {
	return CommonCmd("HDEL", pre_key, val)
}

// HEXISTS
func CommonHexists(pre_key string, val interface{}) bool {

	res, _ := redis.Int(CommonCmd("HEXISTS", pre_key, val))
	if res == 1 {
		return true
	}

	return false

}

func CommonHIncrby(key string, field string, val any) (interface{}, error) {
	return CommonCmd("HINCRBY", key, field, val)

}

func CommonIncrby(key string, val any) (interface{}, error) {
	return CommonCmd("INCRBY", key, val)

}

// floatIncrby
func CommonIncrbyFloat(key string, val any) (interface{}, error) {

	return CommonCmd("INCRBYFLOAT", key, val)

}

func CommonHMget(key string, field []interface{}) (interface{}, error) {
	fields := make([]interface{}, 0)
	fields = append(fields, key)
	fields = append(fields, field...)
	return CommonCmd("HMGET", fields...)
}

func CommonHgetAll(key string) (interface{}, error) {
	return CommonCmd("HGETALL", key)

}

func CommonLPush(key string, val interface{}) (interface{}, error) {
	return CommonCmd("LPUSH", key, val)
}

func CommonRPop(key string) (interface{}, error) {
	return CommonCmd("RPOP", key)
}

// 该命令会阻塞
func CommonBRpop(key string, timeout int) (interface{}, error) {
	return CommonCmd("BRPOP", key, timeout)
}

func CommonLLen(key string) (interface{}, error) {
	return CommonCmd("LLEN", key)
}
