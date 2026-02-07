# SnowRedis

[中文文档](README_CN.md) | [English Documentation](README.md)

A distributed ID generator utility using the Snowflake algorithm with Redis for coordination.

## Features

- Standard Snowflake algorithm implementation
- Redis coordination for distributed environments
- Interface-based Redis client for pluggable implementations (works with any Redis client library: go-redis, redigo, rueidis, etc.)
- Support for custom Redis clients
- Builder pattern for flexible configuration
- Thread-safe ID generation
- High performance
- Strict mode with Redis assistance to prevent duplicates

## Installation

```bash
go get github.com/sunquakes/snowredis
```

## Usage

### Auto-allocation Mode (Recommended)

Automatically allocate datacenterID and workerID from Redis. This is the recommended approach for most use cases:

```go
package main

import (
	"fmt"
	"log"

	"github.com/sunquakes/snowredis/redis"
	"github.com/sunquakes/snowredis/snowflake"
)

func main() {
	// Configure Redis connection
	cfg := &redis.Config{
		Addr: "localhost:6379",  // Redis address
		Pwd:  "",                // Redis password (if any)
		Db:   0,                 // Redis database number
	}
	
	// Create a Redis client
	redisClient, err := redis.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Redis client: %v", err)
	}

	// Automatically allocate IDs from Redis
	sf, err := snowflake.NewBuilder().
		SetRedisClient(redisClient).  // Redis client required for auto-allocation
		Build()
	if err != nil {
		log.Fatalf("Failed to initialize with auto-allocation: %v", err)
	}
	defer sf.Cleanup()

	// Generate unique IDs
	for i := 0; i < 5; i++ {
		id, err := sf.Generate()
		if err != nil {
			log.Printf("Error generating ID: %v", err)
			continue
		}
		fmt.Printf("Generated ID: %d\n", id)
	}
}
```

### Manual Configuration Mode

Explicitly set datacenter ID and worker ID:

```go
// Create snowflake algorithm instance with manual configuration
cfg := &redis.Config{
	Addr: "localhost:6379",  // Redis address
	Pwd:  "",                // Redis password (if any)
	Db:   0,                 // Redis database number
}

redisClient, err := redis.NewClient(cfg)
if err != nil {
	log.Fatalf("Failed to create Redis client: %v", err)
}

sf, err := snowflake.NewBuilder().
	SetRedisClient(redisClient).
	SetDatacenterID(1).
	SetWorkerID(1).
	Build()
if err != nil {
	log.Fatalf("Failed to initialize Redis Snowflake: %v", err)
}
defer sf.Cleanup()

// Generate unique IDs
for i := 0; i < 5; i++ {
	id, err := sf.Generate()
	if err != nil {
		log.Printf("Error generating ID: %v", err)
		continue
	}
	fmt.Printf("Generated ID: %d\n", id)
}
```

### Using Strict Mode for Enhanced Uniqueness

Enable strict mode to use Redis assistance for preventing duplicates:

```go
// Create instance with strict mode enabled
cfg := &redis.Config{
	Addr: "localhost:6379",  // Redis address
	Pwd:  "",                // Redis password (if any)
	Db:   0,                 // Redis database number
}
redisClient, err := redis.NewClient(cfg)
if err != nil {
    log.Fatalf("Failed to create Redis client: %v", err)
}

sf, err := snowflake.NewBuilder().
    SetRedisClient(redisClient).
    SetDatacenterID(1).
    SetWorkerID(1).
    SetStrictMode(true).  // Enable strict mode for additional duplicate prevention
    Build()
if err != nil {
    log.Fatalf("Failed to initialize in strict mode: %v", err)
}
defer sf.Cleanup()

id, err := sf.Generate()  // This will use Redis to ensure uniqueness when strict mode is enabled
}
```

### Default Values Mode

Use default values without Redis coordination:

```go
// Create instance with default values (no Redis client needed)
sf, err := snowflake.NewBuilder().
    Build()  // Uses default datacenter ID and worker ID
if err != nil {
    log.Fatalf("Failed to initialize with defaults: %v", err)
}
defer sf.Cleanup()

id, err := sf.Generate()
if err != nil {
    log.Printf("Error generating ID: %v", err)
} else {
    fmt.Printf("Generated ID: %d\n", id)
}
```

## API Methods

### Builder Methods
- `NewBuilder()` - Creates a new builder instance
- `SetRedisClient(client)` - Sets the Redis client
- `SetDatacenterID(id)` - Sets the datacenter ID
- `SetWorkerID(id)` - Sets the worker ID
- `SetStrictMode(strict)` - Enables/disables strict mode
- `Build()` - Builds the snowflake instance

### Instance Methods
- `Generate()` - Generates a unique ID
- `Cleanup()` - Cleans up resources

## Configuration

The library supports three main configuration approaches:

1. **Auto-allocation (Recommended)**: Let Redis assign IDs automatically
2. **Manual Configuration**: Explicitly set datacenter ID and worker ID
3. **Default Values**: Use built-in default values

## Performance

The library is designed for high performance:
- Local ID generation after initialization
- Thread-safe operation
- Minimal Redis interaction in normal mode
- Additional Redis checks in strict mode for enhanced uniqueness

## Custom Redis Client Implementation

The library uses an interface-based approach for Redis clients, allowing you to implement your own Redis client that conforms to the Client interface:

```go
type Client interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Del(ctx context.Context, keys ...string) (int64, error)
}
```

You can implement this interface with any Redis client library of your choice (such as go-redis, redigo, rueidis, goredis/redismock for testing, etc.) and pass your custom implementation to the builder:

```go
type MyCustomRedisClient struct {
	// Your Redis client implementation using your preferred library
}

func (c *MyCustomRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	// Implementation using your preferred Redis client library
}

func (c *MyCustomRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	// Implementation using your preferred Redis client library
}

func (c *MyCustomRedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	// Implementation using your preferred Redis client library
}

// Usage
customRedisClient := &MyCustomRedisClient{}
sf, err := snowflake.NewBuilder().
    SetRedisClient(customRedisClient).
    SetDatacenterID(1).
    SetWorkerID(1).
    Build()
```

This approach provides flexibility to use different Redis client libraries based on your performance, feature, or dependency requirements. The default implementation uses go-redis, but you're not limited to it.

## Notes

- In normal mode, ID generation is completely local after initialization. Redis is only used during setup to coordinate unique identifiers.
- In strict mode, Redis is additionally used during ID generation to provide an extra layer of duplicate prevention.