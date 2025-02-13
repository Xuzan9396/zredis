package zredis_test

import (
	"github.com/Xuzan9396/zredis"
	"github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

const GIFT_DRAW_POOL = `
local amount_to_deduct = tonumber(ARGV[1]) -- 扣除的金额

local coin = tonumber(redis.call('HGET', KEYS[1], 'coin') or '0')
local coin_free = tonumber(redis.call('HGET', KEYS[1], 'coin_free') or '0')
if coin + coin_free >= amount_to_deduct then
    if coin_free >= amount_to_deduct then
        coin_free = coin_free - amount_to_deduct
        redis.call('HSET', KEYS[1], 'coin_free', coin_free)
        return { 0, coin, amount_to_deduct, coin_free}
    end
    local coin_free_deducted = math.min(coin_free, amount_to_deduct)
    local coin_deducted = amount_to_deduct - coin_free_deducted
    local new_coin_free = 0 -- coin_free 已经全部使用
    local new_coin = coin - coin_deducted
    redis.call('HSET', KEYS[1], 'coin_free', new_coin_free)
    redis.call('HSET', KEYS[1], 'coin', new_coin)
  	return { coin_deducted, new_coin, coin_free_deducted, new_coin_free}

else
    return {}
end
`

func TestRedisPool_CommonCmd(t *testing.T) {
	conn := "127.0.0.1:6379"
	passwd := "27252725"
	dbnum := 0
	zredis.Conn(conn, passwd, dbnum, zredis.WithIdleTime(60*time.Second))
	zredis.CommonHset("test_lua_hset", "coin", 100)
	zredis.CommonHset("test_lua_hset", "coin_free", 10)
	resBytes, err := redis.Values(zredis.CommonLuaScript(GIFT_DRAW_POOL, "test_lua_hset", 2))
	zredis.CommonEXPIRE("test_lua_hset", 3600)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resBytes, len(resBytes))
}
