package tests

import (
	"testing"
	"time"

	"snowredis/internal/snowflake"
	"snowredis/tests/mock"
)

// TestConcurrentIDGeneration Tests ID generation in concurrent scenarios
func TestConcurrentIDGeneration(t *testing.T) {
	mockRedis := mock.NewMockRedisClient()

	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake: %v", err)
	}
	defer sf.Cleanup()

	// Generate IDs concurrently
	results := make(chan int64, 100)
	errors := make(chan error, 100)

	// Start multiple goroutines to generate IDs simultaneously
	numWorkers := 10
	idsPerWorker := 10

	for w := 0; w < numWorkers; w++ {
		go func() {
			for i := 0; i < idsPerWorker; i++ {
				id, err := sf.Generate()
				if err != nil {
					errors <- err
					return
				}
				results <- id
			}
		}()
	}

	// Collect results
	ids := make([]int64, 0, numWorkers*idsPerWorker)
	errCount := 0

	for i := 0; i < numWorkers*idsPerWorker; i++ {
		select {
		case id := <-results:
			ids = append(ids, id)
		case err := <-errors:
			t.Errorf("Error generating ID: %v", err)
			errCount++
		case <-time.After(time.Second): // Prevent deadlock
			t.Fatal("Timeout waiting for ID generation")
		}
	}

	// Check ID count
	if len(ids) != numWorkers*idsPerWorker {
		t.Errorf("Expected %d IDs, got %d", numWorkers*idsPerWorker, len(ids))
	}

	// Check ID uniqueness
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		seen[id] = true

		// Check if ID is positive
		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}
	}
}

// TestTimestampExtraction Extracts and validates timestamp from generated ID
func TestTimestampExtraction(t *testing.T) {
	mockRedis := mock.NewMockRedisClient()

	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake: %v", err)
	}
	defer sf.Cleanup()

	// Record start time
	startTime := time.Now().UnixNano() / 1000000 // Convert to milliseconds

	id, err := sf.Generate()
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}

	endTime := time.Now().UnixNano() / 1000000 // Convert to milliseconds

	// Extract timestamp part from ID (top 41 bits)
	timestamp := (id >> 22) + snowflake.Epoch

	// Validate timestamp is within reasonable range
	if timestamp < startTime || timestamp > endTime+1 {
		t.Errorf("Extracted timestamp %d is outside expected range [%d, %d]", timestamp, startTime, endTime)
	}
}

// TestUniqueWithinSameMillisecond Tests uniqueness within the same millisecond
func TestUniqueWithinSameMillisecond(t *testing.T) {
	mockRedis := mock.NewMockRedisClient()

	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake: %v", err)
	}
	defer sf.Cleanup()

	// Generate multiple IDs within the same millisecond
	ids := make([]int64, 10)
	for i := 0; i < 10; i++ {
		id, err := sf.Generate()
		if err != nil {
			t.Fatalf("Failed to generate ID: %v", err)
		}
		ids[i] = id
	}

	// Check all IDs are unique
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		seen[id] = true
	}

	// Check IDs are incrementing (sequence number increments within same millisecond)
	for i := 1; i < len(ids); i++ {
		if ids[i] <= ids[i-1] {
			t.Errorf("IDs should be strictly increasing, got %d followed by %d", ids[i-1], ids[i])
		}
	}
}
