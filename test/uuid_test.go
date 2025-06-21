package test // 声明包名为 test，通常用于存放测试代码

import (
	"fmt"                            // 导入 fmt 包，用于格式化输出
	uuid "github.com/satori/go.uuid" // 导入 go.uuid 库，并将其别名为 uuid，用于生成UUID
	"regexp"                         // 导入 regexp 包，用于正则表达式匹配，验证UUID格式
	"sync"                           // 导入 sync 包，用于同步原语，例如 WaitGroup，在高并发测试中使用
	"testing"                        // 导入 testing 包，Go语言内置的测试框架
)

// uuidV4Regex 是一个正则表达式，用于验证UUID版本4的格式。
// UUID V4 格式：xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
// 其中 x 是十六进制数字，y 是 8, 9, A, 或 B。
var uuidV4Regex = regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$")

// TestGenerateUUID_BasicAssertions 是一个单元测试函数，用于基本验证UUID的生成。
func TestGenerateUUID_BasicAssertions(t *testing.T) { // 定义一个名为 TestGenerateUUID_BasicAssertions 的测试函数，遵循Go测试函数命名规范
	t.Parallel() // 标记该测试可以并行运行，提高测试效率

	id := uuid.NewV4().String() // 调用 uuid 库生成一个UUID字符串

	// 单元测试 - 长度检查
	expectedLength := 36           // UUID字符串的标准长度是36（32个十六进制字符 + 4个连字符）
	if len(id) != expectedLength { // 检查生成的UUID字符串长度是否符合预期
		t.Errorf("生成的UUID长度不正确，期望 %d，得到 %d", expectedLength, len(id)) // 如果长度不符，则报错
	}

	// 单元测试 - 格式检查（使用正则表达式）
	if !uuidV4Regex.MatchString(id) { // 使用预定义的正则表达式检查UUID字符串是否符合版本4的格式
		t.Errorf("生成的UUID格式不正确: %s", id) // 如果格式不符，则报错
	}

	// 单元测试 - 验证非空
	if id == "" { // 检查生成的UUID字符串是否为空
		t.Error("生成的UUID为空") // 如果为空，则报错
	}
}

// TestGenerateUUID_Uniqueness 是一个单元测试函数，用于验证UUID的唯一性（在一定程度上）。
// 严格的唯一性保证需要通过大量生成和冲突检测，这里做基本验证。
func TestGenerateUUID_Uniqueness(t *testing.T) { // 定义一个名为 TestGenerateUUID_Uniqueness 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	uuids := make(map[string]struct{}) // 创建一个map，用于存储生成的UUID，并利用map的键唯一性来检测重复

	numToGenerate := 1000                // 设定要生成的UUID数量
	for i := 0; i < numToGenerate; i++ { // 循环生成指定数量的UUID
		id := uuid.NewV4().String()         // 生成一个UUID字符串
		if _, exists := uuids[id]; exists { // 检查当前生成的UUID是否已经存在于map中
			t.Errorf("发现重复的UUID: %s", id) // 如果发现重复，则报错
		}
		uuids[id] = struct{}{} // 将生成的UUID添加到map中
	}
	t.Logf("成功生成 %d 个唯一UUID", numToGenerate) // 记录成功生成的UUID数量
}

// BenchmarkGenerateUUID_Performance 是一个简单的性能测试函数，用于评估UUID生成的效率。
// 注意：性能测试函数必须以 Benchmark 开头，并且接收 *testing.B 参数。
func BenchmarkGenerateUUID_Performance(b *testing.B) { // 定义一个名为 BenchmarkGenerateUUID_Performance 的基准测试函数，接收 *testing.B 参数
	b.ResetTimer() // 重置计时器，排除Setup代码的执行时间

	b.Run("Sequential Generation", func(b *testing.B) { // 运行一个子基准测试，用于顺序生成性能测试
		for i := 0; i < b.N; i++ { // 循环生成UUID，b.N 会根据测试运行时间自动调整
			_ = uuid.NewV4().String() // 生成UUID，_ 表示忽略返回值
		}
	})

	b.Run("Concurrent Generation (10 goroutines)", func(b *testing.B) { // 运行一个子基准测试，用于并发生成性能测试
		var wg sync.WaitGroup // 声明一个 WaitGroup 用于等待所有goroutine完成
		numGoroutines := 10   // 设定并发的goroutine数量

		if b.N < numGoroutines {
			b.N = numGoroutines // 至少每个goroutine执行一次
		}
		uuidsPerGoroutine := b.N / numGoroutines

		b.ResetTimer()                       // 在并发测试的内部再次重置计时器，确保只测量并发部分的性能
		for i := 0; i < numGoroutines; i++ { // 启动指定数量的goroutine
			wg.Add(1)   // 增加 WaitGroup 计数
			go func() { // 启动一个新的goroutine
				defer wg.Done()                          // goroutine结束后调用 Done() 减少 WaitGroup 计数
				for j := 0; j < uuidsPerGoroutine; j++ { // 每个goroutine内部循环生成UUID
					_ = uuid.NewV4().String() // 生成UUID
				}
			}()
		}
		wg.Wait() // 等待所有goroutine完成
	})
}

// TestSystemUUIDGeneration 是一个系统测试/集成测试示例，模拟高并发下UUID的生成和冲突检测。
// 它可以作为系统测试的一部分，验证在更接近生产环境的负载下的行为。
func TestSystemUUIDGeneration(t *testing.T) { // 定义一个名为 TestSystemUUIDGeneration 的测试函数
	// 注意：系统测试通常会更复杂，可能涉及外部依赖或模拟服务。这里只是一个简化示例。
	t.Log("开始系统测试：高并发UUID生成与冲突检测") // 记录测试开始信息

	numWorkers := 100                         // 并发工作者（goroutine）的数量
	uuidsPerWorker := 10000                   // 每个工作者生成的UUID数量
	totalUUIDs := numWorkers * uuidsPerWorker // 总共要生成的UUID数量

	var mu sync.Mutex                           // 声明一个互斥锁，用于保护共享资源（map）
	generatedUUIDs := make(map[string]struct{}) // 存储所有生成的UUID，用于检测冲突

	var wg sync.WaitGroup // 声明一个 WaitGroup 用于等待所有工作者完成

	for i := 0; i < numWorkers; i++ { // 启动指定数量的并发工作者
		wg.Add(1)               // 增加 WaitGroup 计数
		go func(workerID int) { // 启动一个新的goroutine，传入 workerID
			defer wg.Done()                              // goroutine结束后调用 Done() 减少 WaitGroup 计数
			localUUIDs := make([]string, uuidsPerWorker) // 每个工作者自己的局部UUID切片

			for j := 0; j < uuidsPerWorker; j++ { // 每个工作者生成指定数量的UUID
				localUUIDs[j] = uuid.NewV4().String() // 生成UUID并存储在局部切片中
			}

			// 将局部生成的UUID安全地添加到共享map中
			mu.Lock()                       // 获取互斥锁，保护共享资源
			for _, id := range localUUIDs { // 遍历局部生成的UUID
				if _, exists := generatedUUIDs[id]; exists { // 检查共享map中是否已存在该UUID
					t.Errorf("在高并发下发现重复的UUID: %s (Worker %d)", id, workerID) // 如果发现重复，则报错
				}
				generatedUUIDs[id] = struct{}{} // 将UUID添加到共享map中
			}
			mu.Unlock()                                              // 释放互斥锁
			t.Logf("工作者 %d 完成生成 %d 个UUID", workerID, uuidsPerWorker) // 记录工作者完成信息
		}(i) // 传入当前循环变量 i 作为 workerID
	}

	wg.Wait() // 等待所有工作者完成

	if len(generatedUUIDs) != totalUUIDs { // 检查最终生成的唯一UUID数量是否等于总预期数量
		t.Errorf("系统测试失败：预期生成 %d 个唯一UUID，实际生成 %d 个", totalUUIDs, len(generatedUUIDs)) // 如果不一致，则报错
	}
	t.Logf("系统测试完成：成功生成并验证 %d 个唯一UUID", totalUUIDs) // 记录测试完成信息
}

// TestUUIDFormatValidation 是一个验证测试函数，专门用于严格验证UUID V4的格式。
func TestUUIDFormatValidation(t *testing.T) { // 定义一个名为 TestUUIDFormatValidation 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	numToValidate := 1000                // 设定要验证的UUID数量
	for i := 0; i < numToValidate; i++ { // 循环生成并验证UUID
		id := uuid.NewV4().String()       // 生成一个UUID字符串
		if !uuidV4Regex.MatchString(id) { // 使用正则表达式严格检查UUID格式
			t.Errorf("验证测试失败：生成的UUID '%s' 不符合V4格式。", id) // 如果不符合，则报错
		}
	}
	t.Logf("成功验证 %d 个UUID的V4格式", numToValidate) // 记录成功验证的UUID数量
}

// TestUUIDConsistency 是一个回归测试函数，确保在未来的代码变更后，UUID生成行为保持一致。
// 例如，我们可以确保即使底层库更新，其核心功能（如长度和格式）也不会改变。
func TestUUIDConsistency(t *testing.T) { // 定义一个名为 TestUUIDConsistency 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 在这里，我们可以存储一些“已知良好”的UUID（如果可能的话，但对于随机UUID意义不大），
	// 或者只是重新运行基本断言，确保它们不会因为未来的代码修改而失效。

	// 回归测试通常会重复运行现有的单元测试和集成测试，以确保没有引入新的缺陷。
	// 这里我们重复调用 TestGenerateUUID_BasicAssertions 和 TestGenerateUUID_Uniqueness。

	t.Run("Re-run Basic Assertions", func(t *testing.T) { // 运行子测试，重新执行基本断言
		TestGenerateUUID_BasicAssertions(t) // 调用 TestGenerateUUID_BasicAssertions 函数
	})

	t.Run("Re-run Uniqueness Test", func(t *testing.T) { // 运行子测试，重新执行唯一性测试
		TestGenerateUUID_Uniqueness(t) // 调用 TestGenerateUUID_Uniqueness 函数
	})

	// 模拟将来可能发生的库版本升级或自定义逻辑添加
	// 例如，如果 GetUUID 函数将来增加了额外的逻辑，我们可以在这里添加特定的回归测试
	// 假设未来 GetUUID 可能会被修改为：
	// func GetUUID() string {
	//    baseUUID := uuid.NewV4().String()
	//    return "prefix-" + baseUUID // 模拟添加了一个前缀
	// }
	// 如果发生了这样的变更，则需要新的回归测试来验证新行为。
	// 这里我们只验证当前的简单行为
	sampleID := uuid.NewV4().String()                              // 生成一个样本UUID
	if len(sampleID) != 36 || !uuidV4Regex.MatchString(sampleID) { // 再次检查长度和格式
		t.Errorf("回归测试失败：UUID生成行为不一致。长度或格式不符。") // 如果不符，则报错
	}
	fmt.Printf("生成的UUID: %s, 长度: %d\n", sampleID, len(sampleID)) // 打印生成的UUID和长度
}
