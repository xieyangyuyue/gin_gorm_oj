package service

import (
	"gin_gorm_oj/define"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetTestCase
// @Tags 管理员私有方法
// @Summary 测试案例列表
// @Param authorization header string true "authorization"
// @Param identity query string true "问题唯一标识"
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/test-case [get]
func GetTestCase(c *gin.Context) {
	// 获取URL查询参数中的size，如果未提供则使用define.DefaultSize
	// _ 忽略了可能出现的错误，因为DefaultQuery已经处理了默认值
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Printf("分页参数size转换错误: %v\n", err) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分页参数size错误：" + err.Error(),
		})
		return
	}

	// 获取URL查询参数中的page，如果未提供则使用define.DefaultPage
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Printf("分页参数page转换错误: %v\n", err) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分页参数page错误：" + err.Error(),
		})
		return
	}

	// 计算数据库查询的偏移量
	// 例如：page=1, size=10 -> offset=0
	// 例如：page=2, size=10 -> offset=10
	page = (page - 1) * size

	// 获取URL查询参数中的problemIdentity（问题唯一标识）
	problemIdentity := c.Query("identity")
	// 检查problemIdentity是否为空
	if problemIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}

	// 用于存储测试案例总数
	var count int64
	// 创建TestCase切片用于存储查询结果
	data := make([]*models.TestCase, 0)

	// 构建数据库查询，首先筛选出特定problem_identity的测试案例，并计算总数
	tx := models.DB.Model(new(models.TestCase)).Where("problem_identity = ?", problemIdentity).Count(&count)

	// 如果查询参数中提供了size（意味着需要分页），则应用分页逻辑
	// 修正：Limit参数应该使用size，而不是page
	if c.Query("size") != "" {
		tx = tx.Offset(page).Limit(size) // 应用偏移量和限制数量
	}

	// 执行查询，将结果存储到data切片中
	err = tx.Find(&data).Error
	// 检查数据库查询是否出错
	if err != nil {
		log.Printf("获取测试案例列表失败: %v, problemIdentity: %s\n", err, problemIdentity) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取测试案例列表失败：" + err.Error(),
		})
		return
	}

	// 返回查询到的测试案例列表和总数
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  data,  // 测试案例列表
			"count": count, // 总测试案例数
		},
	})
}
