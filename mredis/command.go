package mredis

import (
	"github.com/garyburd/redigo/redis"
)

// -------------------------  公众函数  -----------------------
func CommonGet(name, key string) (interface{}, error) {
	return CommonCmd(name, "GET", key)

}

func CommonExists(name, key string) (interface{}, error) {
	return CommonCmd(name, "EXISTS", key)

}

func CommonSetNx(name, key string, val interface{}) (interface{}, error) {
	return CommonCmd(name, "SETNX", key, val)
}

func CommonSetNxEx(name, key string, val interface{}, timeExpire int) (interface{}, error) {
	// 返回OK 和 nil
	return CommonCmd(name, "SET", key, val, "EX", timeExpire, "NX")
}

func CommonSet(name, key string, val interface{}) (interface{}, error) {
	return CommonCmd(name, "SET", key, val)

}

func CommonSetEx(name, key string, val interface{}, timeExpire int64) (interface{}, error) {
	return CommonCmd(name, "SETEX", key, timeExpire, val)
}

func CommonSISMEMBER(name, key string, val interface{}) (bools bool, err error) {

	var (
		existsInt64 int64
	)

	existsInt64, err = redis.Int64(CommonCmd(name, "SISMEMBER", key, val))

	if err != nil {
		return
	}

	if existsInt64 == 1 {
		bools = true
		return
	}

	return
}

func CommonZADD(name, key string, score, val interface{}) (err error) {

	_, err = CommonCmd(name, "ZADD", key, score, val)

	if err != nil {
		return
	}

	return
}

func CommonZADDBool(name, key string, score, val interface{}) (boolInt int, err error) {

	boolInt, err = redis.Int(CommonCmd(name, "ZADD", key, score, val))

	if err != nil {
		return
	}

	return
}

func CommonSADD(name, key string, val interface{}) (err error) {

	_, err = CommonCmd(name, "SADD", key, val)

	if err != nil {
		return
	}

	return
}

func CommonEXPIRE(name, key string, timeInt int) (err error) {

	_, err = CommonCmd(name, "EXPIRE", key, timeInt)

	if err != nil {
		return
	}

	return
}

// 倒序获取榜单
func CommonZrevrank(name, key string, val interface{}) (interface{}, error) {
	return CommonCmd(name, "ZREVRANK", key, val)
}

func CommonZCARD(name, key string) (interface{}, error) {
	//avgArr := dealRedis(constvar.ACTIVE_520_HEARTBEAT_VAL_RANK, avgs...)
	return CommonCmd(name, "ZCARD", key)

}

func CommonSCARD(name, key string) (interface{}, error) {
	//avgArr := dealRedis(constvar.ACTIVE_520_HEARTBEAT_VAL_RANK, avgs...)
	return CommonCmd(name, "SCARD", key)

}

func CommonHset(name, key string, key_one, val interface{}) (interface{}, error) {
	return CommonCmd(name, "HSET", key, key_one, val)
}
func CommonZRange(name, key string, start, end int, withScore bool) (interface{}, error) {
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return CommonCmd(name, "ZRANGE", arr...)

}

// 分数排序
func CommonZRevRange(name, key string, start, end int, withScore bool) (interface{}, error) {
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return CommonCmd(name, "ZREVRANGE", arr...)

}

func CommonZrem(name, key string, val interface{}) (interface{}, error) {

	return CommonCmd(name, "ZREM", key, val)
}

func CommonSrem(name, key string, val interface{}) (interface{}, error) {

	return CommonCmd(name, "SREM", key, val)
}

func CommonZRANGEBYSCORE(name, key string, start, end interface{}, withScore bool) (interface{}, error) {

	//ZRANGEBYSCORE RED_RANK:UNION_UNBIND:20201202 -inf 4 WITHSCORES
	arr := []interface{}{
		key, start, end,
	}
	if withScore {
		arr = append(arr, "WITHSCORES")
	}

	return CommonCmd(name, "ZRANGEBYSCORE", arr...)

}

func CommonZscore(name, key string, val interface{}) (interface{}, error) {

	return CommonCmd(name, "ZSCORE", key, val)
}

func CommonSMEMBERS(name, key string) (interface{}, error) {
	return CommonCmd(name, "SMEMBERS", key)

}

func CommonSetBit(name, key string, offset, val interface{}) (interface{}, error) {

	return CommonCmd(name, "SETBIT", key, offset, val)

}

func CommonGetBit(name, key string, offset interface{}) (interface{}, error) {

	return CommonCmd(name, "GETBIT", key, offset)

}

func CommonBitCount(name, key string) (interface{}, error) {
	return CommonCmd(name, "BITCOUNT", key)
}

func CommonDel(name, key string) (interface{}, error) {
	return CommonCmd(name, "del", key)
}

func CommonZIncrBy(name, key string, offest interface{}, val interface{}) (interface{}, error) {

	return CommonCmd(name, "ZINCRBY", key, offest, val)

}

func CommonZIncrByExpire(name, key string, offest interface{}, val interface{}, expireInt int) (interface{}, error) {
	res, err := CommonCmd(name, "ZINCRBY", key, offest, val)
	CommonCmd(name, "EXPIRE", key, expireInt)
	return res, err

}

func CommonIncrBy(name, key string) (interface{}, error) {
	return CommonCmd(name, "INCR", key)
}

func CommonKeys(name, key string) (interface{}, error) {
	return CommonCmd(name, "KEYS", key)

}

func CommonHget(name, key string, val string) (interface{}, error) {
	return CommonCmd(name, "HGET", key, val)
}

func CommonDECRBYByNum(name, key string, num interface{}) (interface{}, error) {
	return CommonCmd(name, "DECRBY", key, num)

}

// hdel
func CommonHdel(name, key string, val interface{}) (interface{}, error) {
	return CommonCmd(name, "HDEL", key, val)
}

// HEXISTS
func CommonHexists(name, key string, val interface{}) bool {

	res, _ := redis.Int(CommonCmd(name, "HEXISTS", key, val))
	if res == 1 {
		return true
	}

	return false

}

func CommonHIncrby(name, key string, field string, val any) (interface{}, error) {
	return CommonCmd(name, "HINCRBY", key, field, val)

}

func CommonIncrby(name, key string, val any) (interface{}, error) {
	return CommonCmd(name, "INCRBY", key, val)

}

// floatIncrby
func CommonIncrbyFloat(name, key string, val any) (interface{}, error) {

	return CommonCmd(name, "INCRBYFLOAT", key, val)

}

func CommonHMget(name, key string, field []interface{}) (interface{}, error) {
	fields := make([]interface{}, 0)
	fields = append(fields, key)
	fields = append(fields, field...)
	return CommonCmd(name, "HMGET", fields...)
}

func CommonHgetAll(name, key string) (interface{}, error) {
	return CommonCmd(name, "HGETALL", key)

}

func CommonLPush(name, key string, val interface{}) (interface{}, error) {
	return CommonCmd(name, "LPUSH", key, val)
}

func CommonRPop(name, key string) (interface{}, error) {
	return CommonCmd(name, "RPOP", key)
}

// 该命令会阻塞
func CommonBRpop(name, key string, timeout int) (interface{}, error) {
	return CommonCmd(name, "BRPOP", key, timeout)
}

func CommonLLen(name, key string) (interface{}, error) {
	return CommonCmd(name, "LLEN", key)
}
