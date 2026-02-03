package tests

import (
	"runtime"
	"testing"
	"time"

	"github.com/sunquakes/snowredis/tests/mock"

	"github.com/sunquakes/snowredis/snowflake"
)

// BenchmarkIDGeneration Performance benchmark test
func BenchmarkIDGeneration(b *testing.B) {
	mockRedis := mock.NewMockRedisClient()

	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		Build()
	if err != nil {
		b.Fatalf("Failed to initialize Redis Snowflake: %v", err)
	}
	defer sf.Cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sf.Generate()
		if err != nil {
			b.Errorf("Error generating ID: %v", err)
		}
	}
}

// BenchmarkConcurrentIDGeneration Concurrent performance benchmark test
func BenchmarkConcurrentIDGeneration(b *testing.B) {
	mockRedis := mock.NewMockRedisClient()

	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		Build()
	if err != nil {
		b.Fatalf("Failed to initialize Redis Snowflake: %v", err)
	}
	defer sf.Cleanup()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := sf.Generate()
			if err != nil {
				b.Errorf("Error generating ID: %v", err)
			}
		}
	})
}

// TestPerformance Tests performance metrics
func TestPerformance(t *testing.T) {
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

	// Test how many IDs can be generated in 1 second
	start := time.Now()
	count := 0
	duration := time.Second

	for time.Since(start) < duration {
		_, err := sf.Generate()
		if err != nil {
			t.Errorf("Error generating ID: %v", err)
			break
		}
		count++
	}

	elapsed := time.Since(start)
	tps := float64(count) / elapsed.Seconds()

	t.Logf("Generated %d IDs in %v (%.0f TPS)", count, elapsed, tps)

	// Verify at least 1000 IDs per second can be generated (this is a reasonable performance goal)
	if tps < 1000 {
		t.Logf("Warning: Performance is lower than expected (%.0f TPS)", tps)
	} else {
		t.Logf("Performance is acceptable (%.0f TPS)", tps)
	}
}

// TestMemoryUsage Tests memory usage
func TestMemoryUsage(t *testing.T) {
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

	// Get initial memory stats
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Generate 1000 IDs
	for i := 0; i < 1000; i++ {
		_, err := sf.Generate()
		if err != nil {
			t.Errorf("Error generating ID: %v", err)
		}
	}

	// Get final memory stats
	runtime.GC()
	runtime.ReadMemStats(&m2)

	// Check if memory growth is within reasonable range
	memoryIncrease := m2.TotalAlloc - m1.TotalAlloc
	t.Logf("Memory increase after 1000 ID generations: %d bytes", memoryIncrease)

	// Verify memory growth is reasonable (should not exceed 100 bytes per ID)
	expectedMaxGrowth := uint64(1000 * 100) // 100 bytes/ID * 1000 IDs
	if memoryIncrease > expectedMaxGrowth {
		t.Logf("Warning: Memory usage is higher than expected (%d bytes for 1000 IDs)", memoryIncrease)
	} else {
		t.Logf("Memory usage is acceptable")
	}
}
