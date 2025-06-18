// 本测试包用于验证Redis客户端操作及models包中的Redis配置
package test

import (
	"context" // 上下文管理，用于控制请求生命周期（如超时）
	"fmt"     // 格式化输出，用于测试结果展示
	"gin_gorm_oj/models"
	"gin_gorm_oj/utils"
	"github.com/go-redis/redis/v8" // Redis客户端库
	"testing"                      // Go测试框架，提供测试功能
	"time"                         // 时间处理，用于设置键值过期时间
)

// 全局上下文对象
// 在Redis操作中传递上下文，可用于实现超时控制、取消操作等（当前使用默认空上下文）
var ctx = context.Background()

// Redis客户端实例
// 使用NewClient方法创建连接池，配置参数说明：
// Addr:     Redis服务器地址（默认本地6379端口）
// Password: 认证密码，此处需根据实际环境修改
// DB:       选择使用的数据库编号（0-15，默认0）
var rdb = redis.NewClient(&redis.Options{
	Addr:     fmt.Sprintf("%s:%s", utils.RedisHost, utils.RedisPort), // 使用格式化拼接确保有冒号
	Password: utils.RedisPassWord,
	DB:       utils.RedisNumber,
})

// TestRedisSet 测试Redis的字符串写入功能
// 功能：设置键值对，验证基础写入功能与过期时间配置
// 测试逻辑：
// 1. 使用Set方法写入键"name"，值"mmc"
// 2. 设置30秒过期时间（注意：实际代码中应为30*time.Second）
// 3. 错误检查：失败时终止测试
// 注意：参数应使用*testing.T指针（当前代码存在类型错误）
func TestRedisSet(t *testing.T) {
	err := rdb.Set(ctx, "name", "mmc", 60*time.Second).Err()
	if err != nil {
		t.Fatal("Redis Set操作失败:", err) // 严重错误，立即终止
	}
}

// TestRedisGet 测试Redis的字符串读取功能
// 功能：验证TestRedisSet写入的数据一致性
// 测试逻辑：
// 1. 使用Get方法获取键"name"的值
// 2. Result()方法同时返回值和错误信息
// 3. 错误检查：读取失败时终止测试
// 4. 成功时打印结果，预期输出"mmc"
func TestRedisGet(t *testing.T) {
	v, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal("Redis Get操作失败:", err)
	}
	fmt.Println("获取到的值:", v) // 控制台输出验证（测试成功时）
}

// TestRedisGetByModels 验证models包的Redis配置
// 功能：确保项目内其他模块通过统一配置能正常访问Redis
// 测试逻辑：
// 1. 使用models包导出的RDB实例进行操作
// 2. 读取同一键值，验证配置一致性
// 3. 错误检查：失败说明models包配置错误
// 4. 成功时打印结果，预期与TestRedisGet结果一致
func TestRedisGetByModels(t *testing.T) {
	v, err := models.RDB.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal("通过models.RDB获取数据失败:", err)
	}
	fmt.Println("通过models包获取的值:", v)
}
