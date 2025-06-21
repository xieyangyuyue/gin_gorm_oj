package test

import (
	"context" // 导入 context 包，用于上下文管理，控制请求生命周期（如超时）
	"errors"
	"fmt"                          // 导入 fmt 包，用于格式化输出，用于测试结果展示
	"gin_gorm_oj/models"           // 导入 models 包，包含 Redis 客户端的初始化
	"gin_gorm_oj/utils"            // 导入 utils 包，用于获取 Redis 配置
	"github.com/go-redis/redis/v8" // 导入 Redis 客户端库
	"sync"                         // 导入 sync 包，用于并发测试中的同步原语，如 WaitGroup
	"testing"                      // 导入 Go 测试框架，提供测试功能
	"time"                         // 导入 time 包，用于时间处理，用于设置键值过期时间
)

// 全局上下文对象
// 在 Redis 操作中传递上下文，可用于实现超时控制、取消操作等（当前使用默认空上下文）。
var ctx = context.Background() //

// Redis 客户端实例
// 使用 NewClient 方法创建连接池，配置参数说明：
// Addr:     Redis 服务器地址（默认本地6379端口）
// Password: 认证密码，此处需根据实际环境修改
// DB:       选择使用的数据库编号（0-15，默认0）
var rdb = redis.NewClient(&redis.Options{ //
	Addr:     fmt.Sprintf("%s:%s", utils.RedisHost, utils.RedisPort), // Redis 服务器地址，使用格式化拼接确保有冒号
	Password: utils.RedisPassWord,                                    // 认证密码
	DB:       utils.RedisNumber,                                      // 选择使用的数据库编号
})

// setupRedisTestKey 在每个测试开始前生成一个唯一的键，并确保测试结束后清理。
// 返回生成的键名。
func setupRedisTestKey(t *testing.T) string {
	// 修正：使用 t.Name() 获取当前测试的完整名称，确保在并行测试中的唯一性。
	// t.Name() 会返回 TestFunc/SubTestName 的形式。
	testKey := fmt.Sprintf("test_key_%s_%d", t.Name(), time.Now().UnixNano()) // 生成一个基于测试名和纳秒时间戳的唯一键
	t.Cleanup(func() {                                                        // 注册一个清理函数，在当前测试函数结束后执行
		err := rdb.Del(ctx, testKey).Err() // 删除测试键
		if err != nil {
			t.Logf("Failed to clean up test key %s: %v", testKey, err) // 记录清理失败日志
		}
	})
	return testKey
}

// TestRedisSet_Unit 测试 Redis 的字符串写入功能（单元测试维度）。
// 功能：设置键值对，验证基础写入功能与过期时间配置。
func TestRedisSet_Unit(t *testing.T) { // 定义一个名为 TestRedisSet_Unit 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	testKey := setupRedisTestKey(t) // 获取一个唯一的测试键，并设置清理

	testValue := "unit_test_value" // 定义要设置的值
	// 设置键值对，并设置 60 秒的过期时间。
	err := rdb.Set(ctx, testKey, testValue, 60*time.Second).Err() //
	if err != nil {                                               // 检查 Redis Set 操作是否出错
		t.Fatalf("Redis Set operation failed: %v", err) // 如果出错，则致命错误并终止测试
	}
	t.Logf("Successfully set key %s with value %s", testKey, testValue) // 记录成功设置信息
}

// TestRedisGet_Unit 测试 Redis 的字符串读取功能（单元测试维度）。
// 功能：验证 TestRedisSet 写入的数据一致性。
func TestRedisGet_Unit(t *testing.T) { // 定义一个名为 TestRedisGet_Unit 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	testKey := setupRedisTestKey(t) // 获取一个唯一的测试键，并设置清理
	testValue := "get_test_value"   // 定义要设置的值

	// 首先设置一个键值对，确保有数据可读
	err := rdb.Set(ctx, testKey, testValue, 10*time.Second).Err() // 设置键值对
	if err != nil {                                               // 检查设置是否出错
		t.Fatalf("Setup for Get test failed: %v", err) // 如果设置失败，则致命错误
	}

	// 使用 Get 方法获取键的值。
	// Result() 方法同时返回值和错误信息。
	v, err := rdb.Get(ctx, testKey).Result() // 获取键的值
	if err != nil {                          // 检查 Redis Get 操作是否出错
		t.Fatalf("Redis Get operation failed: %v", err) // 如果出错，则致命错误
	}

	// 验证获取到的值是否与预期一致。
	if v != testValue { // 检查获取到的值是否与预期值匹配
		t.Errorf("Expected value %q, got %q", testValue, v) // 如果不匹配，则报错
	}
	t.Logf("Successfully got value for key %s: %s", testKey, v) // 记录成功获取信息
}

// TestRedisGetByModels_Integration 验证 models 包的 Redis 配置（集成测试维度）。
// 功能：确保项目内其他模块通过统一配置能正常访问 Redis。
func TestRedisGetByModels_Integration(t *testing.T) { // 定义一个名为 TestRedisGetByModels_Integration 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	testKey := setupRedisTestKey(t)       // 获取一个唯一的测试键，并设置清理
	testValue := "integration_test_value" // 定义要设置的值

	// 使用直接的 rdb 实例设置数据，模拟外部系统写入
	err := rdb.Set(ctx, testKey, testValue, 10*time.Second).Err() // 设置键值对
	if err != nil {                                               // 检查设置是否出错
		t.Fatalf("Setup for models integration test failed: %v", err) // 如果设置失败，则致命错误
	}

	// 使用 models 包导出的 RDB 实例进行操作，验证其配置和可用性。
	v, err := models.RDB.Get(ctx, testKey).Result() // 通过 models.RDB 获取键的值
	if err != nil {                                 // 检查通过 models.RDB 获取数据是否出错
		t.Fatalf("Getting data via models.RDB failed: %v", err) // 如果失败，则致命错误
	}

	// 验证获取到的值是否与预期一致。
	if v != testValue { // 检查获取到的值是否与预期值匹配
		t.Errorf("Expected value %q from models.RDB, got %q", testValue, v) // 如果不匹配，则报错
	}
	t.Logf("Successfully got value for key %s via models.RDB: %s", testKey, v) // 记录成功获取信息
}

// TestRedis_SystemConcurrency 是一个系统测试函数，模拟高并发环境下对 Redis 的读写操作。
func TestRedis_SystemConcurrency(t *testing.T) { // 定义一个名为 TestRedis_SystemConcurrency 的测试函数
	t.Log("开始系统测试：高并发 Redis 读写场景") // 记录测试开始信息

	// 修正：使用 t.Name() 获取当前测试的完整名称作为前缀的一部分
	testKeyPrefix := fmt.Sprintf("concurrency_%s_%d_", t.Name(), time.Now().UnixNano()) // 为并发测试生成一个键前缀
	numOperations := 10                                                                 // 每个并发操作的次数
	numWorkers := 100                                                                   // 并发工作者（goroutine）的数量

	var wg sync.WaitGroup                   // 声明一个 WaitGroup
	errChan := make(chan error, numWorkers) // 用于收集并发操作中的错误

	t.Cleanup(func() { // 注册一个清理函数，删除所有测试相关的键
		// 使用 SCAN 命令查找并删除所有带前缀的键
		iter := rdb.Scan(ctx, 0, testKeyPrefix+"*", 0).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			err := rdb.Del(ctx, key).Err()
			if err != nil {
				t.Logf("Failed to clean up concurrency key %s: %v", key, err)
			}
		}
		if err := iter.Err(); err != nil {
			t.Logf("Error during scan for cleanup: %v", err)
		}
		t.Logf("Cleaned up keys with prefix %s", testKeyPrefix)
	})

	for i := 0; i < numWorkers; i++ { // 启动多个并发工作者
		wg.Add(1)               // 增加 WaitGroup 计数
		go func(workerID int) { // 启动新的 goroutine
			defer wg.Done() // goroutine 结束后调用 Done()

			for j := 0; j < numOperations; j++ { // 每个工作者执行多次读写操作
				key := fmt.Sprintf("%s%d_%d", testKeyPrefix, workerID, j) // 生成唯一的键名
				value := fmt.Sprintf("value_w%d_op%d", workerID, j)       // 生成对应的值

				// 模拟写入操作
				err := rdb.Set(ctx, key, value, 5*time.Second).Err() // 设置键值对
				if err != nil {                                      // 检查写入是否出错
					errChan <- fmt.Errorf("worker %d: Set operation failed for key %s: %w", workerID, key, err) // 将错误发送到错误通道
					continue
				}

				// 模拟读取操作
				retrievedValue, err := rdb.Get(ctx, key).Result() // 获取键的值
				if err != nil {                                   // 检查读取是否出错
					errChan <- fmt.Errorf("worker %d: Get operation failed for key %s: %w", workerID, key, err) // 将错误发送到错误通道
					continue
				}

				// 验证读取到的值
				if retrievedValue != value { // 检查读取到的值是否与写入的值一致
					errChan <- fmt.Errorf("worker %d: Value mismatch for key %s. Expected %q, got %q", workerID, key, value, retrievedValue) // 如果不一致，则发送错误
				}
			}
			t.Logf("工作者 %d 完成 %d 个并发 Redis 操作", workerID, numOperations) // 记录工作者完成信息
		}(i) // 传入 workerID
	}

	wg.Wait()      // 等待所有工作者完成
	close(errChan) // 关闭错误通道

	// 检查是否有错误发生
	for err := range errChan { // 遍历错误通道中的所有错误
		t.Error(err) // 记录错误
	}
	if t.Failed() { // 检查是否有测试失败
		t.Error("系统测试失败：在高并发 Redis 场景下发现错误。") // 如果有失败，则报错
	} else {
		t.Log("系统测试完成：高并发 Redis 读写操作通过。") // 记录测试通过
	}
}

// TestRedis_Validation 是一个验证测试函数，用于严格验证 Redis 客户端的特定行为。
func TestRedis_Validation(t *testing.T) { // 定义一个名为 TestRedis_Validation 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	testKey := setupRedisTestKey(t) // 获取一个唯一的测试键，并设置清理

	t.Run("Set with Zero Expiration", func(t *testing.T) { // 验证设置零秒过期时间（应视为永久）
		err := rdb.Set(ctx, testKey, "zero_ttl_value", 0).Err() // 设置键值，过期时间为 0
		if err != nil {                                         // 检查设置是否出错
			t.Fatalf("Failed to set key with zero expiration: %v", err) // 如果出错，则致命错误
		}
		// 验证键是否存在，且 TTL 应该为 -1 (永久)
		ttl, err := rdb.TTL(ctx, testKey).Result() // 获取键的 TTL
		if err != nil {                            // 检查获取 TTL 是否出错
			t.Fatalf("Failed to get TTL for key %s: %v", testKey, err) // 如果出错，则致命错误
		}
		if ttl != -1 { // 检查 TTL 是否为 -1 (表示永久)
			t.Errorf("Expected TTL -1 for zero expiration, got %v", ttl) // 如果不为 -1，则报错
		}
		t.Logf("Set with zero expiration validated. Key %s, TTL: %v", testKey, ttl) // 记录验证成功
	})

	t.Run("Get Non-existent Key", func(t *testing.T) { // 验证获取不存在的键
		nonExistentKey := setupRedisTestKey(t)            // 获取一个不存在的键 (确保不会被写入)
		val, err := rdb.Get(ctx, nonExistentKey).Result() // 获取不存在的键
		if err != nil && !errors.Is(err, redis.Nil) {     // 如果有错误且不是 redis.Nil 错误
			t.Errorf("Expected redis.Nil for non-existent key, got %v", err) // 报错，因为预期是 redis.Nil
		}
		if val != "" { // 检查返回值是否为空字符串
			t.Errorf("Expected empty string for non-existent key, got %q", val) // 报错，因为预期是空字符串
		}
		t.Logf("Getting non-existent key %s returned expected nil error.", nonExistentKey) // 记录验证成功
	})
}

// TestRedis_Regression 是一个回归测试函数，确保在未来的代码变更后，Redis 操作行为仍然一致且正确。
func TestRedis_Regression(t *testing.T) { // 定义一个名为 TestRedis_Regression 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	testKey := setupRedisTestKey(t) // 获取一个唯一的测试键，并设置清理
	testValue := "regression_value" // 定义测试值

	t.Run("Basic Set/Get Consistency", func(t *testing.T) { // 验证基本的 Set/Get 操作一致性
		err := rdb.Set(ctx, testKey, testValue, 10*time.Second).Err() // 设置键值
		if err != nil {                                               // 检查设置是否出错
			t.Fatalf("Regression test failed: Set operation error: %v", err) // 如果出错，则致命错误
		}

		retrievedVal, err := rdb.Get(ctx, testKey).Result() // 获取键值
		if err != nil {                                     // 检查获取是否出错
			t.Fatalf("Regression test failed: Get operation error: %v", err) // 如果出错，则致命错误
		}

		if retrievedVal != testValue { // 检查值是否一致
			t.Errorf("Regression test failed: Value mismatch. Expected %q, got %q", testValue, retrievedVal) // 如果不一致，则报错
		}
		t.Logf("Regression test: Basic Set/Get consistency passed for key %s.", testKey) // 记录一致性
	})

	t.Run("TTL Behavior Consistency", func(t *testing.T) { // 验证 TTL 行为一致性
		testKeyTTL := setupRedisTestKey(t) // 获取另一个唯一键
		// 设置一个短 TTL 的键
		err := rdb.Set(ctx, testKeyTTL, "temp_value", 2*time.Second).Err() // 设置一个 2 秒 TTL 的键
		if err != nil {                                                    // 检查设置是否出错
			t.Fatalf("Regression test failed: Set with TTL error: %v", err) // 如果出错，则致命错误
		}

		// 检查 TTL 是否大于 0
		ttl, err := rdb.TTL(ctx, testKeyTTL).Result() // 获取 TTL
		if err != nil {                               // 检查获取 TTL 是否出错
			t.Fatalf("Regression test failed: TTL check error: %v", err) // 如果出错，则致命错误
		}
		if ttl <= 0 || ttl > 2*time.Second { // 检查 TTL 是否在预期范围内
			t.Errorf("Regression test failed: Unexpected TTL value. Expected >0 and <=2s, got %v", ttl) // 如果不在预期范围，则报错
		}

		time.Sleep(3 * time.Second) // 等待键过期

		// 再次尝试获取，预期会是 redis.Nil 错误
		_, err = rdb.Get(ctx, testKeyTTL).Result() // 再次获取键
		if !errors.Is(err, redis.Nil) {            // 检查错误是否为 redis.Nil
			t.Errorf("Regression test failed: Expected redis.Nil after expiration, got %v", err) // 如果不是 redis.Nil，则报错
		}
		t.Logf("Regression test: TTL behavior consistent for key %s.", testKeyTTL) // 记录一致性
	})
	t.Log("Redis 回归测试完成。") // 记录回归测试完成信息
}
