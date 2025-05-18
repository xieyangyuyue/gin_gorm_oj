package define

import (
	"os"
)

// DefaultPage 是默认的分页页码
var DefaultPage = "1"

// DefaultSize 是默认的每页显示数量
var DefaultSize = "20"

// MailPassword 从环境变量中获取邮件密码
var MailPassword = os.Getenv("MailPassword")

// MysqlDNS 是 MySQL 数据库的连接字符串
var MysqlDNS = "xieyangyuyue:xieyangyuyue@tcp(111.170.11.164:10500)"

// ProblemBasic 表示问题基础信息的结构体
type ProblemBasic struct {
	// Identity 是问题表的唯一标识
	Identity string `json:"identity"`
	// Title 是问题标题
	Title string `json:"title"`
	// Content 是问题内容
	Content string `json:"content"`
	// ProblemCategories 是关联问题分类表的 ID 列表
	ProblemCategories []int `json:"problem_categories"`
	// MaxRuntime 是最大运行时长
	MaxRuntime int `json:"max_runtime"`
	// MaxMem 是最大运行内存
	MaxMem int `json:"max_mem"`
	// TestCases 是关联测试用例表的列表
	TestCases []*TestCase `json:"test_cases"`
}

// TestCase 表示测试用例的结构体
type TestCase struct {
	// Input 是测试用例的输入
	Input string `json:"input"`
	// Output 是测试用例的输出
	Output string `json:"output"`
}

// ContestBasic 表示竞赛基础信息的结构体
type ContestBasic struct {
	// Identity 是竞赛的唯一标识
	Identity string `json:"identity"`
	// Name 是竞赛名称
	Name string `json:"name"`
	// Content 是竞赛描述
	Content string `json:"content"`
	// ProblemBasics 是关联题目表的 ID 列表
	ProblemBasics []int `json:"problem_basic"`
	// StartAt 是竞赛开启时间
	StartAt int64 `json:"start_at"`
	// EndAt 是竞赛关闭时间
	EndAt int64 `json:"end_at"`
}

// DateLayout 是日期时间的格式化布局
var DateLayout = "2006-01-02 15:04:05"

// ValidGolangPackageMap 是有效的 Go 语言包名映射
var ValidGolangPackageMap = map[string]struct{}{
	"bytes":   {},
	"fmt":     {},
	"math":    {},
	"sort":    {},
	"strings": {},
}
