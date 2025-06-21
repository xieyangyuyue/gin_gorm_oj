package test

import (
	"gin_gorm_oj/models" // 导入自定义的 models 包，其中包含 DB 实例和数据模型。请确保路径正确
	"testing"            // 导入 Go 语言内置的 testing 包，用于编写测试函数
)

// TestGormConnectionAndDataRetrieval 是一个集成测试函数，
// 用于验证 GORM 数据库连接的正确性以及从数据库中检索数据的能力。
// 它依赖于 models 包中已初始化的 DB 实例。
func TestGormConnectionAndDataRetrieval(t *testing.T) { // 定义一个名为 TestGormConnectionAndDataRetrieval 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 直接使用 models 包中已初始化的全局 DB 实例。
	// 这避免了在测试中重复数据库连接逻辑，使测试更专注于数据操作的验证。
	db := models.DB // 获取 models 包中已初始化好的 GORM DB 实例
	if db == nil {  // 检查 DB 实例是否为空 (理论上不应该，除非 models.Init() 失败)
		t.Fatal("models.DB is nil, database connection might have failed to initialize.") // 如果 DB 为空，则致命错误并退出测试
	}

	// 声明一个切片，用于存储从 ProblemBasic 表中查询到的数据。
	// ProblemBasic 应该是一个 GORM 模型结构体，定义在 models 包中。
	data := make([]*models.ProblemBasic, 0) // 创建一个 ProblemBasic 结构体指针切片

	// 执行数据库查询操作：从 ProblemBasic 表中查找所有记录。
	// .Error 用于检查查询过程中是否发生错误。
	err := db.Find(&data).Error // 查询 ProblemBasic 表中的所有数据
	if err != nil {             // 检查查询是否出错
		t.Fatalf("Failed to query ProblemBasic data: %v", err) // 如果查询失败，则致命错误并退出测试
	}

	// 验证查询结果：检查是否至少查询到了一条数据。
	// 如果数据库为空，这个测试可能需要调整预期或插入测试数据。
	if len(data) == 0 { // 检查查询结果是否为空
		t.Log("No ProblemBasic data found in the database. Consider populating test data.") // 如果没有数据，则记录日志提示
		// t.Error("No ProblemBasic data found in the database. Test expects at least one record.") // 如果测试要求必须有数据，则使用 Error
	}

	// 遍历查询到的数据，并打印每一条记录。
	// 仅用于调试和查看数据，不作为断言依据。
	for _, v := range data { // 遍历查询到的 ProblemBasic 记录
		t.Logf("Problem ==> %v", v) // 使用 t.Log 打印数据，而不是 fmt.Printf，以便与测试日志集成
	}

	t.Log("GORM database connection and data retrieval test completed successfully.") // 记录测试成功信息
}

// TestGorm_Regression 是一个回归测试函数，用于确保 GORM 相关操作在代码变更后仍然保持预期行为。
func TestGorm_Regression(t *testing.T) { // 定义一个名为 TestGorm_Regression 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	t.Run("DB Instance Availability", func(t *testing.T) { // 验证 DB 实例的可用性
		// 确保 models.DB 实例在任何时候都能够被访问并且不是 nil
		if models.DB == nil { // 检查 models.DB 是否为 nil
			t.Errorf("Regression test failed: models.DB instance is nil.") // 如果是 nil，则报错
		}
		t.Log("Regression test: models.DB instance is available.") // 记录可用性
	})

	t.Run("QueryWithoutError", func(t *testing.T) { // 验证基本查询操作无错误
		// 再次执行一个简单的查询，确保没有引入导致查询失败的回归问题
		var count int64
		err := models.DB.Model(&models.ProblemBasic{}).Count(&count).Error // 查询 ProblemBasic 表的记录数
		if err != nil {                                                    // 检查查询是否出错
			t.Errorf("Regression test failed: Basic count query returned error: %v", err) // 如果出错，则报错
		}
		t.Logf("Regression test: Basic count query successful. Count: %d", count) // 记录成功
	})
	t.Log("GORM 回归测试完成。") // 记录回归测试完成信息
}

// TestGorm_SystemHealthCheck 是一个系统测试级别的检查，模拟应用启动时对数据库连接的健康检查。
func TestGorm_SystemHealthCheck(t *testing.T) { // 定义一个名为 TestGorm_SystemHealthCheck 的测试函数
	t.Log("开始系统健康检查：GORM 数据库连接") // 记录测试开始信息

	// 模拟应用启动时的健康检查，尝试 ping 数据库
	sqlDB, err := models.DB.DB() // 获取底层的 *sql.DB 实例
	if err != nil {              // 检查获取 sql.DB 实例是否出错
		t.Fatalf("System Health Check failed: Could not get underlying *sql.DB: %v", err) // 如果失败，则致命错误
	}

	// Ping 数据库以验证连接是否活跃
	err = sqlDB.Ping() // Ping 数据库
	if err != nil {    // 检查 Ping 是否出错
		t.Fatalf("System Health Check failed: Database connection ping failed: %v", err) // 如果失败，则致命错误
	}

	t.Log("系统健康检查完成：GORM 数据库连接正常。") // 记录健康检查成功信息
}
