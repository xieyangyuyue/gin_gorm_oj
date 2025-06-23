package service

import (
	"errors"
	"gin_gorm_oj/define"
	"gin_gorm_oj/models"
	"gin_gorm_oj/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) { // 定义 GetProblemList 函数，它是一个 Gin 框架的 HTTP 请求处理函数。
	// 从请求查询参数中获取 'size'，如果不存在则使用 define.DefaultSize 作为默认值，并转换为 int。
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil { // 检查 size 参数转换是否发生错误。
		log.Printf("GetProblemList: 分页参数size转换错误: %v\n", err) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{                          // 返回JSON格式错误响应
			"code": -1,
			"msg":  "分页参数size错误：" + err.Error(),
		})
		return // 终止函数执行。
	}

	// 从请求查询参数中获取 'page'，如果不存在则使用 define.DefaultPage 作为默认值，并转换为 int。
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil { // 检查 page 参数转换是否发生错误。
		log.Printf("GetProblemList: 分页参数page转换错误: %v\n", err) // 如果发生错误，则打印日志。
		c.JSON(http.StatusOK, gin.H{                          // 返回JSON格式错误响应
			"code": -1,
			"msg":  "分页参数page错误：" + err.Error(),
		})
		return // 终止函数执行。
	}
	page = (page - 1) * size                         // 计算分页查询的偏移量（例如，第一页偏移量为 0）。
	var count int64                                  // 声明一个 int64 类型的变量 count，用于存储问题总数。
	keyword := c.Query("keyword")                    // 从请求查询参数中获取 'keyword'。
	categoryIdentity := c.Query("category_identity") // 从请求查询参数中获取 'category_identity'。

	list := make([]*models.ProblemBasic, 0) // 初始化一个 ProblemBasic 结构体指针的切片，用于存放查询到的问题列表。
	// 调用 models 包的方法获取问题列表的查询构建器，然后使用 Distinct() 确保按 ID 去重，并计算符合条件的问题总数，将结果存储到 count 中。
	err = models.GetProblemList(keyword, categoryIdentity).Distinct("`problem_basic`.`id`").Count(&count).Error
	if err != nil { // 检查统计总数时是否发生数据库错误。
		log.Printf("GetProblemList: 统计问题总数错误: %v\n", err) // 如果发生错误，则打印日志。
		c.JSON(http.StatusOK, gin.H{                      // 返回JSON格式错误响应
			"code": -1,
			"msg":  "获取问题总数失败：" + err.Error(),
		})
		return // 终止函数执行。
	}
	// 再次调用 models 包的方法获取查询构建器，应用分页（偏移量和限制数量），并执行查询将结果填充到 list 中。
	err = models.GetProblemList(keyword, categoryIdentity).Offset(page).Limit(size).Find(&list).Error
	if err != nil { // 检查查询问题列表时是否发生数据库错误。
		log.Printf("GetProblemList: 获取问题列表错误: %v\n", err) // 如果发生错误，则打印日志。
		c.JSON(http.StatusOK, gin.H{                      // 返回JSON格式错误响应
			"code": -1,
			"msg":  "获取问题列表失败：" + err.Error(),
		})
		return // 终止函数执行。
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
// @Summary 问题详情
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) { // 定义 GetProblemDetail 函数，它是一个 Gin 框架的 HTTP 请求处理函数。
	identity := c.Query("identity") // 从请求查询参数中获取 'identity'（问题唯一标识）。
	if identity == "" {             // 检查 identity 是否为空。
		c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
			"code": -1,           // 设置自定义错误码为 -1。
			"msg":  "问题唯一标识不能为空", // 设置错误信息。
		})
		return // 终止函数执行。
	}
	data := new(models.ProblemBasic) // 初始化一个 ProblemBasic 结构体指针，用于存放查询到的问题详情。
	// 使用 GORM 构建查询，查找 identity 字段与给定值匹配的问题。
	// 预加载关联的 ProblemCategories 和 ProblemCategories 下的 CategoryBasic 信息。
	// 执行查询，尝试获取第一条匹配的记录，并将结果填充到 data 中。
	err := models.DB.Where("identity = ?", identity).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").First(&data).Error
	if err != nil { // 检查查询是否发生错误。
		if errors.Is(err, gorm.ErrRecordNotFound) { // 如果错误是 GORM 的记录未找到错误。
			c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
				"code": -1,      // 设置自定义错误码为 -1。
				"msg":  "问题不存在", // 设置错误信息。
			})
			return // 终止函数执行。
		}
		log.Printf("GetProblemDetail: 查询问题详情错误: %v, identity: %s\n", err, identity) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{                                                // 对于其他类型的数据库错误，返回 JSON 响应。
			"code": -1,                        // 设置自定义错误码为 -1。
			"msg":  "获取问题详情失败：" + err.Error(), // 设置包含原始错误信息的错误信息。
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
// @Summary 问题创建
// @Accept json
// @Param authorization header string true "authorization"
// @Param data body define.ProblemBasic true "ProblemBasic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func ProblemCreate(c *gin.Context) { // 定义 ProblemCreate 函数，它是一个 Gin 框架的 HTTP 请求处理函数，用于创建问题。
	in := new(define.ProblemBasic) // 初始化一个 define.ProblemBasic 结构体指针，用于接收请求体中的 JSON 数据。
	err := c.ShouldBindJSON(in)    // 尝试将请求体中的 JSON 数据绑定到 in 结构体。
	if err != nil {                // 检查 JSON 绑定是否发生错误。
		log.Printf("[ProblemCreate JsonBind Error] : %v\n", err) // 如果发生错误，则打印 JSON 绑定错误日志。
		c.JSON(http.StatusOK, gin.H{                             // 返回 JSON 格式的响应。
			"code": -1,       // 设置自定义错误码为 -1。
			"msg":  "参数解析错误", // 设置错误信息。
		})
		return // 终止函数执行。
	}

	// 检查所有必填字段是否为空或零值。
	if in.Title == "" || in.Content == "" || len(in.ProblemCategories) == 0 || len(in.TestCases) == 0 || in.MaxRuntime == 0 || in.MaxMem == 0 {
		c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
			"code": -1,            // 设置自定义错误码为 -1。
			"msg":  "必填参数不能为空或零值", // 设置错误信息。
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

	// 处理分类
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

	// 处理测试用例
	testCaseBasics := make([]*models.TestCase, 0) // 初始化一个 TestCase 结构体指针的切片，用于存放测试用例。
	for _, v := range in.TestCases {              // 遍历输入中提供的测试用例列表。
		// 举个例子 {"input":"1 2\n","output":"3\n"}
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

	// 创建问题
	err = models.DB.Create(data).Error // 使用 GORM 的 Create 方法将 data（包含问题、分类和测试用例）保存到数据库中。
	if err != nil {                    // 检查数据库创建操作是否发生错误。
		log.Printf("ProblemCreate Error: %v\n", err) // 记录详细错误日志
		c.JSON(http.StatusOK, gin.H{                 // 返回 JSON 格式的响应。
			"code": -1,                      // 设置自定义错误码为 -1。
			"msg":  "问题创建失败：" + err.Error(), // 设置包含原始错误信息的错误信息。
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
func ProblemModify(c *gin.Context) { // 定义 ProblemModify 函数，用于处理问题修改的HTTP请求
	in := new(define.ProblemBasic) // 初始化一个 define.ProblemBasic 结构体指针，用于接收请求体中的JSON数据
	err := c.ShouldBindJSON(in)    // 尝试将请求体中的JSON数据绑定到in结构体
	if err != nil {                // 检查JSON绑定是否发生错误
		log.Printf("[ProblemModify JsonBind Error] : %v\n", err) // 记录JSON绑定错误日志
		c.JSON(http.StatusOK, gin.H{                             // 返回JSON格式错误响应
			"code": -1,       // 设置自定义错误码
			"msg":  "参数解析错误", // 设置错误信息
		})
		return // 终止函数执行
	}

	// 检查所有必填字段是否为空或零值
	if in.Identity == "" || in.Title == "" || in.Content == "" || len(in.ProblemCategories) == 0 || len(in.TestCases) == 0 || in.MaxRuntime == 0 || in.MaxMem == 0 {
		c.JSON(http.StatusOK, gin.H{ // 返回JSON格式错误响应
			"code": -1,            // 设置自定义错误码
			"msg":  "必填参数不能为空或零值", // 设置错误信息
		})
		return // 终止函数执行
	}

	// 使用GORM事务确保所有数据库操作的原子性
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 问题基础信息保存 problem_basic
		problemBasic := &models.ProblemBasic{ // 创建ProblemBasic结构体实例用于更新
			Identity:   in.Identity,               // 问题唯一标识
			Title:      in.Title,                  // 问题标题
			Content:    in.Content,                // 问题内容
			MaxRuntime: in.MaxRuntime,             // 最大运行时间
			MaxMem:     in.MaxMem,                 // 最大内存限制
			UpdatedAt:  models.MyTime(time.Now()), // 更新时间
		}
		// 根据identity更新problem_basic表
		err := tx.Where("identity = ?", in.Identity).Updates(problemBasic).Error
		if err != nil { // 检查更新是否出错
			log.Printf("ProblemModify: 更新问题基本信息错误: %v, identity: %s\n", err, in.Identity) // 记录详细错误日志
			return err                                                                    // 返回错误，触发事务回滚
		}

		// 查询问题详情，以便获取其ID用于关联表的更新
		err = tx.Where("identity = ?", in.Identity).Find(problemBasic).Error
		if err != nil { // 检查查询是否出错
			log.Printf("ProblemModify: 查询问题ID错误: %v, identity: %s\n", err, in.Identity) // 记录详细错误日志
			return err                                                                  // 返回错误，触发事务回滚
		}

		// 关联问题分类的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil { // 检查删除是否出错
			log.Printf("ProblemModify: 删除旧分类关联错误: %v, problem_id: %d\n", err, problemBasic.ID) // 记录详细错误日志
			return err                                                                         // 返回错误，触发事务回滚
		}
		// 2、新增新的关联关系
		pcs := make([]*models.ProblemCategory, 0) // 创建ProblemCategory切片
		for _, id := range in.ProblemCategories { // 遍历新的分类ID
			pcs = append(pcs, &models.ProblemCategory{ // 添加新的关联关系
				ProblemId:  problemBasic.ID,           // 问题ID
				CategoryId: uint(id),                  // 分类ID
				CreatedAt:  models.MyTime(time.Now()), // 创建时间
				UpdatedAt:  models.MyTime(time.Now()), // 更新时间
			})
		}
		err = tx.Create(&pcs).Error // 批量创建新的关联关系
		if err != nil {             // 检查创建是否出错
			log.Printf("ProblemModify: 创建新分类关联错误: %v, problem_id: %d\n", err, problemBasic.ID) // 记录详细错误日志
			return err                                                                         // 返回错误，触发事务回滚
		}

		// 关联测试案例的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_identity = ?", in.Identity).Delete(new(models.TestCase)).Error
		if err != nil { // 检查删除是否出错
			log.Printf("ProblemModify: 删除旧测试案例错误: %v, problem_identity: %s\n", err, in.Identity) // 记录详细错误日志
			return err                                                                           // 返回错误，触发事务回滚
		}
		// 2、增加新的关联关系
		tcs := make([]*models.TestCase, 0) // 创建TestCase切片
		for _, v := range in.TestCases {   // 遍历新的测试案例
			// 举个例子 {"input":"1 2\n","output":"3\n"}
			tcs = append(tcs, &models.TestCase{ // 添加新的测试案例
				Identity:        utils.GetUUID(),           // 测试案例唯一标识
				ProblemIdentity: in.Identity,               // 所属问题标识
				Input:           v.Input,                   // 输入数据
				Output:          v.Output,                  // 输出数据
				CreatedAt:       models.MyTime(time.Now()), // 创建时间
				UpdatedAt:       models.MyTime(time.Now()), // 更新时间
			})
		}
		err = tx.Create(tcs).Error // 批量创建新的测试案例
		if err != nil {            // 检查创建是否出错
			log.Printf("ProblemModify: 创建新测试案例错误: %v, problem_identity: %s\n", err, in.Identity) // 记录详细错误日志
			return err                                                                           // 返回错误，触发事务回滚
		}
		return nil // 事务成功，返回nil
	}); err != nil { // 检查事务是否出错
		c.JSON(http.StatusOK, gin.H{ // 返回JSON格式错误响应
			"code": -1,                      // 设置自定义错误码
			"msg":  "问题修改失败：" + err.Error(), // 设置包含原始错误信息的错误信息
		})
		return // 终止函数执行
	}
	c.JSON(http.StatusOK, gin.H{ // 如果问题修改成功，返回JSON响应
		"code": 200,      // 设置响应状态码为200，表示成功
		"msg":  "问题修改成功", // 设置成功信息
	})
}
