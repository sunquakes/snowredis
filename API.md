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
id, err := sf.Generate()
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
id, err := redisSnowflake.Generate()
```

## Core Components

- `RedisSnowflake` - Main ID generator structure
- `Client` - Redis client interface
- `NewBuilder` - Create snowflake algorithm builder instance
- `SetRedisClient` - Set Redis client (Builder method)
- `SetDatacenterID` - Set datacenter ID (Builder method)
- `SetWorkerID` - Set worker ID (Builder method)
- `SetStrictMode` - Set strict mode to use Redis assistance for preventing duplicates (Builder method)
- `Build` - Build final snowflake algorithm instance (Builder method)
- `Generate` - Generate a unique ID based on configuration (local or with Redis assistance)
- `Cleanup` - Perform cleanup operations for the RedisSnowflake instance