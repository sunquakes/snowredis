package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// Config Defines Redis connection configuration
type Config struct {
	Addr string
	Pwd  string
	Db   int
}

// RedisWrapper Wraps *redis.Client to implement RedisClient interface
type RedisWrapper struct {
	client *redis.Client
}

// SetNX Implements RedisClient interface
func (rw *RedisWrapper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return rw.client.SetNX(ctx, key, value, expiration).Result()
}

// Incr Implements RedisClient interface
func (rw *RedisWrapper) Incr(ctx context.Context, key string) (int64, error) {
	return rw.client.Incr(ctx, key).Result()
}

// Del Implements RedisClient interface
func (rw *RedisWrapper) Del(ctx context.Context, keys ...string) (int64, error) {
	return rw.client.Del(ctx, keys...).Result()
}

// NewClient Creates and returns a Redis client instance
func NewClient(cfg *Config) (*RedisWrapper, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.Db,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	redisWrapper := &RedisWrapper{client: client}

	fmt.Println("Connected to Redis successfully")
	return redisWrapper, nil
}

// NewRedisClient Creates and returns a Redis client instance (implements RedisClient interface)
func NewRedisClient(cfg *Config) (RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Pwd,
		DB:       cfg.Db,
	})

	// Test connection
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	redisWrapper := &RedisWrapper{client: client}

	fmt.Println("Connected to Redis successfully")
	return redisWrapper, nil
}
