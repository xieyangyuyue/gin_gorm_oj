package models // 声明包名为 models，通常用于存放数据模型和数据库操作相关代码

import (
	"fmt"                          // 导入 fmt 包，用于格式化字符串
	"gin_gorm_oj/utils"            // 导入自定义的 utils 包，用于获取数据库和 Redis 配置信息。请确保路径正确
	"github.com/go-redis/redis/v8" // 导入 go-redis/redis v8 库，用于 Redis 客户端操作
	"gorm.io/driver/mysql"         // 导入 GORM 的 MySQL 驱动
	"gorm.io/gorm"                 // 导入 GORM 核心库
	"gorm.io/gorm/logger"          // 导入 GORM 的日志模块
	"gorm.io/gorm/schema"          // 导入 GORM 的 schema 模块，用于命名策略
	"log"                          // 导入 log 包，用于日志输出
	"time"                         // 导入 time 包，用于时间相关的操作
)

// DB 是全局的 GORM 数据库实例，通过 Init 函数初始化。
// 外部模块可以通过 models.DB 访问数据库连接。
var DB = Init() // 在包初始化时调用 Init 函数来获取 GORM 数据库实例

// RDB 是全局的 Redis 数据库实例，通过 InitRedisDB 函数初始化。
// 外部模块可以通过 models.RDB 访问 Redis 连接。
var RDB = InitRedisDB() // 在包初始化时调用 InitRedisDB 函数来获取 Redis 客户端实例

// Init 初始化 GORM 数据库连接。
// 使用 utils 包中定义的数据库连接信息来构建 DSN。
//
// 返回值:
//
//	*gorm.DB: 返回一个 GORM 数据库连接实例。如果在连接过程中发生严重错误，会导致程序退出。
func Init() *gorm.DB { //
	// 构造数据库连接的 DSN (Data Source Name) 字符串。
	// DSN 包含了连接 MySQL 数据库所需的所有参数。
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", //
		utils.DbUser,     // 数据库用户名
		utils.DbPassWord, // 数据库密码
		utils.DbHost,     // 数据库主机地址
		utils.DbPort,     // 数据库端口
		utils.DbName,     // 数据库名称
	)
	// 打开数据库连接。
	// mysql.Open(dns) 使用 MySQL 驱动打开连接。
	// &gorm.Config{} 配置 GORM 的行为。
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{ // gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Silent), // 配置 GORM 日志模式为静默 (Silent)，生产环境推荐，可根据需要调整为 Info 或 Warn
		// 外键约束：在进行数据库迁移时禁用外键约束，可以避免循环依赖问题。
		DisableForeignKeyConstraintWhenMigrating: true, //
		// 禁用默认事务：提高某些操作的性能，但需在业务逻辑中手动管理事务。
		SkipDefaultTransaction: true, //
		// 命名策略：配置数据库表名和字段名的转换规则。
		NamingStrategy: schema.NamingStrategy{ //
			// 使用单数表名：启用该选项，此时，`User` 结构体对应的表名应该是 `user`。
			SingularTable: true, //
		},
	})
	if err != nil { // 检查数据库连接是否出错
		// 记录数据库初始化错误信息，并在错误发生时终止程序，因为数据库是核心依赖。
		log.Fatalf("GORM Init Error: %v", err) // 使用 Fatal 终止程序，而不是 Println，因为连接失败是致命错误
	}
	// 迁移数据表。
	// 注意: 数据库表结构初次运行后可注销此行，或仅在开发环境开启，生产环境慎用或手动迁移。
	// _ = db.AutoMigrate(&CategoryBasic{}, &ContestBasic{}, &ContestProblem{},
	// 	&ContestUser{}, &ProblemBasic{}, &ProblemCategory{}, &SubmitBasic{}, &TestCase{}, &UserBasic{})

	// 获取底层的 *sql.DB 实例，用于配置连接池参数。
	sqlDB, err := db.DB()
	if err != nil { // 检查获取 sql.DB 实例是否出错
		log.Fatalf("Failed to get underlying *sql.DB: %v", err) // 获取失败也是致命错误
	}

	// SetMaxIdleConns 设置连接池中的最大闲置连接数。
	// 建议根据应用负载和数据库性能进行调优。
	sqlDB.SetMaxIdleConns(10) //

	// SetMaxOpenConns 设置数据库的最大连接数量。
	// 建议根据数据库服务器的最大连接数和应用需求进行调优。
	sqlDB.SetMaxOpenConns(100) //

	// SetConnMaxLifetime 设置连接的最大可复用时间。
	// 连接在达到此时间后会被关闭，确保连接的刷新，防止因连接时间过长导致的问题。
	sqlDB.SetConnMaxLifetime(10 * time.Second) //

	return db // 返回初始化后的 GORM 数据库实例
}

// InitRedisDB 初始化 Redis 数据库连接。
// 使用 utils 包中定义的 Redis 连接信息。
//
// 返回值:
//
//	*redis.Client: 返回一个 Redis 客户端实例。
func InitRedisDB() *redis.Client { //
	return redis.NewClient(&redis.Options{ //
		Addr:     fmt.Sprintf("%s:%s", utils.RedisHost, utils.RedisPort), // Redis 服务器地址，使用格式化拼接确保有冒号
		Password: utils.RedisPassWord,                                    // Redis 认证密码
		DB:       utils.RedisNumber,                                      // 选择使用的数据库编号（0-15，默认0）
	})
}
