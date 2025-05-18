package models

import (
	"fmt"
	"gin_gorm_oj/helper"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
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
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		// 记录数据库初始化错误信息
		log.Println("gorm Init Error : ", err)
	}
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
