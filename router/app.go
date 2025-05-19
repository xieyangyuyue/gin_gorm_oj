package router

import (
	_ "gin_gorm_oj/docs"
	"gin_gorm_oj/service"
	"gin_gorm_oj/utils"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() {
	// 设置Gin运行模式（从配置中读取）
	// utils.AppMode 可能的值：debug/test/release
	gin.SetMode(utils.AppMode)

	// 创建默认路由引擎（自带Logger和Recovery中间件）
	r := gin.Default()
	//r := gin.New()
	// 设置信任网络 []string
	// nil 为不计算，避免性能消耗，上线应当设置
	_ = r.SetTrustedProxies(nil)

	// Swagger 配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 路由规则

	//// 公有方法
	//// 问题
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)
	//// 用户
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)
	//// 排行榜
	//r.GET("/rank-list", service.GetRankList)
	//// 提交记录
	r.GET("/submit-list", service.GetSubmitList)
	//// 分类列表
	//r.GET("/category-list", service.GetCategoryList)
	//// 竞赛列表
	//r.GET("/contest-list", service.GetContestList)
	//r.GET("/contest-detail", service.GetContestDetail)
	//
	//// 管理员私有方法
	//authAdmin := r.Group("/admin", middlewares.AuthAdminCheck())
	////authAdmin := r.Group("/admin")
	//// 问题创建
	//authAdmin.POST("/problem-create", service.ProblemCreate)
	//// 问题修改
	//authAdmin.PUT("/problem-modify", service.ProblemModify)
	//// 分类创建
	//authAdmin.POST("/category-create", service.CategoryCreate)
	//// 分类修改
	//authAdmin.PUT("/category-modify", service.CategoryModify)
	//// 分类删除
	//authAdmin.DELETE("/category-delete", service.CategoryDelete)
	//// 获取测试案例
	//authAdmin.GET("/test-case", service.GetTestCase)
	//
	//// 竞赛创建
	//authAdmin.POST("/contest-create", service.ContestCreate)
	//authAdmin.PUT("/contest-modify", service.ContestModify)
	//authAdmin.DELETE("/contest-delete", service.ContestDelete)
	//
	//// 用户私有方法
	//authUser := r.Group("/user", middlewares.AuthUserCheck())
	//// 代码提交
	//authUser.POST("/submit", service.Submit)
	//authUser.POST("/contest-registration", service.ContestRegistration)
	err := r.Run(utils.HttpPort)
	if err != nil {
		return
	}
}
