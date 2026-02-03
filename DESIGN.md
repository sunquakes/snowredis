# Snowflake-Redis Design Document

## Overview

Snowflake-Redis is a distributed ID generation system that combines the Twitter Snowflake algorithm with Redis coordination to ensure unique ID generation across distributed systems. The system generates 64-bit unique IDs composed of timestamp, datacenter ID, worker ID, and sequence number.

## Architecture

### Core Components

#### 1. Snowflake Algorithm Implementation
- **Node Structure**: Contains datacenter ID, worker ID, sequence number, and last timestamp
- **ID Format**: 64-bit integer with the following bit distribution:
  - 1 bit: Sign (always 0)
  - 41 bits: Timestamp (milliseconds since epoch)
  - 5 bits: Datacenter ID
  - 5 bits: Worker ID
  - 12 bits: Sequence number

#### 2. Redis Coordination Layer
- **RedisClient Interface**: Abstracts Redis operations for pluggable implementations
- **ID Allocation**: Automatic allocation of datacenter and worker IDs via Redis
- **Node Registration**: Ensures unique node identification across the cluster

#### 3. Builder Pattern
- **Flexible Configuration**: Allows multiple ways to configure the snowflake instance
- **Method Chaining**: Provides clean, readable API for setup
- **Multiple Modes**: Supports auto-allocation, manual configuration, and custom client modes

## Key Features

### 1. High Performance
- Local ID generation after initialization
- Minimal Redis interaction
- Thread-safe operations with mutex protection
- Capacity for millions of IDs per second

### 2. Distributed Uniqueness
- Coordinated ID allocation during startup
- Unique datacenter and worker IDs per node
- Sequence number handling for same-millisecond IDs
- Clock drift protection

### 3. Flexible Deployment
- Multiple initialization modes
- Pluggable Redis client implementations
- Configurable datacenter and worker IDs
- Standalone operation after initialization

### 4. Strict Mode for Enhanced Uniqueness
- Optional Redis-assisted duplicate prevention
- Additional layer of uniqueness verification
- Configurable via SetStrictMode method
- Falls back to local generation if Redis unavailable

## Implementation Details

### ID Generation Process
1. Acquire mutex for thread safety
2. Check for clock drift (time moving backwards)
3. Handle sequence overflow within same millisecond
4. Combine timestamp, datacenter ID, worker ID, and sequence
5. Return unique 64-bit ID

### Redis Operations
- **ID Allocation**: INCR operations to assign unique datacenter/worker IDs
- **Node Registration**: SETNX operations to ensure node uniqueness
- **Coordination**: Atomic operations to prevent conflicts

### Error Handling
- Clock drift detection and error reporting
- Redis connection failure handling
- Duplicate ID prevention
- Graceful degradation when Redis is unavailable

## Usage Patterns

### Auto-allocation Mode
```go
// Automatically assign IDs from Redis
sf, err := snowflake.NewSnowflake()
```

### Manual Configuration Mode
```go
// Specify IDs manually
sf, err := snowflake.NewRedisSnowflake(cfg.NodeID, datacenterID, workerID)
```

### Builder Mode
```go
// Flexible configuration with builder
sf, err := snowflake.NewRedisSnowflakeBuilder().
    SetRedisClient(client).
    SetDatacenterID(1).
    SetWorkerID(1).
    Build()
```

## Performance Characteristics

- **Throughput**: Over 1 million IDs per second per instance
- **Latency**: Sub-microsecond generation time (after initialization)
- **Scalability**: Up to 32 datacenters × 32 workers per datacenter
- **Uniqueness**: Guaranteed globally unique IDs across all nodes

## Safety Measures

### Clock Drift Protection
- Monitors system time for backward movement
- Returns error when clock rollback is detected
- Prevents duplicate ID generation due to time adjustments

### Sequence Overflow Handling
- Resets sequence to 0 when reaching maximum (4095)
- Waits for next millisecond if sequence exhausted
- Ensures uniqueness within the same millisecond

### Thread Safety
- Mutex protection around critical sections
- Atomic operations for shared state
- Safe concurrent ID generation

## Extensibility

### Custom Redis Clients
The `RedisClient` interface allows plugging in different Redis implementations:
- Standard Redis
- Clustered Redis
- Custom caching layers
- Mock implementations for testing

### Configuration Options
- Adjustable bit allocations for different scales
- Custom epoch timestamps
- Pluggable storage backends

## Testing Strategy

### Unit Tests
- Individual function validation
- Edge case handling
- Error condition testing

### Performance Tests
- Throughput benchmarks
- Concurrency stress testing
- Memory usage validation

### Integration Tests
- Full workflow validation
- Multi-node coordination
- Redis interaction verification

## Future Enhancements

- Health checks and monitoring
- Metrics collection
- Graceful shutdown procedures
- Enhanced error recovery
- Additional ID formats