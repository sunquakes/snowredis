# Snowflake-Redis Distributed ID Generation Library

This is a distributed ID generation library based on the Snowflake algorithm with Redis coordination, capable of generating globally unique 64-bit IDs.

## Features

- Generate unique IDs based on the Snowflake algorithm
- Use Redis for node coordination
- Support multi-node deployment
- Pluggable Redis client implementation

## Usage

### Basic Usage

```go
// Load configuration
cfg := config.LoadConfig()

// Initialize Redis connection
redis.InitRedis(cfg)

// Create snowflake algorithm instance
// Parameter description:
// - cfg.NodeID: Unique identifier of the node in Redis
// - 1: Datacenter ID (for physical partitioning of the ID structure)
// - 1: Worker ID (for physical partitioning of the ID structure)
sf, err := snowflake.NewRedisSnowflake(cfg.NodeID, 1, 1)
if err != nil {
    log.Fatalf("Failed to initialize Redis Snowflake: %v", err)
}
defer sf.Cleanup()  // Clean up resources

// Generate unique ID
id, err := sf.GenerateUsingMutex()
if err != nil {
    log.Printf("Error generating ID: %v", err)
}
fmt.Printf("Generated ID: %d\n", id)
```

### Parameter Description

- **Node ID**: Unique identifier registered in Redis, used to prevent multiple nodes with the same ID from running simultaneously
- **Datacenter ID**: Datacenter identifier for physical partitioning of the ID structure (maintaining compatibility with the standard Snowflake algorithm)
- **Worker ID**: Worker identifier for physical partitioning of the ID structure (maintaining compatibility with the standard Snowflake algorithm)

> Note: In Redis coordination mode, Node ID is the key uniqueness guarantee, while Datacenter ID and Worker ID are mainly used to maintain ID structure compatibility and provide local uniqueness when Redis is unavailable.

### Custom Redis Client

```go
// Create custom Redis client
mockRedis := mock.NewMockRedisClient()

// Create snowflake algorithm instance with custom Redis client
redisSnowflake, err := snowflake.NewRedisSnowflakeWithClient(mockRedis, 1, 1, 1)
if err != nil {
    log.Fatalf("Failed to initialize: %v", err)
}
defer redisSnowflake.Cleanup()

// Generate ID
id, err := redisSnowflake.GenerateUsingMutex()
```

## Core Components

- `RedisSnowflake` - Main ID generator structure
- `Client` - Redis client interface
- `NewSnowflake` - Instance using Redis to automatically assign datacenterID and workerID (deprecated, use NewSnowflakeWithClient)
- `NewSnowflakeWithClient` - Instance using Redis to automatically assign datacenterID and workerID (provides Redis client)
- `NewSnowflakeWithConfig` - Create instance with default Redis client (deprecated, use NewSnowflakeWithConfigAndClient)
- `NewSnowflakeWithConfigAndClient` - Create instance with custom Redis client (requires manual specification of datacenterID and workerID)
- `NewRedisSnowflake` - Instance using Redis to automatically assign datacenterID and workerID (deprecated, use NewRedisSnowflakeWithClientAlias)
- `NewRedisSnowflakeWithClient` - Create instance with custom Redis client (requires manual specification of datacenterID and workerID)

- `NewRedisAllIDSnowflake` - Instance assigning datacenterID and workerID through Redis (deprecated, use NewRedisAllIDSnowflakeWithClient)
- `NewRedisIDSnowflake` - Instance assigning datacenterID and workerID through Redis (deprecated, use NewRedisIDSnowflakeWithClient)
- `NewAutoSnowflake` - Instance with automatic node ID assignment (deprecated, use NewAutoSnowflakeWithClient)
- `NewAutoRedisSnowflake` - Instance with automatic node ID assignment (deprecated, use NewAutoRedisSnowflakeWithClient)
- `NewInitializedSnowflake` - Instance that only uses Redis during initialization

- `NewBuilder` - Create snowflake algorithm builder instance
- `SetRedisClient` - Set Redis client (Builder method)
- `SetDatacenterID` - Set datacenter ID (Builder method)
- `SetWorkerID` - Set worker ID (Builder method)
- `Build` - Build final snowflake algorithm instance (Builder method)
- `NewRedisAllIDSnowflakeWithClient` - Instance assigning datacenterID and workerID through custom Redis client
- `NewRedisIDSnowflakeWithClient` - Instance assigning datacenterID and workerID through custom Redis client
- `NewAutoSnowflakeWithClient` - 通过自定义Redis客户端自动分配节点ID的实例


- `NewSnowflakeWithConfigAndClient` - 通过自定义Redis客户端和指定datacenterID/workerID创建实例
- `Generate` - 生成ID（完全本地化，无Redis访问）
- `GenerateUsingMutex` - 兼容性方法（实际也完全本地化）
- `Cleanup` - 清理资源（清理节点注册信息）
- `Del` - 删除键
- `GetRedisInstance` - 获取全局Redis实例
- `NewRedisClient` - 创建并返回Redis客户端实例
- `SetNX` - 如果键不存在则设置（实现Client接口）
- `Incr` - 原子递增（实现Client接口）
- `Del` - 删除键（实现Client接口）