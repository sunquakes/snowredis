package tests

import (
	"testing"

	"github.com/sunquakes/snowredis/tests/mock"

	"github.com/sunquakes/snowredis/internal/snowflake"
)

// TestStrictModeWithRedis tests using Redis assistance in strict mode to prevent duplicates
func TestStrictModeWithRedis(t *testing.T) {
	mockRedis := mock.NewMockRedisClient()

	// Create a snowflake instance with strict mode enabled
	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		SetStrictMode(true). // Enable strict mode
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake in strict mode: %v", err)
	}
	defer sf.Cleanup()

	// 生成一些ID并验证
	ids := make([]int64, 10)
	for i := 0; i < 10; i++ {
		id, err := sf.Generate()
		if err != nil {
			t.Errorf("Error generating ID in strict mode: %v", err)
			continue
		}
		ids[i] = id

		// Verify the ID is positive
		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}
	}

	// Check ID uniqueness
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found in strict mode: %d", id)
		}
		seen[id] = true
	}
}

// TestStrictModeWithoutRedis tests enabling strict mode without a Redis client
func TestStrictModeWithoutRedis(t *testing.T) {
	// Create a snowflake instance with strict mode enabled but without a Redis client
	sf, err := snowflake.NewBuilder().
		SetDatacenterID(1).
		SetWorkerID(1).
		SetStrictMode(true). // Enable strict mode but without a Redis client
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake in strict mode without Redis: %v", err)
	}
	defer sf.Cleanup()

	// Should fall back to local generation mode
	id, err := sf.Generate()
	if err != nil {
		t.Errorf("Error generating ID in strict mode fallback: %v", err)
	}

	if id <= 0 {
		t.Errorf("Generated ID should be positive, got: %d", id)
	}
}

// TestNormalMode tests normal mode (non-strict mode)
func TestNormalMode(t *testing.T) {
	mockRedis := mock.NewMockRedisClient()

	// Create a snowflake instance without strict mode enabled
	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		SetStrictMode(false). // Disable strict mode (default behavior)
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake in normal mode: %v", err)
	}
	defer sf.Cleanup()

	// 生成一些ID并验证
	ids := make([]int64, 10)
	for i := 0; i < 10; i++ {
		id, err := sf.Generate()
		if err != nil {
			t.Errorf("Error generating ID in normal mode: %v", err)
			continue
		}
		ids[i] = id

		// Verify the ID is positive
		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}
	}

	// Check ID uniqueness
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found in normal mode: %d", id)
		}
		seen[id] = true
	}
}
