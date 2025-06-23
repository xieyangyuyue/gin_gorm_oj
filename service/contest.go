package service

import (
	"errors"                  // errors 包用于处理 Go 语言中的错误类型
	"gin_gorm_oj/middlewares" // 引入自定义的中间件包，可能包含用户认证相关逻辑
	"log"                     // log 包用于打印日志信息
	"net/http"                // net/http 包提供了 HTTP 客户端和服务端实现
	"strconv"                 // strconv 包用于字符串和基本数据类型之间的转换
	"time"                    // time 包用于时间相关的操作

	"gin_gorm_oj/define"       // 引入自定义的 define 包，可能包含常量和结构体定义
	"gin_gorm_oj/models"       // 引入 models 包，包含数据库模型定义和操作
	"gin_gorm_oj/utils"        // 引入自定义的 utils 工具包
	"github.com/gin-gonic/gin" // 引入 Gin Web 框架
	"gorm.io/gorm"             // 引入 GORM ORM 库
)

// ContestCreate
// @Tags 管理员私有方法
// @Summary 竞赛创建
// @Accept json
// @Param authorization header string true "authorization"
// @Param data body define.ContestBasic true "contestBasic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/contest-create [post]
// ContestCreate 函数用于处理管理员创建竞赛的请求
func ContestCreate(c *gin.Context) {
	in := new(define.ContestBasic) // 创建一个 ContestBasic 结构体实例用于接收请求体
	err := c.ShouldBindJSON(in)    // 将请求体绑定到 in 结构体
	if err != nil {
		log.Printf("[JsonBind Error] : %v", err) // 记录 JSON 绑定错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误: JSON 解析失败", // 返回参数错误信息
		})
		return
	}

	// 检查必需参数是否为空
	if in.Name == "" || in.Content == "" || len(in.ProblemBasics) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空: 竞赛名称、内容或题目列表缺失", // 返回参数缺失错误信息
		})
		return
	}

	// 检查开始时间和结束时间是否有效
	if in.StartAt >= in.EndAt {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "竞赛开始时间必须早于结束时间",
		})
		return
	}

	// 生成竞赛的唯一标识
	identity := utils.GetUUID()
	// 创建 ContestBasic 模型实例
	data := &models.ContestBasic{
		Identity:  identity,                                // 设置唯一标识
		Name:      in.Name,                                 // 设置竞赛名称
		Content:   in.Content,                              // 设置竞赛内容
		StartAt:   models.MyTime(utils.ToTime(in.StartAt)), // 转换并设置开始时间
		EndAt:     models.MyTime(utils.ToTime(in.EndAt)),   // 转换并设置结束时间
		CreatedAt: models.MyTime(time.Now()),               // 设置创建时间
		UpdatedAt: models.MyTime(time.Now()),               // 设置更新时间
	}

	// 构建竞赛与问题的关联关系列表
	contestProblems := make([]*models.ContestProblem, 0, len(in.ProblemBasics)) // 预分配切片容量
	for _, id := range in.ProblemBasics {
		// 检查问题ID是否有效（非0）
		if id == 0 {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题ID不能为0",
			})
			return
		}
		contestProblems = append(contestProblems, &models.ContestProblem{
			ContestId: data.ID,                   // 关联竞赛ID（注意：此时 data.ID 可能还未生成，需要依赖 GORM 的 Hooks 或在事务中处理）
			ProblemId: uint(id),                  // 关联问题ID
			CreatedAt: models.MyTime(time.Now()), // 设置创建时间
			UpdatedAt: models.MyTime(time.Now()), // 设置更新时间
		})
	}
	data.ContestProblems = contestProblems // 将关联问题列表赋值给 ContestBasic

	// 使用事务创建竞赛及其关联问题，确保数据一致性
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		// 创建竞赛基础信息
		if err := tx.Create(data).Error; err != nil {
			return errors.New("ContestBasic 创建失败: " + err.Error())
		}

		// 更新 ContestProblems 的 ContestId，因为 data.ID 在上面 Create 之后才生成
		for _, cp := range contestProblems {
			cp.ContestId = data.ID
		}
		// 批量创建竞赛与问题的关联关系
		if len(contestProblems) > 0 { // 只有在有题目时才执行创建操作
			if err := tx.Create(&contestProblems).Error; err != nil {
				return errors.New("ContestProblem 关联创建失败: " + err.Error())
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Contest Create Transaction Error: %v", err) // 记录事务错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "竞赛创建失败：" + err.Error(), // 返回创建失败信息
		})
		return
	}

	// 创建成功，返回竞赛唯一标识
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": data.Identity,
		},
		"msg": "竞赛创建成功", // 添加成功消息
	})
}

// GetContestList
// @Tags 公共方法
// @Summary 竞赛列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /contest-list [get]
// GetContestList 函数用于获取竞赛列表
func GetContestList(c *gin.Context) {
	// 从查询参数中获取每页显示的数量，若未提供则使用 define.DefaultSize
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	// 从查询参数中获取当前页码，若未提供则使用 define.DefaultPage
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Printf("GetContestList Page strconv Error: %v", err) // 记录页码转换错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误：页码非法",
		})
		return
	}
	// 计算数据库查询的偏移量，确保页码从 1 开始计算
	offset := (page - 1) * size
	var count int64               // 用于存储符合条件的记录总数
	keyword := c.Query("keyword") // 从查询参数中获取关键词

	list := make([]*models.ContestBasic, 0) // 初始化竞赛列表切片

	// 使用 models.GetContestList 构建查询，先获取总数
	tx := models.GetContestList(keyword)
	err = tx.Count(&count).Error // 统计总数
	if err != nil {
		log.Printf("GetContestList Count Error: %v", err) // 记录统计错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取竞赛列表总数失败：" + err.Error(),
		})
		return
	}

	// 再次使用 models.GetContestList 构建查询，进行分页查找
	err = tx.Offset(offset).Limit(size).Find(&list).Error // 分页查询数据
	if err != nil {
		log.Printf("Get Contest List Error: %v", err) // 记录查询错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取竞赛列表失败：" + err.Error(),
		})
		return
	}
	// 返回竞赛列表和总数
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
		"msg": "获取竞赛列表成功", // 添加成功消息
	})
}

// GetContestDetail
// @Tags 公共方法
// @Summary 竞赛详情
// @Param identity query string false "contest identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /contest-detail [get]
// GetContestDetail 函数用于获取竞赛详情
func GetContestDetail(c *gin.Context) {
	identity := c.Query("identity") // 从查询参数中获取竞赛唯一标识
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "竞赛唯一标识不能为空", // 返回参数缺失错误信息
		})
		return
	}
	data := new(models.ContestBasic) // 创建 ContestBasic 结构体实例用于接收查询结果
	// 根据唯一标识查询竞赛详情，并预加载关联的问题和用户
	err := models.DB.Where("identity = ?", identity).
		Preload("ContestProblems").Preload("ContestProblems.ProblemBasic"). // 预加载竞赛问题及其对应的问题基本信息
		Preload("ContestUsers").                                            // 预加载竞赛用户
		First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 若记录未找到
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "竞赛不存在", // 返回竞赛不存在信息
			})
			return
		}
		log.Printf("Get Contest Detail Error: %v", err) // 记录数据库查询错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取竞赛详情失败：" + err.Error(), // 返回获取失败信息
		})
		return
	}
	// 返回竞赛详情数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
		"msg":  "获取竞赛详情成功", // 添加成功消息
	})
}

// ContestModify
// @Tags 管理员私有方法
// @Summary 竞赛修改
// @Param authorization header string true "authorization"
// @Param data body define.ContestBasic true "ContestBasic"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/contest-modify [put]
// ContestModify 函数用于处理管理员修改竞赛的请求
func ContestModify(c *gin.Context) {
	in := new(define.ContestBasic) // 创建一个 ContestBasic 结构体实例用于接收请求体
	err := c.ShouldBindJSON(in)    // 将请求体绑定到 in 结构体
	if err != nil {
		log.Printf("[JsonBind Error] : %v", err) // 记录 JSON 绑定错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误: JSON 解析失败", // 返回参数错误信息
		})
		return
	}
	// 检查必需参数是否为空
	if in.Identity == "" || in.Name == "" || in.Content == "" || len(in.ProblemBasics) == 0 || in.StartAt == 0 || in.EndAt == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空: 竞赛唯一标识、名称、内容、题目列表或时间信息缺失", // 返回参数缺失错误信息
		})
		return
	}

	// 检查开始时间和结束时间是否有效
	if in.StartAt >= in.EndAt {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "竞赛开始时间必须早于结束时间",
		})
		return
	}

	// 使用事务进行竞赛信息的修改，确保数据一致性
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 竞赛基础信息保存 contest_basic
		contestBasic := &models.ContestBasic{
			Name:      in.Name,                                 // 更新竞赛名称
			Content:   in.Content,                              // 更新竞赛内容
			StartAt:   models.MyTime(utils.ToTime(in.StartAt)), // 更新开始时间
			EndAt:     models.MyTime(utils.ToTime(in.EndAt)),   // 更新结束时间
			UpdatedAt: models.MyTime(time.Now()),               // 更新更新时间
		}
		// 根据唯一标识更新竞赛基础信息
		err := tx.Where("identity = ?", in.Identity).Updates(contestBasic).Error
		if err != nil {
			return errors.New("竞赛基础信息更新失败: " + err.Error())
		}
		// 查询更新后的竞赛详情，以获取 ID
		err = tx.Where("identity = ?", in.Identity).First(contestBasic).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("竞赛不存在或已被删除")
			}
			return errors.New("查询竞赛详情失败: " + err.Error())
		}

		// 关联问题题目的更新
		// 1、删除已存在的关联关系
		err = tx.Where("contest_id = ?", contestBasic.ID).Delete(new(models.ContestProblem)).Error
		if err != nil {
			return errors.New("删除旧的竞赛问题关联失败: " + err.Error())
		}
		// 2、新增新的关联关系
		cps := make([]*models.ContestProblem, 0, len(in.ProblemBasics)) // 预分配切片容量
		for _, id := range in.ProblemBasics {
			// 检查问题ID是否有效（非0）
			if id == 0 {
				return errors.New("问题ID不能为0")
			}
			cps = append(cps, &models.ContestProblem{
				ContestId: contestBasic.ID,           // 关联竞赛ID
				ProblemId: uint(id),                  // 关联问题ID
				CreatedAt: models.MyTime(time.Now()), // 设置创建时间
				UpdatedAt: models.MyTime(time.Now()), // 设置更新时间
			})
		}
		if len(cps) > 0 { // 只有在有题目时才执行创建操作
			err = tx.Create(&cps).Error
			if err != nil {
				return errors.New("新增竞赛问题关联失败: " + err.Error())
			}
		}
		return nil // 事务成功
	}); err != nil {
		log.Printf("Contest Modify Transaction Error: %v", err) // 记录事务错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "竞赛修改失败：" + err.Error(), // 返回修改失败信息
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "竞赛修改成功", // 返回成功信息
	})
}

// ContestDelete
// @Tags 管理员私有方法
// @Summary 竞赛删除
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/contest-delete [delete]
// ContestDelete 函数用于处理管理员删除竞赛的请求
func ContestDelete(c *gin.Context) {
	identity := c.Query("identity") // 从查询参数中获取竞赛唯一标识
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空: 竞赛唯一标识缺失", // 返回参数缺失错误信息
		})
		return
	}
	// 使用事务进行竞赛的删除操作，包括所有关联数据
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		cbs := &models.ContestBasic{}
		// 查询竞赛是否存在，以便获取其 ID
		err := tx.Where("identity = ?", identity).First(cbs).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("竞赛不存在或已被删除")
			}
			return errors.New("查询竞赛失败: " + err.Error())
		}

		// 删除竞赛与问题的关联
		err = tx.Where("contest_id = ?", cbs.ID).Delete(new(models.ContestProblem)).Error
		if err != nil {
			return errors.New("删除竞赛问题关联失败: " + err.Error())
		}

		// 删除竞赛与用户的关联（报名信息）
		err = tx.Where("contest_id = ?", cbs.ID).Delete(new(models.ContestUser)).Error
		if err != nil {
			return errors.New("删除竞赛用户关联失败: " + err.Error())
		}

		// 删除竞赛基础信息
		err = tx.Where("identity = ?", identity).Delete(cbs).Error
		if err != nil {
			return errors.New("删除竞赛基础信息失败: " + err.Error())
		}
		return nil // 事务成功
	}); err != nil {
		log.Printf("Contest Delete Transaction Error: %v", err) // 记录事务错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "竞赛删除失败：" + err.Error(), // 返回删除失败信息
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功", // 返回成功信息
	})
}

// ContestRegistration
// @Tags 用户私有方法
// @Summary 竞赛报名
// @Param authorization header string true "authorization"
// @Param contest_identity query string true "contest_identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/contest-registration [post]
// ContestRegistration 函数用于处理用户报名竞赛的请求
func ContestRegistration(c *gin.Context) {
	contestIdentity := c.Query("contest_identity") // 从查询参数中获取竞赛唯一标识
	if contestIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空: 竞赛唯一标识缺失", // 返回参数缺失错误信息
		})
		return
	}
	// 查询竞赛是否存在
	cb := &models.ContestBasic{}
	err := models.DB.Where("identity = ?", contestIdentity).First(cb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 竞赛不存在
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "竞赛不存在", // 返回竞赛不存在信息
			})
			return
		}
		log.Printf("Contest Registration Query Contest Error: %v", err) // 记录数据库查询错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常: 查询竞赛失败", // 返回数据库异常信息
		})
		return
	}

	// 判断竞赛是否已开始或已结束 (更严谨的判断)
	now := time.Now()
	if now.Before(time.Time(cb.StartAt)) { // 在开始时间之前不能报名 (如果业务允许预报名则移除此判断)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "竞赛尚未开始报名",
		})
		return
	}
	if now.After(time.Time(cb.EndAt)) { // 在结束时间之后不能报名
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "竞赛已结束，无法报名",
		})
		return
	}

	u, exists := c.Get("user_claims") // 从上下文中获取用户声明
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户认证信息不存在，请重新登录",
		})
		return
	}
	userClaim, ok := u.(*middlewares.UserClaims) // 类型断言
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户认证信息解析失败",
		})
		return
	}

	// 判断用户是否已报名该竞赛
	var contestUser models.ContestUser
	err = models.DB.Where("contest_id = ? AND user_identity = ?", cb.ID, userClaim.Identity).First(&contestUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // 未找到报名记录，可以继续报名
			// 继续执行报名逻辑
		} else { // 数据库异常
			log.Printf("Contest Registration Query User Contest Error: %v", err) // 记录数据库查询错误日志
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "数据库异常: 查询报名信息失败", // 返回数据库异常信息
			})
			return
		}
	} else { // 已找到报名记录，表示已报名
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "您已报名该竞赛", // 返回已报名信息
		})
		return
	}

	// 进行报名：创建新的竞赛用户关联记录
	cu := &models.ContestUser{
		ContestId:    cb.ID,                     // 关联竞赛ID
		UserIdentity: userClaim.Identity,        // 关联用户唯一标识
		CreatedAt:    models.MyTime(time.Now()), // 设置创建时间
		UpdatedAt:    models.MyTime(time.Now()), // 设置更新时间
	}
	err = models.DB.Create(cu).Error // 创建报名记录
	if err != nil {
		log.Printf("Contest Registration Create Error: %v", err) // 记录创建错误日志
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "数据库异常: 创建报名信息失败", // 返回创建失败信息
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "报名成功", // 返回成功信息
	})
}
