# SnowRedis

[中文文档](README_CN.md) | [English Documentation](README.md)

基于Snowflake算法并使用Redis进行协调的分布式ID生成器。

## 特性

- 标准Snowflake算法实现
- Redis协调用于分布式环境
- 基于接口的Redis客户端，支持可插拔实现（适用于任何Redis客户端库：go-redis、redigo、rueidis等）
- 支持自定义Redis客户端
- 构建器模式，灵活配置
- 线程安全的ID生成
- 高性能
- 严格模式配合Redis辅助防重复

## 安装

```bash
go get github.com/sunquakes/snowredis
```

## 使用

### 自动分配模式（推荐）

自动从Redis分配datacenterID和workerID。这是大多数用例的推荐方式：

```go
package main

import (
	"fmt"
	"log"

	"github.com/sunquakes/snowredis/redis"
	"github.com/sunquakes/snowredis/snowflake"
)

func main() {
	// 配置Redis连接
	cfg := &redis.Config{
		Addr: "localhost:6379",  // Redis地址
		Pwd:  "",                // Redis密码（如果有）
		Db:   0,                 // Redis数据库编号
	}
	
	// 创建Redis客户端
	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatalf("创建Redis客户端失败: %v", err)
	}

	// 从Redis自动分配ID
	sf, err := snowflake.NewBuilder().
		SetRedisClient(redisClient).  // 自动分配需要Redis客户端
		Build()
	if err != nil {
		log.Fatalf("自动分配初始化失败: %v", err)
	}
	defer sf.Cleanup()

	// 生成唯一ID
	for i := 0; i < 5; i++ {
		id, err := sf.Generate()
		if err != nil {
			log.Printf("生成ID错误: %v", err)
			continue
		}
		fmt.Printf("生成ID: %d\n", id)
	}
}
```

### 手动配置模式

显式设置数据中心ID和工作ID：

```go
// 使用手动配置创建snowflake算法实例
cfg := &redis.Config{
	Addr: "localhost:6379",  // Redis地址
	Pwd:  "",                // Redis密码（如果有）
	Db:   0,                 // Redis数据库编号
}

redisClient, err := redis.NewClient(cfg)
if err != nil {
	log.Fatalf("创建Redis客户端失败: %v", err)
}

sf, err := snowflake.NewBuilder().
	SetRedisClient(redisClient).
	SetDatacenterID(1).
	SetWorkerID(1).
	Build()
if err != nil {
	log.Fatalf("初始化Redis Snowflake失败: %v", err)
}
defer sf.Cleanup()

// 生成唯一ID
for i := 0; i < 5; i++ {
	id, err := sf.Generate()
	if err != nil {
		log.Printf("生成ID错误: %v", err)
		continue
	}
	fmt.Printf("生成ID: %d\n", id)
}
```

### 使用严格模式增强唯一性

启用严格模式使用Redis辅助防重复：

```go
// 创建启用严格模式的实例
cfg := &redis.Config{
	Addr: "localhost:6379",  // Redis地址
	Pwd:  "",                // Redis密码（如果有）
	Db:   0,                 // Redis数据库编号
}
redisClient, err := redis.NewClient(cfg)
if err != nil {
    log.Fatalf("创建Redis客户端失败: %v", err)
}

sf, err := snowflake.NewBuilder().
    SetRedisClient(redisClient).
    SetDatacenterID(1).
    SetWorkerID(1).
    SetStrictMode(true).  // 启用严格模式以额外防重复
    Build()
if err != nil {
    log.Fatalf("严格模式初始化失败: %v", err)
}
defer sf.Cleanup()

id, err := sf.Generate()  // 严格模式启用时这将使用Redis确保唯一性
}
```

### 默认值模式

使用默认值，无需Redis协调：

```go
// 创建使用默认值的实例（不需要Redis客户端）
sf, err := snowflake.NewBuilder().
    Build()  // 使用默认数据中心ID和工作ID
if err != nil {
    log.Fatalf("默认值初始化失败: %v", err)
}
defer sf.Cleanup()

id, err := sf.Generate()
if err != nil {
    log.Printf("生成ID错误: %v", err)
} else {
    fmt.Printf("生成ID: %d\n", id)
}
```

## API方法

### 构建器方法
- `NewBuilder()` - 创建新的构建器实例
- `SetRedisClient(client)` - 设置Redis客户端
- `SetDatacenterID(id)` - 设置数据中心ID
- `SetWorkerID(id)` - 设置工作ID
- `SetStrictMode(strict)` - 启用/禁用严格模式
- `Build()` - 构建snowflake实例

### 实例方法
- `Generate()` - 生成唯一ID
- `Cleanup()` - 清理资源

## 配置

该库支持三种主要配置方式：

1. **自动分配（推荐）**: 让Redis自动分配ID
2. **手动配置**: 显式设置数据中心ID和工作ID
3. **默认值**: 使用内置默认值

## 性能

该库为高性能而设计：
- 初始化后本地ID生成
- 线程安全操作
- 正常模式下最小的Redis交互
- 严格模式下额外的Redis检查以增强唯一性

## 自定义Redis客户端实现

该库对Redis客户端采用基于接口的方法，允许您实现符合RedisClient接口的自己的Redis客户端：

```go
type RedisClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Del(ctx context.Context, keys ...string) (int64, error)
}
```

您可以使用任何您选择的Redis客户端库实现此接口（例如go-redis、redigo、rueidis、goredis/redismock用于测试等），并将您的自定义实现实例传递给构建器：

```go
type MyCustomRedisClient struct {
	// 您的Redis客户端实现，使用您首选的库
}

func (c *MyCustomRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	// 使用您首选的Redis客户端库实现
}

func (c *MyCustomRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	// 使用您首选的Redis客户端库实现
}

func (c *MyCustomRedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	// 使用您首选的Redis客户端库实现
}

// 用法
customRedisClient := &MyCustomRedisClient{}
sf, err := snowflake.NewBuilder().
    SetRedisClient(customRedisClient).
    SetDatacenterID(1).
    SetWorkerID(1).
    Build()
```

这种方法提供了根据您的性能、功能或依赖需求使用不同Redis客户端库的灵活性。默认实现在使用go-redis，但您不限于此。

## 注意事项

- 在正常模式下，初始化后ID生成完全是本地的。Redis仅在设置期间用于协调唯一标识符。
- 在严格模式下，Redis在ID生成期间还会被使用，以提供额外的防重复层。
- 最小内存占用