package zredis

import (
	"github.com/garyburd/redigo/redis"
	"testing"
)

/*
  # 运行所有测试（单元测试 + 集成测试）
  go test -v

  # 只运行单元测试（不需要Redis）
  go test -v -run TestRedisCommands -skip Integration

  # 只运行集成测试（需要Redis服务器）
  go test -v -run TestRedisCommands_Integration

  # 运行特定的集成测试子项
  go test -v -run TestRedisCommands_Integration/BasicOperations

*/
// 模拟的Redis命令执行器用于测试
type mockExecutor struct {
	commands []string
	results  map[string]interface{}
}

func (m *mockExecutor) execute(cmdStr string, keysAndArgs ...interface{}) (interface{}, error) {
	m.commands = append(m.commands, cmdStr)
	if result, exists := m.results[cmdStr]; exists {
		return result, nil
	}
	// 默认返回值
	switch cmdStr {
	case "EXISTS", "SISMEMBER", "HEXISTS":
		return int64(1), nil
	case "SADD", "HSET", "EXPIRE", "DEL", "HDEL":
		return int64(1), nil
	case "ZADD":
		return int64(1), nil
	case "INCR", "INCRBY", "HINCRBY", "DECRBY":
		return int64(10), nil
	case "INCRBYFLOAT":
		return "10.5", nil
	case "ZCARD", "SCARD", "LLEN":
		return int64(3), nil
	case "ZSCORE":
		return "100.0", nil
	case "BITCOUNT":
		return int64(5), nil
	case "SETBIT", "GETBIT":
		return int64(1), nil
	case "SCAN":
		// DelPattern 需要 SCAN 命令返回特定格式
		return []interface{}{"0", []interface{}{}}, nil
	default:
		return "OK", nil
	}
}

func (m *mockExecutor) luaExecute(script string, key string, args ...interface{}) (interface{}, error) {
	m.commands = append(m.commands, "LUA")
	return []interface{}{1, 2, 3}, nil
}

func newMockCommander() RedisCommander {
	mock := &mockExecutor{
		commands: []string{},
		results:  make(map[string]interface{}),
	}
	return NewRedisCommands(mock.execute, mock.luaExecute)
}

// Basic Commands Tests
func TestRedisCommands_Get(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Get("test:key1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_Set(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Set("test:key1", "value1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_SetEx(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.SetEx("test:key1", "value1", 3600)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_SetNx(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.SetNx("test:key1", "value1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_SetNxEx(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.SetNxEx("test:key1", "value1", 3600)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_Del(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Del("test:key1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(1) {
		t.Errorf("Expected 1, got %v", result)
	}
}

func TestRedisCommands_Exists(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Exists("test:key1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(1) {
		t.Errorf("Expected 1, got %v", result)
	}
}

func TestRedisCommands_Expire(t *testing.T) {
	commander := newMockCommander()
	err := commander.Expire("test:key1", 3600)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRedisCommands_ExpireAt(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ExpireAt("test:key1", 1672531200)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_Keys(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Keys("test:*")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

// Counter Commands Tests
func TestRedisCommands_IncrBy(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.IncrBy("test:counter")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(10) {
		t.Errorf("Expected 10, got %v", result)
	}
}

func TestRedisCommands_IncrbyVal(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.IncrbyVal("test:counter", 5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(10) {
		t.Errorf("Expected 10, got %v", result)
	}
}

func TestRedisCommands_IncrbyFloat(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.IncrbyFloat("test:counter", 5.5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "10.5" {
		t.Errorf("Expected '10.5', got %v", result)
	}
}

func TestRedisCommands_DecrByNum(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.DecrByNum("test:counter", 3)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(10) {
		t.Errorf("Expected 10, got %v", result)
	}
}

// Hash Commands Tests
func TestRedisCommands_Hset(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Hset("test:hash", "field1", "value1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(1) {
		t.Errorf("Expected 1, got %v", result)
	}
}

func TestRedisCommands_Hget(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Hget("test:hash", "field1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_HgetAll(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.HgetAll("test:hash")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_Hdel(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Hdel("test:hash", "field1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(1) {
		t.Errorf("Expected 1, got %v", result)
	}
}

func TestRedisCommands_Hexists(t *testing.T) {
	commander := newMockCommander()
	result := commander.Hexists("test:hash", "field1")
	if !result {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestRedisCommands_HIncrby(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.HIncrby("test:hash", "field1", 5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(10) {
		t.Errorf("Expected 10, got %v", result)
	}
}

func TestRedisCommands_HMget(t *testing.T) {
	commander := newMockCommander()
	fields := []interface{}{"field1", "field2"}
	result, err := commander.HMget("test:hash", fields)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

// Set Commands Tests
func TestRedisCommands_SAdd(t *testing.T) {
	commander := newMockCommander()
	err := commander.SAdd("test:set", "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRedisCommands_SRem(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.SRem("test:set", "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_SCard(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.SCard("test:set")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(3) {
		t.Errorf("Expected 3, got %v", result)
	}
}

func TestRedisCommands_SIsMember(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.SIsMember("test:set", "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !result {
		t.Errorf("Expected true, got %v", result)
	}
}

func TestRedisCommands_SMembers(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.SMembers("test:set")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

// ZSet Commands Tests
func TestRedisCommands_ZAdd(t *testing.T) {
	commander := newMockCommander()
	err := commander.ZAdd("test:zset", 100, "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestRedisCommands_ZAddBool(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZAddBool("test:zset", 100, "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1, got %v", result)
	}
}

func TestRedisCommands_ZRem(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZRem("test:zset", "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_ZCard(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZCard("test:zset")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(3) {
		t.Errorf("Expected 3, got %v", result)
	}
}

func TestRedisCommands_ZScore(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZScore("test:zset", "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "100.0" {
		t.Errorf("Expected '100.0', got %v", result)
	}
}

func TestRedisCommands_ZRange(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZRange("test:zset", 0, -1, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_ZRevRange(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZRevRange("test:zset", 0, -1, true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_ZRangeByScore(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZRangeByScore("test:zset", 0, 100, false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_ZRevRank(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZRevRank("test:zset", "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_ZIncrBy(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZIncrBy("test:zset", 10, "member1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_ZIncrByExpire(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.ZIncrByExpire("test:zset", 10, "member1", 3600)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

// List Commands Tests
func TestRedisCommands_LPush(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.LPush("test:list", "value1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_RPop(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.RPop("test:list")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_BRPop(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.BRPop("test:list", 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

func TestRedisCommands_LLen(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.LLen("test:list")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(3) {
		t.Errorf("Expected 3, got %v", result)
	}
}

// Bit Commands Tests
func TestRedisCommands_SetBit(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.SetBit("test:bit", 0, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(1) {
		t.Errorf("Expected 1, got %v", result)
	}
}

func TestRedisCommands_GetBit(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.GetBit("test:bit", 0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(1) {
		t.Errorf("Expected 1, got %v", result)
	}
}

func TestRedisCommands_BitCount(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.BitCount("test:bit")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != int64(5) {
		t.Errorf("Expected 5, got %v", result)
	}
}

// Pattern Delete Test
func TestRedisCommands_DelPattern(t *testing.T) {
	commander := newMockCommander()
	err := commander.DelPattern("test:*")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// Lua Script Test
func TestRedisCommands_LuaScript(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.LuaScript("return KEYS[1]", "test:lua", 1, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result == nil {
		t.Errorf("Expected non-nil result")
	}
}

// Core Cmd Test
func TestRedisCommands_Cmd(t *testing.T) {
	commander := newMockCommander()
	result, err := commander.Cmd("PING")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "OK" {
		t.Errorf("Expected 'OK', got %v", result)
	}
}

// Integration Tests with Real Redis Connection
func TestRedisCommands_Integration(t *testing.T) {
	// 检查Redis连接是否可用
	if !isRedisAvailable() {
		t.Skip("Redis server not available, skipping integration tests")
		return
	}

	// 初始化Redis连接
	Conn("127.0.0.1:6379", "27252725", 0)
	defer cleanupTestData()

	// 获取全局命令实例
	initGlobalCommander()
	commander := globalCommander

	t.Run("BasicOperations", func(t *testing.T) {
		testBasicOperations(t, commander)
	})

	t.Run("HashOperations", func(t *testing.T) {
		testHashOperations(t, commander)
	})

	t.Run("SetOperations", func(t *testing.T) {
		testSetOperations(t, commander)
	})

	t.Run("ZSetOperations", func(t *testing.T) {
		testZSetOperations(t, commander)
	})

	t.Run("ListOperations", func(t *testing.T) {
		testListOperations(t, commander)
	})

	t.Run("BitOperations", func(t *testing.T) {
		testBitOperations(t, commander)
	})

	t.Run("CounterOperations", func(t *testing.T) {
		testCounterOperations(t, commander)
	})

	t.Run("LuaScript", func(t *testing.T) {
		testLuaScript(t, commander)
	})
}

// 检查Redis是否可用
func isRedisAvailable() bool {
	defer func() {
		if r := recover(); r != nil {
			// Redis连接失败
		}
	}()

	// 尝试建立连接
	Conn("127.0.0.1:6379", "27252725", 0)

	// 测试ping命令
	result, err := CommonCmd("PING")
	if err != nil {
		return false
	}

	return result == "PONG"
}

// 清理测试数据
func cleanupTestData() {
	if globalCommander != nil {
		globalCommander.DelPattern("test:*")
	}
}

// 测试基础操作
func testBasicOperations(t *testing.T, commander RedisCommander) {
	// Set和Get
	_, err := commander.Set("test:basic:key1", "value1")
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	result, err := commander.Get("test:basic:key1")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	resultStr, _ := redis.String(result, nil)
	if resultStr != "value1" {
		t.Errorf("Expected 'value1', got %v", resultStr)
	}

	// SetEx
	_, err = commander.SetEx("test:basic:key2", "value2", 3600)
	if err != nil {
		t.Errorf("SetEx failed: %v", err)
	}

	// SetNx
	_, err = commander.SetNx("test:basic:key3", "value3")
	if err != nil {
		t.Errorf("SetNx failed: %v", err)
	}

	// Exists
	exists, err := commander.Exists("test:basic:key1")
	if err != nil {
		t.Errorf("Exists failed: %v", err)
	}
	if exists == 0 {
		t.Errorf("Expected key to exist")
	}

	// Del
	deleted, err := commander.Del("test:basic:key1")
	if err != nil {
		t.Errorf("Del failed: %v", err)
	}
	if deleted == 0 {
		t.Errorf("Expected key to be deleted")
	}
}

// 测试Hash操作
func testHashOperations(t *testing.T, commander RedisCommander) {
	hashKey := "test:hash:user"

	// Hset
	_, err := commander.Hset(hashKey, "name", "John")
	if err != nil {
		t.Errorf("Hset failed: %v", err)
	}

	_, err = commander.Hset(hashKey, "age", "30")
	if err != nil {
		t.Errorf("Hset failed: %v", err)
	}

	// Hget
	name, err := commander.Hget(hashKey, "name")
	if err != nil {
		t.Errorf("Hget failed: %v", err)
	}
	nameStr, _ := redis.String(name, nil)
	if nameStr != "John" {
		t.Errorf("Expected 'John', got %v", nameStr)
	}

	// Hexists
	exists := commander.Hexists(hashKey, "name")
	if !exists {
		t.Errorf("Expected field to exist")
	}

	// HIncrby
	_, err = commander.HIncrby(hashKey, "score", 100)
	if err != nil {
		t.Errorf("HIncrby failed: %v", err)
	}

	// HgetAll
	_, err = commander.HgetAll(hashKey)
	if err != nil {
		t.Errorf("HgetAll failed: %v", err)
	}
}

// 测试Set操作
func testSetOperations(t *testing.T, commander RedisCommander) {
	setKey := "test:set:tags"

	// SAdd
	err := commander.SAdd(setKey, "golang")
	if err != nil {
		t.Errorf("SAdd failed: %v", err)
	}

	err = commander.SAdd(setKey, "redis")
	if err != nil {
		t.Errorf("SAdd failed: %v", err)
	}

	// SIsMember
	isMember, err := commander.SIsMember(setKey, "golang")
	if err != nil {
		t.Errorf("SIsMember failed: %v", err)
	}
	if !isMember {
		t.Errorf("Expected 'golang' to be a member")
	}

	// SCard
	count, err := commander.SCard(setKey)
	if err != nil {
		t.Errorf("SCard failed: %v", err)
	}
	if count.(int64) != 2 {
		t.Errorf("Expected 2 members, got %v", count)
	}

	// SMembers
	_, err = commander.SMembers(setKey)
	if err != nil {
		t.Errorf("SMembers failed: %v", err)
	}
}

// 测试ZSet操作
func testZSetOperations(t *testing.T, commander RedisCommander) {
	zsetKey := "test:zset:scores"

	// ZAdd
	err := commander.ZAdd(zsetKey, 100, "player1")
	if err != nil {
		t.Errorf("ZAdd failed: %v", err)
	}

	err = commander.ZAdd(zsetKey, 85, "player2")
	if err != nil {
		t.Errorf("ZAdd failed: %v", err)
	}

	// ZScore
	score, err := commander.ZScore(zsetKey, "player1")
	if err != nil {
		t.Errorf("ZScore failed: %v", err)
	}
	scoreStr, _ := redis.String(score, nil)
	if scoreStr != "100" {
		t.Errorf("Expected score '100', got %v", scoreStr)
	}

	// ZCard
	count, err := commander.ZCard(zsetKey)
	if err != nil {
		t.Errorf("ZCard failed: %v", err)
	}
	if count.(int64) != 2 {
		t.Errorf("Expected 2 members, got %v", count)
	}

	// ZRange
	_, err = commander.ZRange(zsetKey, 0, -1, false)
	if err != nil {
		t.Errorf("ZRange failed: %v", err)
	}

	// ZRevRange
	_, err = commander.ZRevRange(zsetKey, 0, -1, true)
	if err != nil {
		t.Errorf("ZRevRange failed: %v", err)
	}
}

// 测试List操作
func testListOperations(t *testing.T, commander RedisCommander) {
	listKey := "test:list:messages"

	// LPush
	_, err := commander.LPush(listKey, "message1")
	if err != nil {
		t.Errorf("LPush failed: %v", err)
	}

	_, err = commander.LPush(listKey, "message2")
	if err != nil {
		t.Errorf("LPush failed: %v", err)
	}

	// LLen
	length, err := commander.LLen(listKey)
	if err != nil {
		t.Errorf("LLen failed: %v", err)
	}
	if length.(int64) != 2 {
		t.Errorf("Expected length 2, got %v", length)
	}

	// RPop
	_, err = commander.RPop(listKey)
	if err != nil {
		t.Errorf("RPop failed: %v", err)
	}
}

// 测试Bit操作
func testBitOperations(t *testing.T, commander RedisCommander) {
	bitKey := "test:bit:flags"

	// SetBit
	_, err := commander.SetBit(bitKey, 0, 1)
	if err != nil {
		t.Errorf("SetBit failed: %v", err)
	}

	_, err = commander.SetBit(bitKey, 7, 1)
	if err != nil {
		t.Errorf("SetBit failed: %v", err)
	}

	// GetBit
	bit, err := commander.GetBit(bitKey, 0)
	if err != nil {
		t.Errorf("GetBit failed: %v", err)
	}
	if bit.(int64) != 1 {
		t.Errorf("Expected bit value 1, got %v", bit)
	}

	// BitCount
	count, err := commander.BitCount(bitKey)
	if err != nil {
		t.Errorf("BitCount failed: %v", err)
	}
	if count.(int64) != 2 {
		t.Errorf("Expected bit count 2, got %v", count)
	}
}

// 测试计数器操作
func testCounterOperations(t *testing.T, commander RedisCommander) {
	counterKey := "test:counter:visits"

	// IncrBy
	_, err := commander.IncrBy(counterKey)
	if err != nil {
		t.Errorf("IncrBy failed: %v", err)
	}

	// IncrbyVal
	result, err := commander.IncrbyVal(counterKey, 5)
	if err != nil {
		t.Errorf("IncrbyVal failed: %v", err)
	}
	if result.(int64) != 6 { // 1 + 5
		t.Errorf("Expected counter value 6, got %v", result)
	}

	// IncrbyFloat
	floatKey := "test:counter:float"
	_, err = commander.IncrbyFloat(floatKey, 3.14)
	if err != nil {
		t.Errorf("IncrbyFloat failed: %v", err)
	}
}

// 测试Lua脚本
func testLuaScript(t *testing.T, commander RedisCommander) {
	script := `
		local key = KEYS[1]
		local value = ARGV[1]
		redis.call('SET', key, value)
		return redis.call('GET', key)
	`

	result, err := commander.LuaScript(script, "test:lua:key", "test_value")
	if err != nil {
		t.Errorf("LuaScript failed: %v", err)
	}

	resultStr, _ := redis.String(result, nil)
	if resultStr != "test_value" {
		t.Errorf("Expected 'test_value', got %v", resultStr)
	}
}
