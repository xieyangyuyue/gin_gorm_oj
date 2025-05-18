package models

import (
	"fmt"
	"gin_gorm_oj/helper"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

// DB 是全局的 GORM 数据库实例，通过 Init 函数初始化
var DB = Init()

// RDB 是全局的 Redis 数据库实例，通过 InitRedisDB 函数初始化
var RDB = InitRedisDB()

// Init 初始化 GORM 数据库连接
// 使用 define.MysqlDNS 定义的数据库连接信息
func Init() *gorm.DB {
	// 构造数据库连接的 DSN 字符串
	// 构建DSN（数据源名称）
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		helper.DbUser,
		helper.DbPassWord,
		helper.DbHost,
		helper.DbPort,
		helper.DbName,
	)
	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{ // gorm日志模式：silent
		Logger: logger.Default.LogMode(logger.Silent),
		// 外键约束
		DisableForeignKeyConstraintWhenMigrating: true,
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// 使用单数表名，启用该选项，此时，`User` 的表名应该是 `user`
			SingularTable: true,
		}})
	if err != nil {
		// 记录数据库初始化错误信息
		log.Println("gorm Init Error : ", err)
	}
	// 迁移数据表，在没有数据表结构变更时候，建议注释不执行
	// 注意:初次运行后可注销此行
	//_ = db.AutoMigrate(&CategoryBasic{}, &ContestBasic{}, &ContestProblem{},
	//	&ContestUser{}, &ProblemBasic{}, &ProblemCategory{}, &SubmitBasic{}, &TestCase{}, &UserBasic{})

	sqlDB, _ := db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	return db
}

// InitRedisDB 初始化 Redis 数据库连接
// 使用本地 Redis 服务器，端口为 6379，无密码，使用默认数据库
func InitRedisDB() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
