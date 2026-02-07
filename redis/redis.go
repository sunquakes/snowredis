package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config Defines Redis connection configuration.
type Config struct {
	Addr string
	Pwd  string
	Db   int
}

// Wrapper Wraps *redis.Client to implement Client interface.
type Wrapper struct {
	client *redis.Client
}

// SetNX sets a key-value pair if the key does not exist
// @param ctx - context for the operation
// @param key - string representing the key to set
// @param value - interface{} representing the value to set
// @param expiration - time.Duration representing the expiration time
// @return bool - true if the key was set, false if it already existed
// @return error - error if any occurred during the operation
func (r *Wrapper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

// Incr increments the value of a key
// @param ctx - context for the operation
// @param key - string representing the key to increment
// @return int64 - the new value after incrementing
// @return error - error if any occurred during the operation
func (r *Wrapper) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// Del deletes keys
// @param ctx - context for the operation
// @param keys - ...string representing the keys to delete
// @return int64 - the number of keys deleted
// @return error - error if any occurred during the operation
func (r *Wrapper) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(ctx, keys...).Result()
}

// NewClient Creates and returns a Redis client instance
// @param cfg - *Config containing Redis connection configuration
// @return *Wrapper - the created Redis wrapper instance
// @return error - any error that occurred during connection
func NewClient(cfg *Config) (*Wrapper, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.Db,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	redisWrapper := &Wrapper{client: client}

	// Log successful connection (optional - can be removed if logging is not desired)
	// fmt.Printf("Connected to Redis successfully\n")
	return redisWrapper, nil
}

// NewRedisClient Creates and returns a Redis client instance (implements Client interface)
// @param cfg - *Config containing Redis connection configuration
// @return Client - the created Redis client instance that implements the interface
// @return error - any error that occurred during connection
func NewRedisClient(cfg *Config) (Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.Db,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	redisWrapper := &Wrapper{client: client}

	// Log successful connection (optional - can be removed if logging is not desired)
	// fmt.Printf("Connected to Redis successfully\n")
	return redisWrapper, nil
}
