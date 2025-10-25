# ZRedis - 统一的Redis客户端库

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](#测试)

ZRedis 是一个功能强大、易于使用的 Go Redis 客户端库，提供了三种使用模式：全局单例、多实例管理和单实例操作。经过重构后，所有模式都使用统一的命令接口，确保代码的一致性和可维护性。

## 🚀 特性

- **统一接口设计** - 所有Redis命令通过统一的 `RedisCommander` 接口提供
- **三种使用模式** - 支持全局、多实例、单实例三种不同的使用场景
- **完全向后兼容** - 重构后保持所有原有API不变
- **连接池管理** - 内置高效的连接池管理机制
- **TLS支持** - 支持安全的TLS连接
- **Lua脚本支持** - 完整的Redis Lua脚本执行支持
- **丰富的数据类型** - 支持String、Hash、Set、ZSet、List、Bit等所有Redis数据类型
- **缓存辅助函数** - 内置缓存回调函数，简化缓存使用

## 📦 安装

```bash
go get github.com/Xuzan9396/zredis
```

## 🏗️ 架构概览

```
zredis/
├── commands.go          # 统一的Redis命令接口和实现
├── command.go           # 全局模式的包装函数
├── conn.go              # 全局连接池管理
├── mredis/             # 多实例模式
│   ├── command.go      # 多实例命令包装
│   └── conn.go         # 多实例连接池管理
└── sredis/             # 单实例模式
    ├── command.go      # 单实例命令方法
    └── conn.go         # 单实例连接池
```

## 📖 使用方法

### 1. 全局模式 (推荐用于简单应用)

全局模式适用于只需要连接一个Redis实例的应用。

```go
package main

import (
    "fmt"
    "github.com/Xuzan9396/zredis"
    "time"
)

func main() {
    // 初始化全局连接
    zredis.Conn("127.0.0.1:6379", "password", 0, 
        zredis.WithMaxActive(200),
        zredis.WithMaxIdle(100),
        zredis.WithIdleTime(300*time.Second))
    
    // 基础操作
    zredis.CommonSet("user:1", "john")
    value, _ := zredis.CommonGet("user:1")
    fmt.Println("User:", value)
    
    // 设置过期时间
    zredis.CommonSetEx("session:abc", "user_data", 3600)
    zredis.CommonExpire("user:1", 7200)
    
    // Hash操作
    zredis.CommonHset("user:1:profile", "name", "John Doe")
    zredis.CommonHset("user:1:profile", "age", "30")
    profile, _ := zredis.CommonHgetAll("user:1:profile")
    fmt.Println("Profile:", profile)
    
    // Set操作
    zredis.CommonSADD("tags", "golang")
    zredis.CommonSADD("tags", "redis")
    exists, _ := zredis.CommonSISMEMBER("tags", "golang")
    fmt.Println("Tag exists:", exists)
    
    // ZSet操作
    zredis.CommonZADD("leaderboard", 100, "player1")
    zredis.CommonZADD("leaderboard", 85, "player2")
    topPlayers, _ := zredis.CommonZRevRange("leaderboard", 0, 9, true)
    fmt.Println("Top players:", topPlayers)
    
    // List操作
    zredis.CommonLPush("messages", "Hello World")
    message, _ := zredis.CommonRPop("messages")
    fmt.Println("Message:", message)
    
    // 计数器操作
    zredis.CommonIncrBy("visitors")
    count, _ := zredis.CommonIncrby("page_views", 5)
    fmt.Println("Page views:", count)
    
    // Bit操作
    zredis.CommonSetBit("user_online", 123, 1)
    online, _ := zredis.CommonGetBit("user_online", 123)
    fmt.Println("User online:", online)
    
    // Lua脚本执行
    script := `
        local current = redis.call('GET', KEYS[1])
        if current == false then
            redis.call('SET', KEYS[1], ARGV[1])
            return ARGV[1]
        end
        return current
    `
    result, _ := zredis.CommonLuaScript(script, "atomic_key", "new_value")
    fmt.Println("Lua result:", result)
}
```

### 2. 多实例模式 (推荐用于多Redis环境)

多实例模式适用于需要连接多个Redis实例的应用，如读写分离、不同业务使用不同实例等场景。

```go
package main

import (
    "fmt"
    "github.com/Xuzan9396/zredis/mredis"
    "time"
)

func main() {
    // 配置多个Redis实例
    // 主库 - 用于写操作
    mredis.Conn("master", "127.0.0.1:6379", "password", 0,
        mredis.WithMaxActive(300),
        mredis.WithMaxIdle(150))
    
    // 从库 - 用于读操作  
    mredis.Conn("slave", "127.0.0.1:6380", "password", 0,
        mredis.WithMaxActive(200),
        mredis.WithMaxIdle(100))
    
    // 缓存库 - 用于临时数据
    mredis.Conn("cache", "127.0.0.1:6381", "", 1,
        mredis.WithIdleTime(60*time.Second))
    
    // 写入主库
    mredis.CommonSet("master", "user:1001", "alice")
    mredis.CommonHset("master", "user:1001:profile", "email", "alice@example.com")
    mredis.CommonZADD("master", "user_scores", 95, "user:1001")
    
    // 从从库读取
    userData, _ := mredis.CommonGet("slave", "user:1001")
    userEmail, _ := mredis.CommonHget("slave", "user:1001:profile", "email")
    userScore, _ := mredis.CommonZscore("slave", "user_scores", "user:1001")
    
    fmt.Printf("User: %s, Email: %s, Score: %s\n", userData, userEmail, userScore)
    
    // 使用缓存库存储临时数据
    mredis.CommonSetEx("cache", "temp:session:xyz", "temp_data", 300)
    
    // 获取统一命令实例进行更复杂操作
    masterCommander := mredis.GetCommander("master")
    slaveCommander := mredis.GetCommander("slave")
    
    // 使用统一接口
    masterCommander.Set("new_key", "new_value")
    value, _ := slaveCommander.Get("new_key")
    fmt.Println("New value:", value)
    
    // 批量删除模式匹配的key
    mredis.CommonDelPattern("cache", "temp:*")
    
    // 缓存回调函数使用
    cachedData, _ := mredis.CallBackMsgpackCache("cache", "expensive_operation", func() ([]byte, error) {
        // 模拟耗时操作
        return []byte("computed result"), nil
    }, 3600) // 缓存1小时
    
    fmt.Println("Cached data:", string(cachedData))
}
```

### 3. 单实例模式 (推荐用于面向对象设计)

单实例模式适用于需要将Redis连接作为对象管理的场景，提供更好的封装性。

```go
package main

import (
    "fmt"
    "github.com/Xuzan9396/zredis/sredis"
    "time"
)

func main() {
    // 创建Redis实例
    client := sredis.Conn("127.0.0.1:6379", "password", 0,
        sredis.WithMaxActive(100),
        sredis.WithMaxIdle(50),
        sredis.WithIdleTime(300*time.Second))
    
    // 如果需要TLS连接
    // client := sredis.Conn("redis.example.com:6380", "password", 0,
    //     sredis.WithRedisTLS())
    
    // 基础操作
    client.CommonSet("product:1", "iPhone 14")
    product, _ := client.CommonGet("product:1")
    fmt.Println("Product:", product)
    
    // 商品信息存储 (Hash)
    client.CommonHset("product:1:details", "name", "iPhone 14")
    client.CommonHset("product:1:details", "price", "999")
    client.CommonHset("product:1:details", "category", "smartphone")
    
    details, _ := client.CommonHgetAll("product:1:details")
    fmt.Println("Product details:", details)
    
    // 商品标签 (Set)
    client.CommonSADD("product:1:tags", "apple")
    client.CommonSADD("product:1:tags", "phone")
    client.CommonSADD("product:1:tags", "5g")
    
    tags, _ := client.CommonSMEMBERS("product:1:tags")
    fmt.Println("Product tags:", tags)
    
    // 商品评分排行 (ZSet)
    client.CommonZADD("product_ratings", 4.8, "product:1")
    client.CommonZADD("product_ratings", 4.5, "product:2")
    client.CommonZADD("product_ratings", 4.9, "product:3")
    
    topRated, _ := client.CommonZRevRange("product_ratings", 0, 2, true)
    fmt.Println("Top rated products:", topRated)
    
    // 最近查看的商品 (List)
    client.CommonLPush("user:123:recent_views", "product:1")
    client.CommonLPush("user:123:recent_views", "product:5")
    
    recentView, _ := client.CommonRPop("user:123:recent_views")
    fmt.Println("Recent view:", recentView)
    
    // 用户在线状态 (Bit)
    client.CommonSetBit("users_online", 123, 1)  // 用户123在线
    client.CommonSetBit("users_online", 456, 1)  // 用户456在线
    
    onlineCount, _ := client.CommonBitCount("users_online")
    fmt.Println("Online users count:", onlineCount)
    
    // 页面访问计数
    client.CommonIncrby("page:home:views", 1)
    client.CommonHIncrby("page:stats", "unique_visitors", 1)
    
    views, _ := client.CommonGet("page:home:views")
    fmt.Println("Page views:", views)
    
    // 获取统一命令接口进行更复杂操作
    commander := client.GetCommander()
    
    // 使用统一接口
    commander.SetNxEx("lock:operation", "locked", 30) // 30秒锁
    exists, _ := commander.Exists("lock:operation")
    fmt.Println("Lock exists:", exists)
    
    // 缓存函数使用
    expensiveData, _ := client.CallBackMsgpackCache("computation_result", func() ([]byte, error) {
        // 模拟复杂计算
        time.Sleep(100 * time.Millisecond)
        return []byte("complex computation result"), nil
    }, 1800) // 缓存30分钟
    
    fmt.Println("Expensive computation:", string(expensiveData))
    
    // 带内部时间控制的缓存
    timedData, _ := client.CallBackMsgpackCacheIn("timed_cache", func() ([]byte, int64, error) {
        // 返回数据和过期时间戳
        expireAt := time.Now().Add(2 * time.Hour).Unix()
        return []byte("timed data"), expireAt, nil
    })
    
    fmt.Println("Timed data:", string(timedData))
}
```

## 🔧 配置选项

### 连接池配置

```go
// 最大活跃连接数
zredis.WithMaxActive(200)

// 最大空闲连接数  
zredis.WithMaxIdle(100)

// 空闲连接超时时间
zredis.WithIdleTime(300 * time.Second)

// 自定义连接选项
zredis.WithRedisOption(
    redis.DialConnectTimeout(5*time.Second),
    redis.DialReadTimeout(10*time.Second),
    redis.DialWriteTimeout(10*time.Second),
)

// 启用TLS连接
zredis.WithRedisTLS()
```

### 连接池推荐配置

```go
// 高并发场景 (>1000 QPS)
zredis.WithMaxActive(500)
zredis.WithMaxIdle(200)
zredis.WithIdleTime(300 * time.Second)

// 中等负载场景 (100-1000 QPS)
zredis.WithMaxActive(200)
zredis.WithMaxIdle(100)
zredis.WithIdleTime(300 * time.Second)

// 低负载场景 (<100 QPS)
zredis.WithMaxActive(50)
zredis.WithMaxIdle(25)
zredis.WithIdleTime(300 * time.Second)
```

## 📚 Redis命令支持

### 基础命令
- `Get`, `Set`, `SetEx`, `SetNx`, `SetNxEx`
- `Del`, `Exists`, `Expire`, `ExpireAt`, `Keys`

### 计数器命令
- `IncrBy`, `IncrbyVal`, `IncrbyFloat`, `DecrByNum`

### Hash命令
- `Hset`, `Hget`, `HgetAll`, `Hdel`, `Hexists`
- `HIncrby`, `HMget`

### Set命令
- `SAdd`, `SRem`, `SCard`, `SIsMember`, `SMembers`

### ZSet (有序集合) 命令
- `ZAdd`, `ZAddBool`, `ZRem`, `ZCard`, `ZScore`
- `ZRange`, `ZRevRange`, `ZRangeByScore`, `ZRevRank`
- `ZIncrBy`, `ZIncrByExpire`

### List命令
- `LPush`, `RPop`, `BRPop`, `LLen`

### Bit命令
- `SetBit`, `GetBit`, `BitCount`

### 高级功能
- `LuaScript` - Lua脚本执行
- `DelPattern` - 模式匹配批量删除
- `Cmd` - 原始命令执行

## 🔄 缓存辅助函数

### 基础缓存回调

```go
// 全局模式
data, err := zredis.CallBackMsgpackCache("cache_key", func() ([]byte, error) {
    // 执行耗时操作
    result := fetchDataFromDatabase()
    return json.Marshal(result)
}, 3600) // 缓存1小时

// 多实例模式
data, err := mredis.CallBackMsgpackCache("cache_instance", "cache_key", func() ([]byte, error) {
    return expensiveComputation()
}, 7200) // 缓存2小时

// 单实例模式
data, err := client.CallBackMsgpackCache("cache_key", func() ([]byte, error) {
    return fetchRemoteData()
}, 1800) // 缓存30分钟
```

### 带内部时间控制的缓存

```go
// 根据业务逻辑动态设置过期时间
data, err := client.CallBackMsgpackCacheIn("dynamic_cache", func() ([]byte, int64, error) {
    result := computeResult()
    
    // 根据结果决定缓存时间
    var expireAt int64
    if isImportantData(result) {
        expireAt = time.Now().Add(24 * time.Hour).Unix() // 重要数据缓存24小时
    } else {
        expireAt = time.Now().Add(1 * time.Hour).Unix()  // 普通数据缓存1小时
    }
    
    return result, expireAt, nil
})
```

## 🧪 测试

项目包含完整的测试套件，覆盖所有Redis命令和功能。

```bash
# 运行所有测试
go test ./...

# 运行单元测试 (不需要Redis服务器)
go test -run TestRedisCommands

# 运行集成测试 (需要Redis服务器)
go test -run TestRedisPool_CommonCmd

# 运行特定包的测试
go test ./mredis
go test ./sredis
```

### 测试统计
- **48个单元测试** - 覆盖所有Redis命令接口
- **3个集成测试** - 验证实际Redis连接
- **100%通过率** - 所有测试都成功通过

## 🔍 最佳实践

### 1. 选择合适的使用模式

```go
// 简单应用 - 使用全局模式
import "github.com/Xuzan9396/zredis"

// 多Redis环境 - 使用多实例模式  
import "github.com/Xuzan9396/zredis/mredis"

// 面向对象设计 - 使用单实例模式
import "github.com/Xuzan9396/zredis/sredis"
```

### 2. 连接池管理

```go
// 在应用启动时初始化连接
func init() {
    zredis.Conn("127.0.0.1:6379", os.Getenv("REDIS_PASSWORD"), 0,
        zredis.WithMaxActive(200),
        zredis.WithMaxIdle(100))
}

// 在应用关闭时清理资源 (如果需要)
func cleanup() {
    // 连接池会自动管理连接生命周期
}
```

### 3. 错误处理

```go
value, err := zredis.CommonGet("key")
if err != nil {
    log.Printf("Redis error: %v", err)
    // 处理错误，比如返回默认值或从其他数据源获取
    return defaultValue, nil
}
```

### 4. 使用缓存模式

```go
// 缓存模式 - 先查缓存，未命中则计算并缓存
func GetUserProfile(userID string) (*UserProfile, error) {
    cacheKey := fmt.Sprintf("user:profile:%s", userID)
    
    data, err := zredis.CallBackMsgpackCache(cacheKey, func() ([]byte, error) {
        // 从数据库获取用户信息
        profile, err := db.GetUserProfile(userID)
        if err != nil {
            return nil, err
        }
        return json.Marshal(profile)
    }, 3600) // 缓存1小时
    
    if err != nil {
        return nil, err
    }
    
    var profile UserProfile
    err = json.Unmarshal(data, &profile)
    return &profile, err
}
```

### 5. Lua脚本使用

```go
// 原子性操作 - 使用Lua脚本确保操作的原子性
const incrementWithLimit = `
    local key = KEYS[1]
    local limit = tonumber(ARGV[1])
    local increment = tonumber(ARGV[2])
    
    local current = tonumber(redis.call('GET', key) or '0')
    if current + increment <= limit then
        return redis.call('INCRBY', key, increment)
    else
        return current
    end
`

result, err := zredis.CommonLuaScript(incrementWithLimit, "counter:api_calls", 1000, 1)
```

## 🔄 迁移指南

### 从旧版本迁移

如果你正在使用旧版本的代码，重构后的版本完全向后兼容：

```go
// 旧代码仍然可以正常工作
zredis.CommonGet("key")
mredis.CommonSet("instance", "key", "value")
client.CommonHset("hash", "field", "value")

// 新增功能：可以获取统一的命令接口
commander := mredis.GetCommander("instance")
commander.Get("key") // 使用新的统一接口
```



---

**ZRedis** - 让Redis在Go中的使用更加简单、统一、高效！