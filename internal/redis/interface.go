package redis

import (
	"context"
	"time"
)

// RedisClient Interface defining basic Redis operations required for ID allocation
type RedisClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
	Del(ctx context.Context, keys ...string) (int64, error)
}
