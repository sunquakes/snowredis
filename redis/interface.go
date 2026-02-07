// Package redis provides Redis integration for the snowflake ID generator.
package redis

import (
	"context"
	"time"
)

// Client Interface defining basic Redis operations required for ID allocation.
type Client interface {
	/**
	 * SetNX sets a key-value pair if the key does not exist
	 * @param ctx - context for the operation
	 * @param key - string representing the key to set
	 * @param value - interface{} representing the value to set
	 * @param expiration - time.Duration representing the expiration time
	 * @return bool - true if the key was set, false if it already existed
	 * @return error - error if any occurred during the operation
	 */
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)

	/**
	 * Incr increments the value of a key
	 * @param ctx - context for the operation
	 * @param key - string representing the key to increment
	 * @return int64 - the new value after incrementing
	 * @return error - error if any occurred during the operation
	 */
	Incr(ctx context.Context, key string) (int64, error)

	/**
	 * Del deletes keys
	 * @param ctx - context for the operation
	 * @param keys - ...string representing the keys to delete
	 * @return int64 - the number of keys deleted
	 * @return error - error if any occurred during the operation
	 */
	Del(ctx context.Context, keys ...string) (int64, error)
}
