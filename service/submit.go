package service

import (
	"bytes"
	"errors"
	"gin_gorm_oj/define"
	"gin_gorm_oj/middlewares"
	"gin_gorm_oj/models"
	"gin_gorm_oj/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
// GetSubmitList 函数用于获取提交列表，支持分页、按问题标识、用户标识和状态过滤
func GetSubmitList(c *gin.Context) {
	// 从查询参数中获取每页显示的数量，若未提供则使用默认值
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	// 从查询参数中获取当前页码，若未提供则使用默认值
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		// 若页码转换为整数时出错，记录错误日志
		log.Println("GetProblemList Page strconv Error:", err)
		return
	}
	// 计算偏移量
	page = (page - 1) * size

	// 用于存储符合条件的记录总数
	var count int64
	// 用于存储查询到的提交记录列表
	list := make([]models.SubmitBasic, 0)

	// 从查询参数中获取问题标识
	problemIdentity := c.Query("problem_identity")
	// 从查询参数中获取用户标识
	userIdentity := c.Query("user_identity")
	// 从查询参数中获取提交状态，若未提供则默认为 0
	status, _ := strconv.Atoi(c.Query("status"))
	// 调用 models 包中的函数构建查询条件
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	// 执行查询，获取记录总数，并进行分页查询
	err = tx.Count(&count).Offset(page).Limit(size).Find(&list).Error
	if err != nil {
		// 若查询出错，返回错误信息
		log.Println("Get Problem List Error:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Submit List Error:" + err.Error(),
		})
		return
	}
	// 查询成功，返回提交列表和记录总数
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
	})
}

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit [post]
// Submit 函数用于处理用户的代码提交请求
func Submit(c *gin.Context) {
	// 从查询参数中获取问题标识
	problemIdentity := c.Query("problem_identity")
	// 从请求体中读取用户提交的代码
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// 若读取代码出错，返回错误信息
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Read Code Error:" + err.Error(),
		})
		return
	}
	// 调用 helper 包中的函数将代码保存到文件系统
	path, err := utils.CodeSave(code)
	if err != nil {
		// 若代码保存出错，返回错误信息
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Code Save Error:" + err.Error(),
		})
		return
	}
	// 从上下文中获取用户的声明信息
	u, _ := c.Get("user_claims")
	userClaim := u.(*middlewares.UserClaims)
	// 创建一个新的提交记录对象
	sb := &models.SubmitBasic{
		Identity:        utils.GetUUID(),           // 生成唯一标识
		ProblemIdentity: problemIdentity,           // 关联问题标识
		UserIdentity:    userClaim.Identity,        // 关联用户标识
		Path:            path,                      // 代码保存路径
		CreatedAt:       models.MyTime(time.Now()), // 创建时间
		UpdatedAt:       models.MyTime(time.Now()), // 更新时间
	}
	// 从数据库中查询关联的问题信息，并预加载测试用例
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(pb).Error
	if err != nil {
		// 若查询问题信息出错，返回错误信息
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Problem Error:" + err.Error(),
		})
		return
	}
	// 定义不同状态的 channel，用于并发处理测试结果
	// 答案错误的 channel
	WA := make(chan int)
	// 超内存的 channel
	OOM := make(chan int)
	// 编译错误的 channel
	CE := make(chan int)
	// 答案正确的 channel
	AC := make(chan int)
	// 非法代码的 channel
	EC := make(chan struct{})

	// 记录通过的测试用例个数
	passCount := 0
	// 用于并发安全的锁
	var lock sync.Mutex
	// 提示信息
	var msg string

	// 检查代码的合法性
	v, err := utils.CheckGoCodeValid(path)
	if err != nil {
		// 若代码检查出错，返回错误信息
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Code Check Error:" + err.Error(),
		})
		return
	}
	if !v {
		// 若代码不合法，向 EC channel 发送信号
		go func() {
			EC <- struct{}{}
		}()
	} else {
		// 若代码合法，遍历问题的所有测试用例
		for _, testCase := range pb.TestCases {
			testCase := testCase
			go func() {
				// 创建一个执行 Go 代码的命令
				cmd := exec.Command("go", "run", path)
				// 用于存储命令的标准输出和标准错误输出
				var out, stderr bytes.Buffer
				cmd.Stderr = &stderr
				cmd.Stdout = &out
				// 获取命令的标准输入管道
				stdinPipe, err := cmd.StdinPipe()
				if err != nil {
					// 若获取标准输入管道出错，记录错误日志并退出
					log.Fatalln(err)
				}
				// 向标准输入管道写入测试用例的输入数据
				_, _ = io.WriteString(stdinPipe, testCase.Input+"\n")

				// 记录执行前的内存使用情况
				var bm runtime.MemStats
				runtime.ReadMemStats(&bm)
				// 执行命令
				if err := cmd.Run(); err != nil {
					// 若命令执行出错，记录错误信息
					log.Println(err, stderr.String())
					if err.Error() == "exit status 2" {
						// 若错误码为 2，认为是编译错误，向 CE channel 发送信号
						msg = stderr.String()
						CE <- 1
						return
					}
				}
				// 记录执行后的内存使用情况
				var em runtime.MemStats
				runtime.ReadMemStats(&em)

				// 答案错误
				if testCase.Output != out.String() {
					// 若输出结果与测试用例的期望输出不一致，向 WA channel 发送信号
					log.Println(testCase.Output)
					log.Println(out.String())
					WA <- 1
					return
				}
				// 运行超内存
				if em.Alloc/1024-(bm.Alloc/1024) > uint64(pb.MaxMem) {
					// 若内存使用超过问题的最大内存限制，向 OOM channel 发送信号
					OOM <- 1
					return
				}
				// 加锁，保证并发安全
				lock.Lock()
				// 通过的测试用例个数加 1
				passCount++
				if passCount == len(pb.TestCases) {
					// 若所有测试用例都通过，向 AC channel 发送信号
					AC <- 1
				}
				// 解锁
				lock.Unlock()
			}()
		}
	}

	// 监听不同状态的 channel，根据收到的信号更新提交记录的状态和提示信息
	select {
	case <-EC:
		msg = "无效代码"
		sb.Status = 6
	case <-WA:
		msg = "答案错误"
		sb.Status = 2
	case <-OOM:
		msg = "运行超内存"
		sb.Status = 4
	case <-CE:
		sb.Status = 5
	case <-AC:
		msg = "答案正确"
		sb.Status = 1
	case <-time.After(time.Millisecond * time.Duration(pb.MaxRuntime)):
		// 若在规定时间内未完成所有测试用例
		if passCount == len(pb.TestCases) {
			sb.Status = 1
			msg = "答案正确"
		} else {
			sb.Status = 3
			msg = "运行超时"
		}
	}

	// 开启数据库事务，更新提交记录、用户信息和问题信息
	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		// 保存提交记录
		err = tx.Create(sb).Error
		if err != nil {
			return errors.New("SubmitBasic Save Error:" + err.Error())
		}
		// 定义要更新的字段
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if sb.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		// 更新 user_basic
		err = tx.Model(new(models.UserBasic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			return errors.New("UserBasic Modify Error:" + err.Error())
		}
		// 更新 problem_basic
		err = tx.Model(new(models.ProblemBasic)).Where("identity = ?", problemIdentity).Updates(m).Error
		if err != nil {
			return errors.New("ProblemBasic Modify Error:" + err.Error())
		}
		return nil
	}); err != nil {
		// 若事务执行出错，返回错误信息
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Submit Error:" + err.Error(),
		})
		return
	}

	// 提交成功，返回提交结果
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"status": sb.Status,
			"msg":    msg,
		},
	})
}
