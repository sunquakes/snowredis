package tests

import (
	"testing"

	"github.com/sunquakes/snowredis/tests/mock"

	"github.com/sunquakes/snowredis/snowflake"
)

func TestRedisSnowflakeGeneration(t *testing.T) {
	// Create mock Redis client
	mockRedis := mock.NewMockRedisClient()

	// Create snowflake instance using Builder pattern
	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake: %v", err)
	}
	defer sf.Cleanup()

	// Generate some IDs and validate
	ids := make([]int64, 10)
	for i := 0; i < 10; i++ {
		id, err := sf.Generate()
		if err != nil {
			t.Errorf("Error generating ID: %v", err)
			continue
		}
		ids[i] = id

		// Validate ID is positive
		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}

		// Validate ID uniqueness (for consecutive generated IDs)
		if i > 0 && ids[i] <= ids[i-1] {
			t.Errorf("IDs should be increasing, got %d after %d", ids[i], ids[i-1])
		}
	}
}

func TestRedisSnowflakeWithAutoAllocation(t *testing.T) {
	// Create mock Redis client
	mockRedis := mock.NewMockRedisClient()

	// Test auto-allocation mode
	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake with auto allocation: %v", err)
	}
	defer sf.Cleanup()

	// Generate some IDs for testing
	for i := 0; i < 5; i++ {
		id, err := sf.Generate()
		if err != nil {
			t.Errorf("Error generating ID: %v", err)
			continue
		}

		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}
	}
}

func TestRedisSnowflakeWithDefaultValues(t *testing.T) {
	// Test using default values
	sf, err := snowflake.NewBuilder().Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake with defaults: %v", err)
	}
	defer sf.Cleanup()

	// Generate some IDs for testing
	for i := 0; i < 5; i++ {
		id, err := sf.Generate()
		if err != nil {
			t.Errorf("Error generating ID: %v", err)
			continue
		}

		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}
	}
}
