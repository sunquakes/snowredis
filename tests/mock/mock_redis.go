package mock

import (
	"context"
	"time"
)

// MockRedisClient Mock Redis client for testing
type MockRedisClient struct {
	data map[string]interface{}
}

// NewMockRedisClient Creates a mock Redis client
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]interface{}),
	}
}

// SetNX Sets value only if key does not exist
func (m *MockRedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	_, exists := m.data[key]
	if exists {
		return false, nil
	}
	m.data[key] = value
	return true, nil
}

// Incr Increments the value of a key
func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	var val int64 = 1
	if v, exists := m.data[key]; exists {
		val = v.(int64) + 1
	}
	m.data[key] = val
	return val, nil
}

// Del Deletes keys
func (m *MockRedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	count := 0
	for _, key := range keys {
		if _, exists := m.data[key]; exists {
			delete(m.data, key)
			count++
		}
	}
	return int64(count), nil
}
