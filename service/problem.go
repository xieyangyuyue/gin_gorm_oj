package service // Package service 定义包名为 service，表明该文件属于服务层。

import ( // 导入所需的 Go 语言标准库和第三方库。
	"errors"                   // 导入 errors 包，用于错误处理，例如检查特定的错误类型。
	"gin_gorm_oj/define"       // 导入自定义的 define 包，可能包含常量和公共定义，如默认分页大小。
	"gin_gorm_oj/models"       // 导入 models 包，该包定义了数据库模型（如 ProblemBasic）和数据库操作方法。
	"gin_gorm_oj/utils"        // 导入 utils 包，可能包含工具函数，如生成 UUID。
	"github.com/gin-gonic/gin" // 导入 Gin Web 框架，用于构建 HTTP 服务和处理路由。
	"gorm.io/gorm"             // 导入 GORM 库，这是一个 Go 语言的 ORM 框架，用于数据库操作。
	"log"                      // 导入 log 包，用于输出日志信息。
	"net/http"                 // 导入 net/http 包，提供了 HTTP 客户端和服务端的实现，用于处理 HTTP 状态码等。
	"strconv"                  // 导入 strconv 包，用于字符串和基本数据类型之间的转换。
	"time"                     // 导入 time 包，用于处理时间相关的操作，如获取当前时间。
)

// GetProblemList
// @Tags 公共方法
// Swagger 注解：将此 API 归类到“公共方法”标签下，用于接口文档生成。
// @Summary 问题列表
// Swagger 注解：提供 API 的简要概述，表示获取问题列表。
// @Param page query int false "page"
// Swagger 注解：定义一个名为 'page' 的查询参数，类型为 int，非必填，描述为 "page"。
// @Param size query int false "size"
// Swagger 注解：定义一个名为 'size' 的查询参数，类型为 int，非必填，描述为 "size"。
// @Param keyword query string false "keyword"
// Swagger 注解：定义一个名为 'keyword' 的查询参数，类型为 string，非必填，描述为 "keyword"。
// @Param category_identity query string false "category_identity"
// Swagger 注解：定义一个名为 'category_identity' 的查询参数，类型为 string，非必填，描述为 "category_identity"。
// @Success 200 {string} json "{"code":"200","data":""}"
// Swagger 注解：定义成功响应（HTTP 状态码 200）的 JSON 示例。
// @Router /problem-list [get]
// Swagger 注解：指定此 API 对应的 HTTP 方法（GET）和路由路径。
func GetProblemList(c *gin.Context) { // 定义 GetProblemList 函数，它是一个 Gin 框架的 HTTP 请求处理函数。
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))   // 从请求查询参数中获取 'size'，如果不存在则使用 define.DefaultSize 作为默认值，并转换为 int。忽略错误。
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage)) // 从请求查询参数中获取 'page'，如果不存在则使用 define.DefaultPage 作为默认值，并转换为 int。
	if err != nil {                                                       // 检查 page 参数转换是否发生错误。
		log.Println("GetProblemList Page strconv Error:", err) // 如果发生错误，则打印日志。
		return                                                 // 终止函数执行。
	}
	page = (page - 1) * size                         // 计算分页查询的偏移量（例如，第一页偏移量为 0）。
	var count int64                                  // 声明一个 int64 类型的变量 count，用于存储问题总数。
	keyword := c.Query("keyword")                    // 从请求查询参数中获取 'keyword'。
	categoryIdentity := c.Query("category_identity") // 从请求查询参数中获取 'category_identity'。

	list := make([]*models.ProblemBasic, 0)                                                                     // 初始化一个 ProblemBasic 结构体指针的切片，用于存放查询到的问题列表。
	err = models.GetProblemList(keyword, categoryIdentity).Distinct("`problem_basic`.`id`").Count(&count).Error // 调用 models 包的方法获取问题列表的查询构建器，然后使用 Distinct() 确保按 ID 去重，并计算符合条件的问题总数，将结果存储到 count 中。
	if err != nil {                                                                                             // 检查统计总数时是否发生数据库错误。
		log.Println("GetProblemList Count Error:", err) // 如果发生错误，则打印日志。
		return                                          // 终止函数执行。
	}
	err = models.GetProblemList(keyword, categoryIdentity).Offset(page).Limit(size).Find(&list).Error // 再次调用 models 包的方法获取查询构建器，应用分页（偏移量和限制数量），并执行查询将结果填充到 list 中。
	if err != nil {                                                                                   // 检查查询问题列表时是否发生数据库错误。
		log.Println("Get Problem List Error:", err) // 如果发生错误，则打印日志。
		return                                      // 终止函数执行。
	}
	c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
		"code": 200, // 设置响应状态码为 200，表示成功。
		"data": map[string]interface{}{ // 返回一个包含问题列表和总数的 map。
			"list":  list,  // 问题列表数据。
			"count": count, // 问题总数。
		},
	})
}

// GetProblemDetail
// @Tags 公共方法
// Swagger 注解：将此 API 归类到“公共方法”标签下。
// @Summary 问题详情
// Swagger 注解：提供 API 的简要概述，表示获取问题详情。
// @Param identity query string false "problem identity"
// Swagger 注解：定义一个名为 'identity' 的查询参数，类型为 string，非必填，描述为 "problem identity"。
// @Success 200 {string} json "{"code":"200","data":""}"
// Swagger 注解：定义成功响应（HTTP 状态码 200）的 JSON 示例。
// @Router /problem-detail [get]
// Swagger 注解：指定此 API 对应的 HTTP 方法（GET）和路由路径。
func GetProblemDetail(c *gin.Context) { // 定义 GetProblemDetail 函数，它是一个 Gin 框架的 HTTP 请求处理函数。
	identity := c.Query("identity") // 从请求查询参数中获取 'identity'（问题唯一标识）。
	if identity == "" {             // 检查 identity 是否为空。
		c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
			"code": -1,           // 设置自定义错误码为 -1。
			"msg":  "问题唯一标识不能为空", // 设置错误信息。
		})
		return // 终止函数执行。
	}
	data := new(models.ProblemBasic)                  // 初始化一个 ProblemBasic 结构体指针，用于存放查询到的问题详情。
	err := models.DB.Where("identity = ?", identity). // 使用 GORM 构建查询，查找 identity 字段与给定值匹配的问题。
								Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic"). // 预加载关联的 ProblemCategories 和 ProblemCategories 下的 CategoryBasic 信息。
								First(&data).Error                                                       // 执行查询，尝试获取第一条匹配的记录，并将结果填充到 data 中。
	if err != nil { // 检查查询是否发生错误。
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是 GORM 的记录未找到错误。
			c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
				"code": -1,      // 设置自定义错误码为 -1。
				"msg":  "问题不存在", // 设置错误信息。
			})
			return // 终止函数执行。
		}
		c.JSON(http.StatusOK, gin.H{ // 对于其他类型的数据库错误，返回 JSON 响应。
			"code": -1,                                       // 设置自定义错误码为 -1。
			"msg":  "Get ProblemDetail Error:" + err.Error(), // 设置包含原始错误信息的错误信息。
		})
		return // 终止函数执行。
	}
	c.JSON(http.StatusOK, gin.H{ // 如果问题详情查询成功，返回 JSON 响应。
		"code": 200,  // 设置响应状态码为 200，表示成功。
		"data": data, // 返回问题详情数据。
	})
}

// ProblemCreate
// @Tags 管理员私有方法
// Swagger 注解：将此 API 归类到“管理员私有方法”标签下。
// @Summary 问题创建
// Swagger 注解：提供 API 的简要概述，表示创建问题。
// @Accept json
// Swagger 注解：指定此 API 接受 JSON 格式的请求体。
// @Param authorization header string true "authorization"
// Swagger 注解：定义一个名为 'authorization' 的请求头参数，类型为 string，必填，描述为 "authorization"。
// @Param data body define.ProblemBasic true "ProblemBasic"
// Swagger 注解：定义请求体参数，类型为 define.ProblemBasic，必填，描述为 "ProblemBasic"。
// @Success 200 {string} json "{"code":"200","data":""}"
// Swagger 注解：定义成功响应（HTTP 状态码 200）的 JSON 示例。
// @Router /admin/problem-create [post]
// Swagger 注解：指定此 API 对应的 HTTP 方法（POST）和路由路径。
func ProblemCreate(c *gin.Context) { // 定义 ProblemCreate 函数，它是一个 Gin 框架的 HTTP 请求处理函数，用于创建问题。
	in := new(define.ProblemBasic) // 初始化一个 define.ProblemBasic 结构体指针，用于接收请求体中的 JSON 数据。
	err := c.ShouldBindJSON(in)    // 尝试将请求体中的 JSON 数据绑定到 in 结构体。
	if err != nil {                // 检查 JSON 绑定是否发生错误。
		log.Println("[JsonBind Error] : ", err) // 如果发生错误，则打印 JSON 绑定错误日志。
		c.JSON(http.StatusOK, gin.H{            // 返回 JSON 格式的响应。
			"code": -1,     // 设置自定义错误码为 -1。
			"msg":  "参数错误", // 设置错误信息。
		})
		return // 终止函数执行。
	}

	if in.Title == "" || in.Content == "" || len(in.ProblemCategories) == 0 || len(in.TestCases) == 0 || in.MaxRuntime == 0 || in.MaxMem == 0 { // 检查所有必填字段是否为空或零值。
		c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
			"code": -1,       // 设置自定义错误码为 -1。
			"msg":  "参数不能为空", // 设置错误信息。
		})
		return // 终止函数执行。
	}
	identity := utils.GetUUID()   // 调用 utils 包的 GetUUID 函数，生成一个唯一的 UUID 作为问题标识。
	data := &models.ProblemBasic{ // 创建一个新的 models.ProblemBasic 实例。
		Identity:   identity,                  // 设置问题的唯一标识。
		Title:      in.Title,                  // 设置问题标题。
		Content:    in.Content,                // 设置问题内容。
		MaxRuntime: in.MaxRuntime,             // 设置最大运行时间。
		MaxMem:     in.MaxMem,                 // 设置最大内存限制。
		CreatedAt:  models.MyTime(time.Now()), // 设置创建时间为当前时间。
		UpdatedAt:  models.MyTime(time.Now()), // 设置更新时间为当前时间。
	}
	// 处理分类 // 注释：表示接下来是处理问题分类的代码。
	categoryBasics := make([]*models.ProblemCategory, 0) // 初始化一个 ProblemCategory 结构体指针的切片，用于存放问题分类。
	for _, id := range in.ProblemCategories {            // 遍历输入中提供的问题分类 ID 列表。
		categoryBasics = append(categoryBasics, &models.ProblemCategory{ // 将新的 ProblemCategory 实例添加到切片中。
			ProblemId:  data.ID,                   // 设置问题 ID（在问题创建后 GORM 会自动填充）。
			CategoryId: uint(id),                  // 设置分类 ID。
			CreatedAt:  models.MyTime(time.Now()), // 设置创建时间。
			UpdatedAt:  models.MyTime(time.Now()), // 设置更新时间。
		})
	}
	data.ProblemCategories = categoryBasics // 将处理好的分类切片赋值给 data 结构体的 ProblemCategories 字段。
	// 处理测试用例 // 注释：表示接下来是处理问题测试用例的代码。
	testCaseBasics := make([]*models.TestCase, 0) // 初始化一个 TestCase 结构体指针的切片，用于存放测试用例。
	for _, v := range in.TestCases {              // 遍历输入中提供的测试用例列表。
		// 举个例子 {"input":"1 2\n","output":"3\n"} // 示例测试用例的格式。
		testCaseBasic := &models.TestCase{ // 创建一个新的 TestCase 实例。
			Identity:        utils.GetUUID(),           // 为测试用例生成一个唯一的标识。
			ProblemIdentity: identity,                  // 设置测试用例所属的问题标识。
			Input:           v.Input,                   // 设置测试用例的输入。
			Output:          v.Output,                  // 设置测试用例的输出。
			CreatedAt:       models.MyTime(time.Now()), // 设置创建时间。
			UpdatedAt:       models.MyTime(time.Now()), // 设置更新时间。
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic) // 将新的测试用例添加到切片中。
	}
	data.TestCases = testCaseBasics // 将处理好的测试用例切片赋值给 data 结构体的 TestCases 字段。

	// 创建问题 // 注释：表示接下来是创建问题的数据库操作。
	err = models.DB.Create(data).Error // 使用 GORM 的 Create 方法将 data（包含问题、分类和测试用例）保存到数据库中。
	if err != nil {                    // 检查数据库创建操作是否发生错误。
		c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
			"code": -1,                                    // 设置自定义错误码为 -1。
			"msg":  "Problem Create Error:" + err.Error(), // 设置包含原始错误信息的错误信息。
		})
		return // 终止函数执行。
	}
	c.JSON(http.StatusOK, gin.H{ // 如果问题创建成功，返回 JSON 响应。
		"code": 200, // 设置响应状态码为 200，表示成功。
		"data": map[string]interface{}{ // 返回一个包含新创建问题标识的 map。
			"identity": data.Identity, // 返回新创建问题的唯一标识。
		},
	})
}

// ProblemModify
// @Tags 管理员私有方法
// @Summary 问题修改
// @Param authorization header string true "authorization"
// @Param data body define.ProblemBasic true "ProblemBasic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-modify [put]
func ProblemModify(c *gin.Context) {
	in := new(define.ProblemBasic)
	err := c.ShouldBindJSON(in)
	if err != nil {
		log.Println("[JsonBind Error] : ", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		return
	}
	if in.Identity == "" || in.Title == "" || in.Content == "" || len(in.ProblemCategories) == 0 || len(in.TestCases) == 0 || in.MaxRuntime == 0 || in.MaxMem == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}

	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 问题基础信息保存 problem_basic
		problemBasic := &models.ProblemBasic{
			Identity:   in.Identity,
			Title:      in.Title,
			Content:    in.Content,
			MaxRuntime: in.MaxRuntime,
			MaxMem:     in.MaxMem,
			UpdatedAt:  models.MyTime(time.Now()),
		}
		err := tx.Where("identity = ?", in.Identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		// 查询问题详情
		err = tx.Where("identity = ?", in.Identity).Find(problemBasic).Error
		if err != nil {
			return err
		}

		// 关联问题分类的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil {
			return err
		}
		// 2、新增新的关联关系
		pcs := make([]*models.ProblemCategory, 0)
		for _, id := range in.ProblemCategories {
			pcs = append(pcs, &models.ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(id),
				CreatedAt:  models.MyTime(time.Now()),
				UpdatedAt:  models.MyTime(time.Now()),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil {
			return err
		}
		// 关联测试案例的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_identity = ?", in.Identity).Delete(new(models.TestCase)).Error
		if err != nil {
			return err
		}
		// 2、增加新的关联关系
		tcs := make([]*models.TestCase, 0)
		for _, v := range in.TestCases {
			// 举个例子 {"input":"1 2\n","output":"3\n"}
			tcs = append(tcs, &models.TestCase{
				Identity:        utils.GetUUID(),
				ProblemIdentity: in.Identity,
				Input:           v.Input,
				Output:          v.Output,
				CreatedAt:       models.MyTime(time.Now()),
				UpdatedAt:       models.MyTime(time.Now()),
			})
		}
		err = tx.Create(tcs).Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Problem Modify Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "问题修改成功",
	})
}
