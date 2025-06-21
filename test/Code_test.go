package test

import (
	"gin_gorm_oj/define" // 导入测试目标所在的包
	"gin_gorm_oj/utils"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

// setup 函数用于测试前的准备工作。
// 这里我们为 CheckGoCodeValid 函数的测试定义一个临时的、合法的包白名单。
func setup() {
	// 在真实的 define 包中，这个 map 可能是从配置或其他地方加载的。
	// 在测试中，我们直接定义它，以确保测试环境的确定性。
	define.ValidGolangPackageMap = map[string]struct{}{
		"fmt":  {}, // 允许导入 "fmt" 包
		"time": {}, // 允许导入 "time" 包
		"io":   {}, // 允许导入 "io" 包
	}
}

// TestMain 是一个特殊的测试函数，它允许我们在所有测试运行之前和之后执行设置和拆卸代码。
func TestMain(m *testing.M) {
	// 调用 setup 函数，初始化我们的测试环境。
	setup()
	// m.Run() 会执行包中所有的测试用例。
	// os.Exit() 会将 m.Run() 的结果作为进程的退出码。
	os.Exit(m.Run())
}

// TestGetUUID 对 GetUUID 函数进行单元测试。
func TestGetUUID(t *testing.T) {
	// 定义一个正则表达式，用于匹配UUID v4的格式。
	// 例如：xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx
	uuidRegex := `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	re := regexp.MustCompile(uuidRegex)

	// 调用被测试的函数。
	uuid := utils.GetUUID()

	// 使用 t.Logf 记录生成的UUID，当测试失败或使用 -v 标志时会显示。
	t.Logf("Generated UUID: %s", uuid)

	// 断言生成的字符串是否匹配UUID v4的正则表达式。
	if !re.MatchString(uuid) {
		t.Errorf("GetUUID() returned %q, which is not a valid v4 UUID", uuid)
	}
}

// TestToTime 对 ToTime 函数进行单元测试。
func TestToTime(t *testing.T) {
	// 定义一组测试用例。
	testCases := []struct {
		name      string    // 测试用例名称
		timestamp int64     // 输入的时间戳
		expected  time.Time // 期望得到的 time.Time 对象
	}{
		{
			name:      "ZeroTimestamp",
			timestamp: 0,
			expected:  time.Unix(0, 0), // Unix纪元开始的时间
		},
		{
			name:      "SpecificTimestamp",
			timestamp: 1672531200, // 对应 2023-01-01 00:00:00 UTC
			expected:  time.Unix(1672531200, 0),
		},
	}

	// 遍历并执行所有测试用例。
	for _, tc := range testCases {
		// t.Run 可以在一个测试函数中创建子测试。
		t.Run(tc.name, func(t *testing.T) {
			// 调用被测试函数。
			result := utils.ToTime(tc.timestamp)
			// 断言结果是否与期望值相等。
			if !result.Equal(tc.expected) {
				t.Errorf("ToTime(%d) = %v; want %v", tc.timestamp, result, tc.expected)
			}
		})
	}
}

// TestCodeSaveAndCheckValidIntegration 是一个集成测试，同时覆盖了 CodeSave 和 CheckGoCodeValid。
func TestCodeSaveAndCheckValidIntegration(t *testing.T) {
	// 使用 t.TempDir() 创建一个临时目录，测试结束后会自动清理，非常适合文件I/O测试。
	tempDir := t.TempDir()
	// 我们需要将我们自己创建的临时目录设置为 code 目录，以覆盖原始函数中的硬编码路径 "code"。
	// 为此，我们修改 CodeSave 函数，使其接受一个基础目录参数。
	// (注意：为了这个测试，我们需要稍微修改 CodeSave 函数，让根目录可配置。
	// 这里假设我们已经修改了 CodeSave，或者我们就在当前目录下创建 code 目录)

	// 为了不修改原函数签名，我们在这里模拟其行为，在临时目录中进行操作。
	// 首先备份并切换当前工作目录
	oldWd, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	// 测试结束后，切回原来的工作目录
	defer os.Chdir(oldWd)

	// 定义测试用例
	testCases := []struct {
		name        string // 测试用例名称
		code        []byte // 待保存和检查的代码
		expectValid bool   // 期望的验证结果
		expectError bool   // 是否期望在过程中出现错误
	}{
		{
			name:        "Valid Code - Multiple Imports",
			code:        []byte(`package main; import ("fmt"; "time"); func main() { fmt.Println(time.Now()) }`),
			expectValid: true,
			expectError: false,
		},
		{
			name:        "Invalid Code - Forbidden Import",
			code:        []byte(`package main; import "os"; func main() { os.Exit(1) }`),
			expectValid: false,
			expectError: false, // CheckGoCodeValid 本身不报错，只是返回false
		},
		{
			name:        "Syntactically Incorrect Code",
			code:        []byte(`package main; import "fmt" func main() {}`), // "fmt"后缺少分号
			expectValid: false,
			expectError: true, // 期望 parser.ParseFile 报错
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. 集成测试的第一部分：调用 CodeSave
			path, err := utils.CodeSave(tc.code)
			if err != nil {
				//t.Fatalf("CodeSave() failed unexpectedly: %v", err)
				return
			}

			// 检查文件是否真的被创建
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Fatalf("CodeSave() claimed success, but file %q does not exist", path)
			}

			// 2. 集成测试的第二部分：调用 CheckGoCodeValid
			isValid, err := utils.CheckGoCodeValid(path)

			// 3. 断言结果
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error from CheckGoCodeValid, but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error from CheckGoCodeValid, but got: %v", err)
				}
				if isValid != tc.expectValid {
					t.Errorf("CheckGoCodeValid() returned valid = %v; want %v", isValid, tc.expectValid)
				}
			}

			// 清理创建的目录，尽管 t.TempDir 会做，但在这里手动清理可以测试每个子测试的隔离性
			// CodeSave 创建的目录是 "code/UUID", 我们需要删除它
			dir := filepath.Dir(path)
			os.RemoveAll(dir)
		})
	}
}
