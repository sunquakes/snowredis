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

// RedisWrapper Wraps *redis.Client to implement RedisClient interface.
type RedisWrapper struct {
	client *redis.Client
}

/**
 * SetNX Implements RedisClient interface - sets a key-value pair if the key does not exist
 * @param ctx - context for the operation
 * @param key - string representing the key to set
 * @param value - interface{} representing the value to set
 * @param expiration - time.Duration for key expiration
 * @return bool - true if the key was set, false otherwise
 * @return error - any error that occurred during the operation
 */
func (rw *RedisWrapper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return rw.client.SetNX(ctx, key, value, expiration).Result()
}

/**
 * Incr Implements RedisClient interface - increments the value of a key
 * @param ctx - context for the operation
 * @param key - string representing the key to increment
 * @return int64 - the incremented value
 * @return error - any error that occurred during the operation
 */
func (rw *RedisWrapper) Incr(ctx context.Context, key string) (int64, error) {
	return rw.client.Incr(ctx, key).Result()
}

/**
 * Del Implements RedisClient interface - deletes one or more keys
 * @param ctx - context for the operation
 * @param keys - variadic string parameters representing the keys to delete
 * @return int64 - the number of deleted keys
 * @return error - any error that occurred during the operation
 */
func (rw *RedisWrapper) Del(ctx context.Context, keys ...string) (int64, error) {
	return rw.client.Del(ctx, keys...).Result()
}

/**
 * NewClient Creates and returns a Redis client instance
 * @param cfg - *Config containing Redis connection configuration
 * @return *RedisWrapper - the created Redis wrapper instance
 * @return error - any error that occurred during connection
 */
func NewClient(cfg *Config) (*RedisWrapper, error) {
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

	redisWrapper := &RedisWrapper{client: client}

	// Log successful connection (optional - can be removed if logging is not desired)
	// fmt.Printf("Connected to Redis successfully\n")
	return redisWrapper, nil
}

/**
 * NewRedisClient Creates and returns a Redis client instance (implements RedisClient interface)
 * @param cfg - *Config containing Redis connection configuration
 * @return RedisClient - the created Redis client instance that implements the interface
 * @return error - any error that occurred during connection
 */
func NewRedisClient(cfg *Config) (RedisClient, error) {
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

	redisWrapper := &RedisWrapper{client: client}

	// Log successful connection (optional - can be removed if logging is not desired)
	// fmt.Printf("Connected to Redis successfully\n")
	return redisWrapper, nil
}
