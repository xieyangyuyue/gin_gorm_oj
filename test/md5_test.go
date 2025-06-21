package test // 声明包名为 utils_test，这是 Go 语言中进行外部测试（不在原包内）的约定

import (
	"fmt"     // 导入 fmt 包，用于格式化输出
	"regexp"  // 导入 regexp 包，用于正则表达式匹配
	"sync"    // 导入 sync 包，用于并发测试中的同步原语，如 WaitGroup
	"testing" // 导入 Go 语言内置的 testing 包，用于编写测试函数

	"gin_gorm_oj/utils" // 导入你的 utils 包，请将 "your_module_path" 替换为你的实际模块路径
)

// md5Regex 是一个正则表达式，用于验证MD5哈希值的格式。
// MD5哈希值是一个32位的十六进制字符串（16字节 * 2位/字节 = 32位）。
var md5Regex = regexp.MustCompile("^[0-9a-fA-F]{32}$")

// TestGetMd5_Unit 是一个单元测试函数，用于验证 GetMd5 函数的基本功能和正确性。
func TestGetMd5_Unit(t *testing.T) { // 定义一个名为 TestGetMd5_Unit 的测试函数
	t.Parallel() // 标记该测试可以并行运行，提高测试效率

	// 定义一组测试用例，包含输入字符串和期望的MD5哈希值
	tests := []struct {
		input    string // 输入字符串
		expected string // 期望的MD5哈希值（小写十六进制）
	}{
		{"", "d41d8cd98f00b204e9800998ecf8427e"},                                            // 空字符串的MD5值
		{"hello", "5d41402abc4b2a76b9719d911017c592"},                                       // 常用字符串的MD5值
		{"Hello World", "b10a8db164e0754105b7a99be72e3fe5"},                                 // 包含空格和大小写的字符串
		{"1234567890", "e807f1fcf82d132f9bb018ca6738a19f"},                                  // 数字字符串
		{"The quick brown fox jumps over the lazy dog", "9e107d9d372bb6826bd81d3542a419d6"}, // 经典语句的MD5值
		{"你好，世界！", "5082079d92a8ef985f59e001d445ff20"},                                      // 包含中文的字符串
		{"!@#$%^&*()", "05b28d17a7b6e7024b6e5d8cc43a8bf7"},                                  // 包含特殊字符的字符串
	}

	for _, tt := range tests { // 遍历所有测试用例
		t.Run(fmt.Sprintf("Input:%q", tt.input), func(t *testing.T) { // 为每个测试用例创建一个子测试
			actual := utils.GetMd5(tt.input) // 调用 GetMd5 函数计算实际MD5值

			// 单元测试 - 结果对比
			if actual != tt.expected { // 比较实际结果和期望结果
				t.Errorf("GetMd5(%q) 期望结果为 %q, 实际得到 %q", tt.input, tt.expected, actual) // 如果不匹配，则报错
			}

			// 单元测试 - 长度检查
			expectedLength := 32               // MD5哈希值的标准长度是32个十六进制字符
			if len(actual) != expectedLength { // 检查生成的MD5字符串长度是否符合预期
				t.Errorf("生成的MD5哈希值长度不正确，期望 %d，得到 %d", expectedLength, len(actual)) // 如果长度不符，则报错
			}

			// 单元测试 - 格式检查（使用正则表达式）
			if !md5Regex.MatchString(actual) { // 使用预定义的正则表达式检查MD5字符串是否符合十六进制格式
				t.Errorf("生成的MD5哈希值格式不正确: %s", actual) // 如果格式不符，则报错
			}
		})
	}
}

// TestGetMd5_EdgeCases 是一个单元测试函数，用于测试 GetMd5 函数的边缘情况。
func TestGetMd5_EdgeCases(t *testing.T) { // 定义一个名为 TestGetMd5_EdgeCases 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 极端长度的字符串测试
	t.Run("Very Long String", func(t *testing.T) { // 测试非常长的字符串
		longString := ""
		for i := 0; i < 10000; i++ { // 创建一个包含10000个字符的字符串
			longString += "a"
		}
		expectedMD5 := "0d0c9c4db6953fee9e03f528cafd7d3e" // 预期的长字符串MD5值
		actual := utils.GetMd5(longString)
		if actual != expectedMD5 { // 检查MD5值是否正确
			t.Errorf("GetMd5 for long string 期望 %q, 实际 %q", expectedMD5, actual) // 如果不匹配，则报错
		}
	})

	t.Run("String with Null Byte", func(t *testing.T) { // 测试包含空字节的字符串
		input := "hello\x00world"                         // 包含空字节的字符串
		expectedMD5 := "838d3870873a75639041ff8940f397db" // 预期的MD5值
		actual := utils.GetMd5(input)
		if actual != expectedMD5 { // 检查MD5值是否正确
			t.Errorf("GetMd5 for string with null byte 期望 %q, 实际 %q", expectedMD5, actual) // 如果不匹配，则报错
		}
	})
}

// BenchmarkGetMd5_Performance 是一个基准测试函数，用于评估 GetMd5 函数的性能。
func BenchmarkGetMd5_Performance(b *testing.B) { // 定义一个名为 BenchmarkGetMd5_Performance 的基准测试函数
	inputString := "this is a test string for md5 performance benchmark" // 用于测试的输入字符串
	b.ResetTimer()                                                       // 重置计时器，排除Setup代码的执行时间

	b.Run("Sequential Generation", func(b *testing.B) { // 运行一个子基准测试，用于顺序计算MD5的性能
		for i := 0; i < b.N; i++ { // 循环 b.N 次，b.N 会根据测试运行时间自动调整
			_ = utils.GetMd5(inputString) // 调用 GetMd5 函数，_ 表示忽略返回值
		}
	})

	b.Run("Concurrent Generation (10 goroutines)", func(b *testing.B) { // 运行一个子基准测试，用于并发计算MD5的性能
		var wg sync.WaitGroup // 声明一个 WaitGroup 用于等待所有goroutine完成
		numGoroutines := 10   // 设定并发的goroutine数量

		// 确保 b.N 至少能被 numGoroutines 整除，避免除数为零或分配不均的问题
		if b.N < numGoroutines {
			b.N = numGoroutines // 至少每个goroutine执行一次
		}
		iterationsPerGoroutine := b.N / numGoroutines // 每个goroutine需要计算的MD5次数

		b.ResetTimer()                       // 在并发测试的内部再次重置计时器，确保只测量并发部分的性能
		for i := 0; i < numGoroutines; i++ { // 启动指定数量的goroutine
			wg.Add(1)   // 增加 WaitGroup 计数
			go func() { // 启动一个新的goroutine
				defer wg.Done()                               // goroutine结束后调用 Done() 减少 WaitGroup 计数
				for j := 0; j < iterationsPerGoroutine; j++ { // 每个goroutine内部循环计算MD5
					_ = utils.GetMd5(inputString) // 调用 GetMd5 函数
				}
			}()
		}
		wg.Wait() // 等待所有goroutine完成
	})
}

// TestGetMd5_SystemIntegration 是一个系统测试/集成测试示例。
// 在这个特定的MD5函数场景下，它不直接与其他复杂系统集成，
// 但我们可以模拟它在一个更大型的应用中使用的情况，例如计算文件内容的MD5。
func TestGetMd5_SystemIntegration(t *testing.T) {
	// 定义一个名为 TestGetMd5_SystemIntegration 的测试函数
	// 注意：真正的系统集成测试会更复杂，可能涉及文件读写、网络IO或数据库操作。
	// 这里只是一个简化的模拟，用于展示概念。

	t.Log("开始系统集成测试：模拟文件内容MD5计算") // 记录测试开始信息

	// 模拟文件内容
	fileContent1 := "This is the content of file one."
	fileContent2 := "This is the content of file two. It's different."
	fileContent3 := "This is the content of file one." // 与 fileContent1 相同

	md51 := utils.GetMd5(fileContent1) // 计算文件内容1的MD5值
	md52 := utils.GetMd5(fileContent2) // 计算文件内容2的MD5值
	md53 := utils.GetMd5(fileContent3) // 计算文件内容3的MD5值

	// 验证相同内容是否生成相同的MD5
	if md51 != md53 { // 检查相同内容是否生成相同的MD5
		t.Errorf("集成测试失败：相同文件内容生成了不同的MD5值。MD5_1: %s, MD5_3: %s", md51, md53) // 如果不一致，则报错
	} else {
		t.Logf("集成测试：相同内容 '%s' 生成了相同的MD5值 '%s' (符合预期)", fileContent1, md51) // 记录符合预期
	}

	// 验证不同内容是否生成不同的MD5
	if md51 == md52 { // 检查不同内容是否生成不同的MD5
		t.Errorf("集成测试失败：不同文件内容生成了相同的MD5值。MD5_1: %s, MD5_2: %s", md51, md52) // 如果一致，则报错
	} else {
		t.Logf("集成测试：不同内容生成了不同的MD5值 (符合预期)") // 记录符合预期
	}

	t.Log("系统集成测试完成。") // 记录测试完成信息
}

// TestGetMd5_Validation 是一个验证测试函数，专门用于严格验证MD5哈希值的格式和特性。
func TestGetMd5_Validation(t *testing.T) { // 定义一个名为 TestGetMd5_Validation 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 验证MD5哈希值的长度始终为32
	t.Run("Length Validation", func(t *testing.T) { // 运行子测试，验证长度
		input := "any string"
		md5Hash := utils.GetMd5(input)
		if len(md5Hash) != 32 { // 检查MD5哈希值的长度
			t.Errorf("验证失败：MD5哈希值长度不为32，实际为 %d", len(md5Hash)) // 如果不为32，则报错
		}
	})

	// 验证MD5哈希值只包含十六进制字符
	t.Run("Hexadecimal Format Validation", func(t *testing.T) { // 运行子测试，验证十六进制格式
		input := "another string"
		md5Hash := utils.GetMd5(input)
		if !md5Regex.MatchString(md5Hash) { // 使用正则表达式检查MD5哈希值是否只包含十六进制字符
			t.Errorf("验证失败：MD5哈希值 '%s' 包含非十六进制字符或格式不符。", md5Hash) // 如果不符，则报错
		}
	})

	// 验证MD5的确定性：相同输入总是生成相同输出
	t.Run("Determinism Validation", func(t *testing.T) { // 运行子测试，验证确定性
		input := "deterministic test"
		hash1 := utils.GetMd5(input)
		hash2 := utils.GetMd5(input) // 再次计算相同输入的MD5
		if hash1 != hash2 {          // 检查两次计算结果是否一致
			t.Errorf("验证失败：MD5哈希不是确定性的，hash1: %s, hash2: %s", hash1, hash2) // 如果不一致，则报错
		}
	})
}

// TestGetMd5_Regression 是一个回归测试函数，确保在未来的代码变更后，GetMd5 的行为仍然一致且正确。
func TestGetMd5_Regression(t *testing.T) { // 定义一个名为 TestGetMd5_Regression 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 回归测试通常会重复运行现有的单元测试和集成测试，以确保没有引入新的缺陷。
	// 这里我们重复调用 TestGetMd5_Unit 和 TestGetMd5_SystemIntegration 的部分逻辑。

	t.Run("Re-run Basic Unit Cases", func(t *testing.T) { // 运行子测试，重新执行基本单元测试用例
		// 可以选择性地运行 TestGetMd5_Unit 中的部分关键用例，而不是全部重新运行
		input := "regression test string"
		expected := "a59f2df1940bc5d255272bfeaac1620f" // 预期的MD5值
		actual := utils.GetMd5(input)
		if actual != expected { // 检查MD5值是否符合预期
			t.Errorf("回归测试失败：GetMd5(%q) 期望 %q, 实际 %q", input, expected, actual) // 如果不符，则报错
		}
		if len(actual) != 32 || !md5Regex.MatchString(actual) { // 再次检查长度和格式
			t.Errorf("回归测试失败：MD5哈希值格式或长度不一致。") // 如果不符，则报错
		}
	})

	t.Run("Re-check Consistency after Potential Changes", func(t *testing.T) { // 运行子测试，检查潜在变更后的一致性
		// 模拟一个关键的、未来可能被修改的输入
		criticalInput := "password123"
		knownGoodHash := "482c811da5d5b4bc6d497ffa98491e38" // 这个值在代码修改后也不应该变
		currentHash := utils.GetMd5(criticalInput)

		if currentHash != knownGoodHash { // 检查关键输入的MD5值是否保持不变
			t.Errorf("回归测试失败：关键输入MD5哈希值发生变化！期望 %q, 实际 %q", knownGoodHash, currentHash) // 如果变化，则报错
		}
	})
}
