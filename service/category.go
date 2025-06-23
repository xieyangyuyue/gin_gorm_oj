package service

import (
	"gin_gorm_oj/define"
	"gin_gorm_oj/models"
	"gin_gorm_oj/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetCategoryList
// @Tags 公共方法
// @Summary 分类列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-list [get]
func GetCategoryList(c *gin.Context) {
	// 获取URL查询参数中的size，如果未提供则使用define.DefaultSize
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Printf("GetCategoryList: 分页参数size转换错误: %v\n", err) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分页参数size错误：" + err.Error(),
		})
		return
	}

	// 获取URL查询参数中的page，如果未提供则使用define.DefaultPage
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Printf("GetCategoryList: 分页参数page转换错误: %v\n", err) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分页参数page错误：" + err.Error(),
		})
		return
	}

	// 计算数据库查询的偏移量
	page = (page - 1) * size
	// 用于存储分类总数
	var count int64
	// 获取URL查询参数中的keyword
	keyword := c.Query("keyword")

	// 创建CategoryBasic切片用于存储查询结果
	categoryList := make([]*models.CategoryBasic, 0)
	// 构建数据库查询：模糊匹配名称，计算总数，然后应用分页
	err = models.DB.Model(new(models.CategoryBasic)).Where("name like ?", "%"+keyword+"%").
		Count(&count).Limit(size).Offset(page).Find(&categoryList).Error
	// 检查数据库查询是否出错
	if err != nil {
		log.Printf("GetCategoryList Error: %v\n", err) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类列表失败",
		})
		return
	}

	// 返回查询到的分类列表和总数
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  categoryList, // 分类列表
			"count": count,        // 总分类数
		},
	})
}

// CategoryCreate
// @Tags 管理员私有方法
// @Summary 分类创建
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-create [post]
func CategoryCreate(c *gin.Context) {
	// 获取POST表单中的name
	name := c.PostForm("name")
	// 获取POST表单中的parentId，并转换为int类型
	parentId, err := strconv.Atoi(c.PostForm("parentId"))
	if err != nil {
		log.Printf("CategoryCreate: parentId转换错误: %v\n", err) // 记录详细错误日志
		// 即使parentId为可选，如果提供但格式错误，也应返回错误
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "parentId参数格式不正确",
		})
		return
	}

	// 检查name是否为空
	if name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "分类名称不能为空",
		})
		return
	}

	// 创建CategoryBasic结构体实例，并填充分类数据
	category := &models.CategoryBasic{
		Identity:  utils.GetUUID(),           // 生成唯一标识
		Name:      name,                      // 分类名称
		ParentId:  parentId,                  // 父分类ID
		CreatedAt: models.MyTime(time.Now()), // 设置创建时间
		UpdatedAt: models.MyTime(time.Now()), // 设置更新时间
	}
	// 将分类数据插入数据库
	err = models.DB.Create(category).Error
	// 检查数据库插入是否出错
	if err != nil {
		log.Printf("CategoryCreate Error: %v, name: %s\n", err, name) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建分类失败",
		})
		return
	}

	// 返回创建成功信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// CategoryModify
// @Tags 管理员私有方法
// @Summary 分类修改
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param name formData string true "name"
// @Param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-modify [put]
func CategoryModify(c *gin.Context) {
	// 获取POST表单中的identity
	identity := c.PostForm("identity")
	// 获取POST表单中的name
	name := c.PostForm("name")
	// 获取POST表单中的parentId，并转换为int类型
	parentId, err := strconv.Atoi(c.PostForm("parentId"))
	if err != nil {
		log.Printf("CategoryModify: parentId转换错误: %v\n", err) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "parentId参数格式不正确",
		})
		return
	}

	// 检查identity或name是否为空
	if name == "" || identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确，identity和name不能为空",
		})
		return
	}

	// 创建CategoryBasic结构体实例，用于更新操作
	category := &models.CategoryBasic{
		Identity:  identity,                  // 根据identity进行更新
		Name:      name,                      // 更新名称
		ParentId:  parentId,                  // 更新父分类ID
		UpdatedAt: models.MyTime(time.Now()), // 更新时间
	}
	// 在数据库中根据identity更新分类信息
	err = models.DB.Model(new(models.CategoryBasic)).Where("identity = ?", identity).Updates(category).Error
	// 检查数据库更新是否出错
	if err != nil {
		log.Printf("CategoryModify Error: %v, identity: %s, name: %s\n", err, identity, name) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改分类失败",
		})
		return
	}

	// 返回修改成功信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

// CategoryDelete
// @Tags 管理员私有方法
// @Summary 分类删除
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/category-delete [delete]
func CategoryDelete(c *gin.Context) {
	// 获取URL查询参数中的identity
	identity := c.Query("identity")
	// 检查identity是否为空
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确，identity不能为空",
		})
		return
	}

	// 检查该分类下是否已存在问题
	var cnt int64
	// 通过子查询，查询关联到该分类的问题数量
	err := models.DB.Model(new(models.ProblemCategory)).Where("category_id = (SELECT id FROM category_basic WHERE identity = ? LIMIT 1)", identity).Count(&cnt).Error
	// 检查查询关联问题是否出错
	if err != nil {
		log.Printf("CategoryDelete: Get ProblemCategory Error: %v, identity: %s\n", err, identity) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取分类关联的问题失败",
		})
		return
	}
	// 如果该分类下存在问题，则不允许删除
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该分类下面已存在问题，不可删除",
		})
		return
	}

	// 根据identity删除CategoryBasic表中的分类记录
	err = models.DB.Where("identity = ?", identity).Delete(new(models.CategoryBasic)).Error
	// 检查数据库删除是否出错
	if err != nil {
		log.Printf("CategoryDelete: Delete CategoryBasic Error: %v, identity: %s\n", err, identity) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除失败",
		})
		return
	}

	// 返回删除成功信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
