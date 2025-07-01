package mredis

import (
	"github.com/Xuzan9396/zredis"
)

// 获取指定名称的Redis命令实例
func GetCommander(name string) zredis.RedisCommander {
	return zredis.NewRedisCommands(
		func(cmdStr string, keysAndArgs ...interface{}) (interface{}, error) {
			return CommonCmd(name, cmdStr, keysAndArgs...)
		},
		func(script string, key string, args ...interface{}) (interface{}, error) {
			return CommonLuaScript(name, script, key, args...)
		},
	)
}

// -------------------------  公众函数  -----------------------
// 向后兼容的包装函数

func CommonGet(name, key string) (interface{}, error) {
	return GetCommander(name).Get(key)
}

func CommonExists(name, key string) (interface{}, error) {
	return GetCommander(name).Exists(key)
}

func CommonSetNx(name, key string, val interface{}) (interface{}, error) {
	return GetCommander(name).SetNx(key, val)
}

func CommonSetNxEx(name, key string, val interface{}, timeExpire int) (interface{}, error) {
	return GetCommander(name).SetNxEx(key, val, timeExpire)
}

func CommonSet(name, key string, val interface{}) (interface{}, error) {
	return GetCommander(name).Set(key, val)
}

func CommonSetEx(name, key string, val interface{}, timeExpire int64) (interface{}, error) {
	return GetCommander(name).SetEx(key, val, timeExpire)
}

func CommonSISMEMBER(name, key string, val interface{}) (bools bool, err error) {
	return GetCommander(name).SIsMember(key, val)
}

func CommonZADD(name, key string, score, val interface{}) (err error) {
	return GetCommander(name).ZAdd(key, score, val)
}

func CommonZADDBool(name, key string, score, val interface{}) (boolInt int, err error) {
	return GetCommander(name).ZAddBool(key, score, val)
}

func CommonSADD(name, key string, val interface{}) (err error) {
	return GetCommander(name).SAdd(key, val)
}

func CommonEXPIRE(name, key string, timeInt int) (err error) {
	return GetCommander(name).Expire(key, timeInt)
}

func CommonZrevrank(name, key string, val interface{}) (interface{}, error) {
	return GetCommander(name).ZRevRank(key, val)
}

func CommonZCARD(name, key string) (interface{}, error) {
	return GetCommander(name).ZCard(key)
}

func CommonSCARD(name, key string) (interface{}, error) {
	return GetCommander(name).SCard(key)
}

func CommonHset(name, key string, key_one, val interface{}) (interface{}, error) {
	return GetCommander(name).Hset(key, key_one, val)
}

func CommonZRange(name, key string, start, end int, withScore bool) (interface{}, error) {
	return GetCommander(name).ZRange(key, start, end, withScore)
}

func CommonZRevRange(name, key string, start, end int, withScore bool) (interface{}, error) {
	return GetCommander(name).ZRevRange(key, start, end, withScore)
}

func CommonZrem(name, key string, val interface{}) (interface{}, error) {
	return GetCommander(name).ZRem(key, val)
}

func CommonSrem(name, key string, val interface{}) (interface{}, error) {
	return GetCommander(name).SRem(key, val)
}

func CommonZRANGEBYSCORE(name, key string, start, end interface{}, withScore bool) (interface{}, error) {
	return GetCommander(name).ZRangeByScore(key, start, end, withScore)
}

func CommonZscore(name, key string, val interface{}) (interface{}, error) {
	return GetCommander(name).ZScore(key, val)
}

func CommonSMEMBERS(name, key string) (interface{}, error) {
	return GetCommander(name).SMembers(key)
}

func CommonSetBit(name, key string, offset, val interface{}) (interface{}, error) {
	return GetCommander(name).SetBit(key, offset, val)
}

func CommonGetBit(name, key string, offset interface{}) (interface{}, error) {
	return GetCommander(name).GetBit(key, offset)
}

func CommonBitCount(name, key string) (interface{}, error) {
	return GetCommander(name).BitCount(key)
}

func CommonDel(name, key string) (interface{}, error) {
	return GetCommander(name).Del(key)
}

func CommonZIncrBy(name, key string, offest interface{}, val interface{}) (interface{}, error) {
	return GetCommander(name).ZIncrBy(key, offest, val)
}

func CommonZIncrByExpire(name, key string, offest interface{}, val interface{}, expireInt int) (interface{}, error) {
	return GetCommander(name).ZIncrByExpire(key, offest, val, expireInt)
}

func CommonIncrBy(name, key string) (interface{}, error) {
	return GetCommander(name).IncrBy(key)
}

func CommonKeys(name, key string) (interface{}, error) {
	return GetCommander(name).Keys(key)
}

func CommonHget(name, key string, val string) (interface{}, error) {
	return GetCommander(name).Hget(key, val)
}

func CommonDECRBYByNum(name, key string, num interface{}) (interface{}, error) {
	return GetCommander(name).DecrByNum(key, num)
}

func CommonHdel(name, key string, val interface{}) (interface{}, error) {
	return GetCommander(name).Hdel(key, val)
}

func CommonHexists(name, key string, val interface{}) bool {
	return GetCommander(name).Hexists(key, val)
}

func CommonHIncrby(name, key string, field string, val any) (interface{}, error) {
	return GetCommander(name).HIncrby(key, field, val)
}

func CommonIncrby(name, key string, val any) (interface{}, error) {
	return GetCommander(name).IncrbyVal(key, val)
}

func CommonIncrbyFloat(name, key string, val any) (interface{}, error) {
	return GetCommander(name).IncrbyFloat(key, val)
}

func CommonHMget(name, key string, field []interface{}) (interface{}, error) {
	return GetCommander(name).HMget(key, field)
}

func CommonHgetAll(name, key string) (interface{}, error) {
	return GetCommander(name).HgetAll(key)
}

func CommonLPush(name, key string, val interface{}) (interface{}, error) {
	return GetCommander(name).LPush(key, val)
}

func CommonRPop(name, key string) (interface{}, error) {
	return GetCommander(name).RPop(key)
}

func CommonBRpop(name, key string, timeout int) (interface{}, error) {
	return GetCommander(name).BRPop(key, timeout)
}

func CommonLLen(name, key string) (interface{}, error) {
	return GetCommander(name).LLen(key)
}

// 批量删除key
func CommonDelPattern(name, patternKey string) (err error) {
	return GetCommander(name).DelPattern(patternKey)
}

// 缓存相关函数
func CallBackMsgpackCache(name, redisKey string, funcs zredis.FuncType, opt ...int64) ([]byte, error) {
	var timeOut int64 = 86400
	if len(opt) > 0 {
		timeOut = opt[0]
	}
	return zredis.CallBackMsgpackCacheWithCommander(GetCommander(name), redisKey, funcs, timeOut)
}

func CallBackMsgpackCacheIn(name, redisKey string, funcs zredis.FuncTypeInt) ([]byte, error) {
	return zredis.CallBackMsgpackCacheInWithCommander(GetCommander(name), redisKey, funcs)
}
