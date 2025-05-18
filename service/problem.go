package service

import (
	"gin_gorm_oj/define"
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"  // 分页页码，可选参数
// @Param size query int false "size"  // 每页显示数量，可选参数
// @Param keyword query string false "keyword"  // 查询关键字，可选参数
// @Param category_identity query string false "category_identity"  // 分类标识，可选参数
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	// 从请求中获取每页显示数量，若未提供则使用默认值
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	// 从请求中获取页码，若未提供则使用默认值
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		// 若页码转换出错，记录错误日志
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}
	// 计算偏移量
	page = (page - 1) * size
	// 用于存储查询到的记录总数
	var count int64
	// 从请求中获取查询关键字
	keyword := c.Query("keyword")
	// 从请求中获取分类标识
	//categoryIdentity := c.Query("category_identity")

	// 初始化问题列表
	list := make([]*models.ProblemBasic, 0)
	// 统计符合条件的记录总数
	//err = models.GetProblemList(keyword, categoryIdentity).Distinct("`problem_basic`.`id`").Count(&count).Error
	err = models.GetProblemList(keyword).Distinct("`problem_basic`.`id`").Count(&count).Error
	if err != nil {
		// 若统计记录总数出错，记录错误日志
		log.Println("GetProblemList Count Error:", err)
		return
	}
	// 查询符合条件的问题列表
	//err = models.GetProblemList(keyword, categoryIdentity).Offset(page).Limit(size).Find(&list).Error
	err = models.GetProblemList(keyword).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		// 若查询问题列表出错，记录错误日志
		log.Println("Get Problem List Error:", err)
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}
