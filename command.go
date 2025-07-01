package zredis

import (
	"sync"
)

// 全局Redis命令实例，基于统一的命令接口
var (
	globalCommander RedisCommander
	once            sync.Once
)

// 初始化全局命令实例，使用sync.Once确保线程安全
func initGlobalCommander() {
	once.Do(func() {
		globalCommander = NewRedisCommands(CommonCmd, CommonLuaScript)
	})
}

// -------------------------  公众函数  -----------------------
// 使用统一命令接口的包装函数，保持向后兼容性
func CommonGet(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.Get(key)
}

func CommonExists(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.Exists(key)
}

func CommonSetNx(key string, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.SetNx(key, val)
}

func CommonSetNxEx(key string, val interface{}, timeExpire int) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.SetNxEx(key, val, timeExpire)
}

func CommonSet(key string, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.Set(key, val)
}

func CommonSetEx(key string, val interface{}, timeExpire int64) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.SetEx(key, val, timeExpire)
}

func CommonSISMEMBER(key string, val interface{}) (bools bool, err error) {
	initGlobalCommander()
	return globalCommander.SIsMember(key, val)
}

func CommonZADD(key string, score, val interface{}) (err error) {
	initGlobalCommander()
	return globalCommander.ZAdd(key, score, val)
}

func CommonZADDBool(key string, score, val interface{}) (boolInt int, err error) {
	initGlobalCommander()
	return globalCommander.ZAddBool(key, score, val)
}

func CommonSADD(key string, val interface{}) (err error) {
	initGlobalCommander()
	return globalCommander.SAdd(key, val)
}

func CommonEXPIRE(key string, timeInt int) (err error) {
	initGlobalCommander()
	return globalCommander.Expire(key, timeInt)
}

func CommonEXPIREAT(key string, timestampInt int64) (res interface{}, err error) {
	initGlobalCommander()
	return globalCommander.ExpireAt(key, timestampInt)
}

// 倒序获取榜单
func CommonZrevrank(key string, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZRevRank(key, val)
}

func CommonZCARD(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZCard(key)
}

func CommonSCARD(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.SCard(key)
}

func CommonHset(key string, key_one, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.Hset(key, key_one, val)
}

func CommonZRange(key string, start, end int, withScore bool) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZRange(key, start, end, withScore)
}

// 分数排序
func CommonZRevRange(key string, start, end int, withScore bool) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZRevRange(key, start, end, withScore)
}

func CommonZrem(key string, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZRem(key, val)
}

func CommonSrem(key string, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.SRem(key, val)
}

func CommonZRANGEBYSCORE(key string, start, end interface{}, withScore bool) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZRangeByScore(key, start, end, withScore)
}

func CommonZscore(redisKey string, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZScore(redisKey, val)
}

func CommonSMEMBERS(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.SMembers(key)
}

func CommonSetBit(key string, offset, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.SetBit(key, offset, val)
}

func CommonGetBit(key string, offset interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.GetBit(key, offset)
}

func CommonBitCount(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.BitCount(key)
}

func CommonDel(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.Del(key)
}

func CommonZIncrBy(key string, offest interface{}, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZIncrBy(key, offest, val)
}

func CommonZIncrByExpire(key string, offest interface{}, val interface{}, expireInt int) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.ZIncrByExpire(key, offest, val, expireInt)
}

func CommonIncrBy(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.IncrBy(key)
}

func CommonKeys(pre_key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.Keys(pre_key)
}

func CommonHget(pre_key string, val string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.Hget(pre_key, val)
}

func CommonDECRBYByNum(key string, num interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.DecrByNum(key, num)
}

// hdel
func CommonHdel(pre_key string, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.Hdel(pre_key, val)
}

// HEXISTS
func CommonHexists(pre_key string, val interface{}) bool {
	initGlobalCommander()
	return globalCommander.Hexists(pre_key, val)
}

func CommonHIncrby(key string, field string, val any) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.HIncrby(key, field, val)
}

func CommonIncrby(key string, val any) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.IncrbyVal(key, val)
}

// floatIncrby
func CommonIncrbyFloat(key string, val any) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.IncrbyFloat(key, val)
}

func CommonHMget(key string, field []interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.HMget(key, field)
}

func CommonHgetAll(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.HgetAll(key)
}

func CommonLPush(key string, val interface{}) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.LPush(key, val)
}

func CommonRPop(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.RPop(key)
}

// 该命令会阻塞
func CommonBRpop(key string, timeout int) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.BRPop(key, timeout)
}

func CommonLLen(key string) (interface{}, error) {
	initGlobalCommander()
	return globalCommander.LLen(key)
}

func CallBackMsgpackCache[T int64](redisKey string, funcs FuncType, opt ...T) ([]byte, error) {
	initGlobalCommander()
	var timeOut int64 = 86400
	if len(opt) > 0 {
		timeOut = int64(opt[0])
	}
	return CallBackMsgpackCacheWithCommander(globalCommander, redisKey, funcs, timeOut)
}

// 获取到内部数据设置缓存，根据里面的数据设置缓存时间
func CallBackMsgpackCacheIn(redisKey string, funcs FuncTypeInt) ([]byte, error) {
	initGlobalCommander()
	return CallBackMsgpackCacheInWithCommander(globalCommander, redisKey, funcs)
}

// 批量删除key
func CommonDelPattern(patternKey string) (err error) {
	initGlobalCommander()
	return globalCommander.DelPattern(patternKey)
}
