package snowflake

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sunquakes/snowredis/redis"
)

const (
	DefaultDatacenterID    = 1
	DefaultWorkerID        = 1
	AutoAllocateFlag       = true
	NoAllocationFlag       = false
	ZeroValue              = 0
	MaxDatacenterIDPlusOne = 32 // Used to limit datacenter ID to 5 bits
	MaxWorkerIDPlusOne     = 32 // Used to limit worker ID to 5 bits
)

// RedisSnowflake Redis-based snowflake algorithm implementation
type RedisSnowflake struct {
	node          *Node
	redisClient   redis.RedisClient
	ctx           context.Context
	lastTimestamp int64
	strictMode    bool // Strict mode, use Redis assistance to prevent duplicates
}

// RedisSnowflakeBuilder Builder for Redis-based snowflake instance
type RedisSnowflakeBuilder struct {
	client       redis.RedisClient
	datacenterID int64
	workerID     int64
	strictMode   bool // Whether to use strict mode with Redis assistance
}

/**
 * NewBuilder creates a new RedisSnowflakeBuilder instance
 * @return *RedisSnowflakeBuilder - a new builder instance
 */
func NewBuilder() *RedisSnowflakeBuilder {
	return &RedisSnowflakeBuilder{}
}

/**
 * SetRedisClient sets the Redis client for the builder
 * @param client - redis.RedisClient interface implementation
 * @return *RedisSnowflakeBuilder - the builder instance for chaining
 */
func (builder *RedisSnowflakeBuilder) SetRedisClient(client redis.RedisClient) *RedisSnowflakeBuilder {
	builder.client = client
	return builder
}

/**
 * SetDatacenterID sets the datacenter ID for the snowflake instance
 * @param id - int64 representing the datacenter ID (should be 0-31 to fit in 5 bits)
 * @return *RedisSnowflakeBuilder - the builder instance for chaining
 */
func (builder *RedisSnowflakeBuilder) SetDatacenterID(id int64) *RedisSnowflakeBuilder {
	builder.datacenterID = id
	return builder
}

/**
 * SetWorkerID sets the worker ID for the snowflake instance
 * @param id - int64 representing the worker ID (should be 0-31 to fit in 5 bits)
 * @return *RedisSnowflakeBuilder - the builder instance for chaining
 */
func (builder *RedisSnowflakeBuilder) SetWorkerID(id int64) *RedisSnowflakeBuilder {
	builder.workerID = id
	return builder
}

/**
 * SetStrictMode sets whether to use strict mode with Redis assistance to prevent duplicates
 * @param strict - bool indicating whether to enable strict mode
 * @return *RedisSnowflakeBuilder - the builder instance for chaining
 */
func (builder *RedisSnowflakeBuilder) SetStrictMode(strict bool) *RedisSnowflakeBuilder {
	builder.strictMode = strict
	return builder
}

/*
/**
  - Build creates and returns a RedisSnowflake instance based on the configured parameters
  - @return *RedisSnowflake - the configured snowflake instance
  - @return error - any error that occurred during construction
*/
func (builder *RedisSnowflakeBuilder) Build() (*RedisSnowflake, error) {
	// Determine how to build the instance based on priority
	datacenterID, workerID, useRedisAllocation := builder.determineConfiguration()
	if useRedisAllocation {
		// Use Redis for automatic ID allocation
		return builder.createRedisAllocatedInstance(builder.client)
	} else if builder.client != nil {
		// Use manual IDs with Redis client
		return builder.createInstanceWithClient(builder.client, datacenterID, workerID)
	} else {
		// Pure local instance
		return builder.createLocalInstance(datacenterID, workerID, nil)
	}
}

/**
 * determineConfiguration determines the configuration based on priority of provided parameters
 * @return int64 - the datacenter ID to use
 * @return int64 - the worker ID to use
 * @return bool - flag indicating whether to use Redis allocation
 */
func (builder *RedisSnowflakeBuilder) determineConfiguration() (int64, int64, bool) {
	// 1. If datacenterID and workerID (non-zero) are set, prioritize manual values
	if builder.datacenterID != ZeroValue && builder.workerID != ZeroValue {
		return builder.datacenterID, builder.workerID, NoAllocationFlag
	} else if builder.client != nil {
		// 2. If Redis client is set but no manual IDs, allocate automatically via Redis
		return ZeroValue, ZeroValue, AutoAllocateFlag // Return special value indicating Redis allocation
	} else {
		// 3. If neither is set, use default values
		return DefaultDatacenterID, DefaultWorkerID, NoAllocationFlag // Use default values
	}
}

/**
 * generateLocally generates an ID locally without Redis coordination
 * @return int64 - the generated unique ID
 * @return error - any error that occurred during generation (e.g. clock rollback)
 */
func (rs *RedisSnowflake) generateLocally() (int64, error) {
	rs.node.Lock()
	defer rs.node.Unlock()

	timestamp := currentTimeMillis()

	// If timestamp is less than last timestamp, clock rollback occurred
	if timestamp < rs.node.lastTimestamp {
		return 0, errors.New("clock rollback error")
	}

	// If generating in the same millisecond, increment sequence number
	if rs.node.lastTimestamp == timestamp {
		rs.node.sequence = (rs.node.sequence + 1) & maxSequence
		// If sequence number overflows, wait for next millisecond
		if rs.node.sequence == 0 {
			for timestamp <= rs.node.lastTimestamp {
				timestamp = currentTimeMillis()
			}
		}
	} else {
		// Different millisecond, reset sequence number
		rs.node.sequence = 0
	}

	rs.node.lastTimestamp = timestamp

	// Calculate ID
	id := ((timestamp - Epoch) << timestampShift) |
		(rs.node.datacenterID << datacenterShift) |
		(rs.node.workerID << workerShift) |
		rs.node.sequence

	return id, nil
}

/**
 * generateWithRedisAssistance generates an ID using Redis assistance to ensure global uniqueness
 * @return int64 - the generated unique ID
 * @return error - any error that occurred during generation
 */
func (rs *RedisSnowflake) generateWithRedisAssistance() (int64, error) {
	// Try multiple times to generate an ID until we successfully obtain a unique one
	maxRetries := 10
	for attempt := 0; attempt < maxRetries; attempt++ {
		timestamp := currentTimeMillis()

		// Combine timestamp, datacenterID, workerID and attempt count to form a unique ID
		id := ((timestamp - Epoch) << timestampShift) |
			(rs.node.datacenterID << datacenterShift) |
			(rs.node.workerID << workerShift) |
			((timestamp + int64(attempt)) & maxSequence) // Use timestamp+attempt count as sequence part

		// Try to record this ID in Redis using distributed lock mechanism
		key := fmt.Sprintf("snowflake:id:%d", id)
		lockAcquired, err := rs.redisClient.SetNX(rs.ctx, key, "1", time.Hour) // Lock for 1 hour
		if err != nil {
			return 0, fmt.Errorf("failed to acquire lock for ID generation: %v", err)
		}

		if lockAcquired {
			// Successfully acquired lock, ID is unique
			return id, nil
		}

		// If unable to acquire lock, ID is already taken, retry
		// Consider brief delay to reduce contention
		if attempt < maxRetries-1 {
			time.Sleep(time.Millisecond * 1) // Brief sleep to reduce competition
		}
	}

	return 0, fmt.Errorf("failed to generate unique ID after %d attempts", maxRetries)
}

/**
 * Generate generates a unique ID based on the configuration (local or with Redis assistance)
 * @return int64 - the generated unique ID
 * @return error - any error that occurred during generation
 */
func (rs *RedisSnowflake) Generate() (int64, error) {
	// If strict mode is enabled and Redis client exists, use Redis assistance
	if rs.strictMode && rs.redisClient != nil {
		return rs.generateWithRedisAssistance()
	}

	// Otherwise use local generation method
	return rs.generateLocally()
}

/**
 * Cleanup performs cleanup operations for the RedisSnowflake instance
 */
func (rs *RedisSnowflake) Cleanup() {
	// Currently no Redis data needs to be cleaned up
	// Previous node registration functionality has been removed
}

/**
 * createInstance creates a RedisSnowflake instance with the given parameters
 * @param datacenterID - int64 representing the datacenter ID
 * @param workerID - int64 representing the worker ID
 * @param client - redis.RedisClient interface implementation (can be nil for local-only mode)
 * @return *RedisSnowflake - the created instance
 * @return error - any error that occurred during creation
 */
func (builder *RedisSnowflakeBuilder) createInstance(datacenterID, workerID int64, client redis.RedisClient) (*RedisSnowflake, error) {
	node, err := NewNode(datacenterID, workerID)
	if err != nil {
		return nil, err
	}

	return &RedisSnowflake{
		node:          node,
		redisClient:   client,
		ctx:           context.Background(),
		lastTimestamp: 0,
		strictMode:    builder.strictMode,
	}, nil
}

/**
 * createLocalInstance creates a local-only instance for ID generation (without Redis coordination)
 * @param datacenterID - int64 representing the datacenter ID
 * @param workerID - int64 representing the worker ID
 * @param client - redis.RedisClient interface implementation (will be nil for local-only mode)
 * @return *RedisSnowflake - the created local instance
 * @return error - any error that occurred during creation
 */
func (builder *RedisSnowflakeBuilder) createLocalInstance(datacenterID, workerID int64, client redis.RedisClient) (*RedisSnowflake, error) {
	return builder.createInstance(datacenterID, workerID, client)
}

/**
 * createInstanceWithClient creates an instance with a Redis client and specified IDs
 * @param client - redis.RedisClient interface implementation
 * @param datacenterID - int64 representing the datacenter ID
 * @param workerID - int64 representing the worker ID
 * @return *RedisSnowflake - the created instance with Redis client
 * @return error - any error that occurred during creation
 */
func (builder *RedisSnowflakeBuilder) createInstanceWithClient(client redis.RedisClient, datacenterID, workerID int64) (*RedisSnowflake, error) {
	return builder.createInstance(datacenterID, workerID, client)
}

/**
 * createRedisAllocatedInstance creates an instance with IDs automatically allocated by Redis
 * @param client - redis.RedisClient interface implementation
 * @return *RedisSnowflake - the created instance with Redis-allocated IDs
 * @return error - any error that occurred during creation or ID allocation
 */
func (builder *RedisSnowflakeBuilder) createRedisAllocatedInstance(client redis.RedisClient) (*RedisSnowflake, error) {
	ctx := context.Background()

	// Get datacenterID and workerID from Redis
	dcKey := "snowflake:next_datacenter_id"
	datacenterID, err := client.Incr(ctx, dcKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get datacenter ID from Redis: %v", err)
	}
	datacenterID = datacenterID % MaxDatacenterIDPlusOne // Limit to 5 bits

	workerKey := "snowflake:next_worker_id"
	workerID, err := client.Incr(ctx, workerKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get worker ID from Redis: %w", err)
	}
	workerID = workerID % MaxWorkerIDPlusOne // Limit to 5 bits

	return builder.createInstance(datacenterID, workerID, client)
}
