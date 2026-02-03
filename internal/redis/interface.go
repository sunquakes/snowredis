package redis

import (
	"context"
	"time"
)

/**
 * RedisClient Interface defining basic Redis operations required for ID allocation
 */
type RedisClient interface {
	/**
	 * SetNX sets a key-value pair if the key does not exist
	 * @param ctx - context for the operation
	 * @param key - string representing the key to set
	 * @param value - interface{} representing the value to set
	 * @param expiration - time.Duration for key expiration
	 * @return bool - true if the key was set, false otherwise
	 * @return error - any error that occurred during the operation
	 */
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)

	/**
	 * Incr increments the value of a key
	 * @param ctx - context for the operation
	 * @param key - string representing the key to increment
	 * @return int64 - the incremented value
	 * @return error - any error that occurred during the operation
	 */
	Incr(ctx context.Context, key string) (int64, error)

	/**
	 * Del deletes one or more keys
	 * @param ctx - context for the operation
	 * @param keys - variadic string parameters representing the keys to delete
	 * @return int64 - the number of deleted keys
	 * @return error - any error that occurred during the operation
	 */
	Del(ctx context.Context, keys ...string) (int64, error)
}
