package sredis

import (
	"github.com/Xuzan9396/zredis"
)

// 为RedisPool添加统一命令接口
func (c *RedisPool) GetCommander() zredis.RedisCommander {
	return zredis.NewRedisCommands(
		func(cmdStr string, keysAndArgs ...interface{}) (interface{}, error) {
			return c.CommonCmd(cmdStr, keysAndArgs...)
		},
		func(script string, key string, args ...interface{}) (interface{}, error) {
			return c.CommonLuaScript(script, key, args...)
		},
	)
}

// -------------------------  公众函数  -----------------------
// 向后兼容的包装方法

func (c *RedisPool) CommonGet(key string) (interface{}, error) {
	return c.GetCommander().Get(key)
}

func (c *RedisPool) CommonExists(key string) (interface{}, error) {
	return c.GetCommander().Exists(key)
}

func (c *RedisPool) CommonSetNx(key string, val interface{}) (interface{}, error) {
	return c.GetCommander().SetNx(key, val)
}

func (c *RedisPool) CommonSetNxEx(key string, val interface{}, timeExpire int) (interface{}, error) {
	return c.GetCommander().SetNxEx(key, val, timeExpire)
}

func (c *RedisPool) CommonSet(key string, val interface{}) (interface{}, error) {
	return c.GetCommander().Set(key, val)
}

func (c *RedisPool) CommonSetEx(key string, val interface{}, timeExpire int64) (interface{}, error) {
	return c.GetCommander().SetEx(key, val, timeExpire)
}

func (c *RedisPool) CommonSISMEMBER(key string, val interface{}) (bools bool, err error) {
	return c.GetCommander().SIsMember(key, val)
}

func (c *RedisPool) CommonZADD(key string, score, val interface{}) (err error) {
	return c.GetCommander().ZAdd(key, score, val)
}

func (c *RedisPool) CommonZADDBool(key string, score, val interface{}) (boolInt int, err error) {
	return c.GetCommander().ZAddBool(key, score, val)
}

func (c *RedisPool) CommonSADD(key string, val interface{}) (err error) {
	return c.GetCommander().SAdd(key, val)
}

func (c *RedisPool) CommonEXPIRE(key string, timeInt int) (err error) {
	return c.GetCommander().Expire(key, timeInt)
}

func (c *RedisPool) CommonEXPIREAT(key string, timestampInt int64) (res interface{}, err error) {
	return c.GetCommander().ExpireAt(key, timestampInt)
}

func (c *RedisPool) CommonZrevrank(key string, val interface{}) (interface{}, error) {
	return c.GetCommander().ZRevRank(key, val)
}

func (c *RedisPool) CommonZCARD(key string) (interface{}, error) {
	return c.GetCommander().ZCard(key)
}

func (c *RedisPool) CommonSCARD(key string) (interface{}, error) {
	return c.GetCommander().SCard(key)
}

func (c *RedisPool) CommonHset(key string, key_one, val interface{}) (interface{}, error) {
	return c.GetCommander().Hset(key, key_one, val)
}

func (c *RedisPool) CommonZRange(key string, start, end int, withScore bool) (interface{}, error) {
	return c.GetCommander().ZRange(key, start, end, withScore)
}

func (c *RedisPool) CommonZRevRange(key string, start, end int, withScore bool) (interface{}, error) {
	return c.GetCommander().ZRevRange(key, start, end, withScore)
}

func (c *RedisPool) CommonZrem(key string, val interface{}) (interface{}, error) {
	return c.GetCommander().ZRem(key, val)
}

func (c *RedisPool) CommonSrem(key string, val interface{}) (interface{}, error) {
	return c.GetCommander().SRem(key, val)
}

func (c *RedisPool) CommonZRANGEBYSCORE(key string, start, end interface{}, withScore bool) (interface{}, error) {
	return c.GetCommander().ZRangeByScore(key, start, end, withScore)
}

func (c *RedisPool) CommonZscore(redisKey string, val interface{}) (interface{}, error) {
	return c.GetCommander().ZScore(redisKey, val)
}

func (c *RedisPool) CommonSMEMBERS(key string) (interface{}, error) {
	return c.GetCommander().SMembers(key)
}

func (c *RedisPool) CommonSetBit(key string, offset, val interface{}) (interface{}, error) {
	return c.GetCommander().SetBit(key, offset, val)
}

func (c *RedisPool) CommonGetBit(key string, offset interface{}) (interface{}, error) {
	return c.GetCommander().GetBit(key, offset)
}

func (c *RedisPool) CommonBitCount(key string) (interface{}, error) {
	return c.GetCommander().BitCount(key)
}

func (c *RedisPool) CommonDel(key string) (interface{}, error) {
	return c.GetCommander().Del(key)
}

func (c *RedisPool) CommonZIncrBy(key string, offest interface{}, val interface{}) (interface{}, error) {
	return c.GetCommander().ZIncrBy(key, offest, val)
}

func (c *RedisPool) CommonZIncrByExpire(key string, offest interface{}, val interface{}, expireInt int) (interface{}, error) {
	return c.GetCommander().ZIncrByExpire(key, offest, val, expireInt)
}

func (c *RedisPool) CommonIncrBy(key string) (interface{}, error) {
	return c.GetCommander().IncrBy(key)
}

func (c *RedisPool) CommonKeys(pre_key string) (interface{}, error) {
	return c.GetCommander().Keys(pre_key)
}

func (c *RedisPool) CommonHget(pre_key string, val string) (interface{}, error) {
	return c.GetCommander().Hget(pre_key, val)
}

func (c *RedisPool) CommonDECRBYByNum(key string, num interface{}) (interface{}, error) {
	return c.GetCommander().DecrByNum(key, num)
}

func (c *RedisPool) CommonHdel(pre_key string, val interface{}) (interface{}, error) {
	return c.GetCommander().Hdel(pre_key, val)
}

func (c *RedisPool) CommonHexists(pre_key string, val interface{}) bool {
	return c.GetCommander().Hexists(pre_key, val)
}

func (c *RedisPool) CommonHIncrby(key string, field string, val any) (interface{}, error) {
	return c.GetCommander().HIncrby(key, field, val)
}

func (c *RedisPool) CommonIncrby(key string, val any) (interface{}, error) {
	return c.GetCommander().IncrbyVal(key, val)
}

func (c *RedisPool) CommonIncrbyFloat(key string, val any) (interface{}, error) {
	return c.GetCommander().IncrbyFloat(key, val)
}

func (c *RedisPool) CommonHMget(key string, field []interface{}) (interface{}, error) {
	return c.GetCommander().HMget(key, field)
}

func (c *RedisPool) CommonHgetAll(key string) (interface{}, error) {
	return c.GetCommander().HgetAll(key)
}

func (c *RedisPool) CommonLPush(key string, val interface{}) (interface{}, error) {
	return c.GetCommander().LPush(key, val)
}

func (c *RedisPool) CommonRPop(key string) (interface{}, error) {
	return c.GetCommander().RPop(key)
}

func (c *RedisPool) CommonBRpop(key string, timeout int) (interface{}, error) {
	return c.GetCommander().BRPop(key, timeout)
}

func (c *RedisPool) CommonLLen(key string) (interface{}, error) {
	return c.GetCommander().LLen(key)
}

// 缓存相关方法
func (c *RedisPool) CallBackMsgpackCache(redisKey string, funcs zredis.FuncType, opt ...int64) ([]byte, error) {
	var timeOut int64 = 86400
	if len(opt) > 0 {
		timeOut = opt[0]
	}
	return zredis.CallBackMsgpackCacheWithCommander(c.GetCommander(), redisKey, funcs, timeOut)
}

func (c *RedisPool) CallBackMsgpackCacheIn(redisKey string, funcs zredis.FuncTypeInt) ([]byte, error) {
	return zredis.CallBackMsgpackCacheInWithCommander(c.GetCommander(), redisKey, funcs)
}

// 批量删除key
func (c *RedisPool) CommonDelPattern(patternKey string) (err error) {
	return c.GetCommander().DelPattern(patternKey)
}
