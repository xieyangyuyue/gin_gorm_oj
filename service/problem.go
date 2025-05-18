package service

import (
	"gin_gorm_oj/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","msg","","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	models.GetProblemList()
	c.String(http.StatusOK, "Get Problem List")
}
