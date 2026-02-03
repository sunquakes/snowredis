package tests

import (
	"testing"

	"snowredis/internal/snowflake"
	"snowredis/tests/mock"
)

// TestStrictModeWithRedis 测试严格模式下使用Redis辅助防止重复
func TestStrictModeWithRedis(t *testing.T) {
	mockRedis := mock.NewMockRedisClient()

	// 创建启用严格模式的雪花算法实例
	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		SetStrictMode(true). // 启用严格模式
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

		// 验证ID是正数
		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}
	}

	// 检查ID唯一性
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found in strict mode: %d", id)
		}
		seen[id] = true
	}
}

// TestStrictModeWithoutRedis 测试在没有Redis客户端的情况下启用严格模式
func TestStrictModeWithoutRedis(t *testing.T) {
	// 创建启用严格模式但没有Redis客户端的雪花算法实例
	sf, err := snowflake.NewBuilder().
		SetDatacenterID(1).
		SetWorkerID(1).
		SetStrictMode(true). // 启用严格模式，但没有Redis客户端
		Build()
	if err != nil {
		t.Fatalf("Failed to initialize Redis Snowflake in strict mode without Redis: %v", err)
	}
	defer sf.Cleanup()

	// 应该退回到本地生成模式
	id, err := sf.Generate()
	if err != nil {
		t.Errorf("Error generating ID in strict mode fallback: %v", err)
	}

	if id <= 0 {
		t.Errorf("Generated ID should be positive, got: %d", id)
	}
}

// TestNormalMode 测试正常模式（非严格模式）
func TestNormalMode(t *testing.T) {
	mockRedis := mock.NewMockRedisClient()

	// 创建未启用严格模式的雪花算法实例
	sf, err := snowflake.NewBuilder().
		SetRedisClient(mockRedis).
		SetDatacenterID(1).
		SetWorkerID(1).
		SetStrictMode(false). // 不启用严格模式（默认行为）
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

		// 验证ID是正数
		if id <= 0 {
			t.Errorf("Generated ID should be positive, got: %d", id)
		}
	}

	// 检查ID唯一性
	seen := make(map[int64]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("Duplicate ID found in normal mode: %d", id)
		}
		seen[id] = true
	}
}