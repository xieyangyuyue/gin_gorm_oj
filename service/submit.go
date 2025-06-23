package service

import (
	"bytes"
	"context" // 引入 context 用于控制并发 goroutine 的超时和取消
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
	"os" // 引入 os 用于文件操作，例如删除临时文件
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync" // 引入 sync 包，用于 WaitGroup 和 Mutex
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
	// 从查询参数中获取每页显示的数量，若未提供则使用默认值。
	// Atoi 不会返回错误，因为 DefaultQuery 确保了字符串总是有效的数字。
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	// 从查询参数中获取当前页码，若未提供则使用默认值。
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		// 若页码转换为整数时出错，记录错误日志并返回。
		log.Printf("GetProblemList Page strconv Error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数错误：页码非法",
		})
		return
	}
	// 计算数据库查询的偏移量。
	offset := (page - 1) * size

	// 用于存储符合条件的记录总数。
	var count int64
	// 用于存储查询到的提交记录列表，初始化为空切片。
	list := make([]models.SubmitBasic, 0)

	// 从查询参数中获取问题标识。
	problemIdentity := c.Query("problem_identity")
	// 从查询参数中获取用户标识。
	userIdentity := c.Query("user_identity")
	// 从查询参数中获取提交状态，若未提供则默认为 0。
	// Atoi 不会返回错误，因为 DefaultQuery 确保了字符串总是有效的数字。
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0")) // 默认状态为 0

	// 调用 models 包中的函数构建查询条件。
	// GetSubmitList 函数应该返回一个 *gorm.DB 实例，以便后续链式调用。
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	// 执行查询，首先获取记录总数，然后进行分页查询。
	err = tx.Count(&count).Offset(offset).Limit(size).Find(&list).Error
	if err != nil {
		// 若查询出错，记录错误日志并返回错误信息。
		log.Printf("Get Submit List Error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取提交列表失败：" + err.Error(),
		})
		return
	}
	// 查询成功，返回提交列表和记录总数。
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,
			"count": count,
		},
		"msg": "获取提交列表成功",
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
	// 从查询参数中获取问题标识，如果不存在则返回错误。
	problemIdentity := c.Query("problem_identity")
	if problemIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题标识不能为空",
		})
		return
	}

	// 从请求体中读取用户提交的代码。
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		// 若读取代码出错，返回错误信息。
		log.Printf("Read Code Error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "读取代码失败：" + err.Error(),
		})
		return
	}
	if len(code) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "提交代码不能为空",
		})
		return
	}

	// 调用 utils 包中的函数将代码保存到文件系统。
	path, err := utils.CodeSave(code)
	if err != nil {
		// 若代码保存出错，返回错误信息。
		log.Printf("Code Save Error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "代码保存失败：" + err.Error(),
		})
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in defer for resource cleanup: %v", r)
		}
		// 先删除文件（如果存在）
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to delete code file %s: %v", path, err)
		}
		// 获取文件所在目录
		dir := filepath.Dir(path)
		// 删除整个目录（包括所有子文件和文件夹）
		if err := os.RemoveAll(dir); err != nil {
			log.Printf("Failed to delete directory %s: %v", dir, err)
		}
	}()
	// 从上下文中获取用户的声明信息。
	u, exists := c.Get("user_claims")
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户认证信息不存在，请重新登录",
		})
		return
	}
	userClaim, ok := u.(*middlewares.UserClaims)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户认证信息解析失败",
		})
		return
	}

	// 创建一个新的提交记录对象。
	sb := &models.SubmitBasic{
		Identity:        utils.GetUUID(),           // 生成唯一标识。
		ProblemIdentity: problemIdentity,           // 关联问题标识。
		UserIdentity:    userClaim.Identity,        // 关联用户标识。
		Path:            path,                      // 代码保存路径。
		CreatedAt:       models.MyTime(time.Now()), // 创建时间。
		UpdatedAt:       models.MyTime(time.Now()), // 更新时间。
	}

	// 从数据库中查询关联的问题信息，并预加载测试用例。
	pb := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(pb).Error
	if err != nil {
		// 若查询问题信息出错，返回错误信息。
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在",
			})
		} else {
			log.Printf("Get Problem Error: %v", err)
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "获取问题信息失败：" + err.Error(),
			})
		}
		return
	}

	// 定义不同状态的 channel，用于并发处理测试结果。
	// 使用带缓冲的 channel，避免 goroutine 写入阻塞。
	WA := make(chan struct{}, 1)  // Wrong Answer 答案错误
	OOM := make(chan struct{}, 1) // Out Of Memory 超内存
	CE := make(chan struct{}, 1)  // Compile Error 编译错误
	//AC := make(chan struct{}, 1)  // Accepted 答案正确
	TLE := make(chan struct{}, 1) // Time Limit Exceeded 运行超时
	EC := make(chan struct{}, 1)  // Error Code/Invalid Code 非法代码

	// 记录通过的测试用例个数。
	passCount := 0
	// 用于并发安全的锁，保护 passCount。
	var lock sync.Mutex
	// 提示信息。
	var msg string
	// 存储最终的提交状态。
	var submitStatus = define.SubmitStatusPending // 默认状态为待判

	// 使用 WaitGroup 等待所有判题 goroutine 完成。
	var wg sync.WaitGroup
	// 用于限制并发判题的 goroutine 数量，例如最多同时运行 4 个判题进程。
	concurrencyLimit := make(chan struct{}, runtime.NumCPU()) // 根据 CPU 核心数限制并发

	// 检查代码的合法性。
	// 这是一个前置检查，如果代码本身非法，则无需执行测试用例。
	valid, err := utils.CheckGoCodeValid(path)
	if err != nil {
		// 若代码检查出错，返回错误信息。
		log.Printf("Code Check Error: %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "代码合法性检查失败：" + err.Error(),
		})
		return
	}
	if !valid {
		// 若代码不合法，向 EC channel 发送信号并直接设置状态。
		// 这里可以直接赋值，因为后续不会再进行判题。
		submitStatus = define.SubmitStatusInvalidCode
		msg = "无效代码：包含非法操作或关键字"
	} else {
		// 若代码合法，遍历问题的所有测试用例进行判题。
		if len(pb.TestCases) == 0 {
			// 如果没有测试用例，默认视为正确（或者根据实际业务逻辑处理）。
			submitStatus = define.SubmitStatusAccepted
			msg = "无测试用例，默认正确"
		} else {
			// 为每个测试用例启动一个 goroutine。
			for _, testCase := range pb.TestCases {
				// 避免闭包问题，将 testCase 拷贝一份。
				tc := testCase
				wg.Add(1)                      // 增加 WaitGroup 计数器
				concurrencyLimit <- struct{}{} // 获取一个并发槽位
				go func() {
					defer wg.Done()                       // goroutine 完成时，减少 WaitGroup 计数器
					defer func() { <-concurrencyLimit }() // 释放并发槽位

					// 创建一个上下文，用于控制命令的超时。
					ctx, cancel := context.WithTimeout(context.Background(), time.Duration(pb.MaxRuntime)*time.Millisecond)
					defer cancel() // 确保在 goroutine 退出时取消上下文

					// 创建一个执行 Go 代码的命令。
					cmd := exec.CommandContext(ctx, "go", "run", path) // 使用 CommandContext 控制超时
					// 用于存储命令的标准输出和标准错误输出。
					var out, stderr bytes.Buffer
					cmd.Stderr = &stderr
					cmd.Stdout = &out
					// 获取命令的标准输入管道。
					stdinPipe, pipeErr := cmd.StdinPipe()
					if pipeErr != nil {
						// 若获取标准输入管道出错，记录错误日志并发送编译错误信号。
						log.Printf("Failed to get stdin pipe: %v", pipeErr)
						select {
						case CE <- struct{}{}:
						default:
						}
						return
					}

					// 向标准输入管道写入测试用例的输入数据。
					// 写入操作应在 cmd.Start() 之后，但关闭操作应在 Wait() 之前。
					// 为了避免死锁，通常在独立的 goroutine 中写入并关闭 stdin。
					go func() {
						defer stdinPipe.Close() // 确保关闭 stdinPipe
						if _, err := io.WriteString(stdinPipe, tc.Input+"\n"); err != nil {
							log.Printf("Failed to write to stdin: %v", err)
						}
					}()

					// 记录执行前的内存使用情况。
					// 注意：Go runtime 的 MemStats 统计的是 Go 进程自身的内存使用，
					// 而不是整个进程组（包括子进程）的内存。对于判题系统，
					// 最好使用 cgroup 或 os 级别的工具来限制和监控内存。
					var bm runtime.MemStats
					runtime.ReadMemStats(&bm)

					// 执行命令。
					cmdErr := cmd.Run() // 阻塞直到命令完成或超时

					// 记录执行后的内存使用情况。
					var em runtime.MemStats
					runtime.ReadMemStats(&em)

					// 检查命令执行结果。
					if cmdErr != nil {
						// 检查是否是上下文超时引起的错误。
						if errors.Is(cmdErr, context.DeadlineExceeded) {
							// 运行超时。
							log.Printf("Run Time Out for test case: %s", tc.Input)
							select {
							case TLE <- struct{}{}:
							default:
							}
							return
						}

						// 检查是否是编译错误（Go 程序的 exit status 2 通常表示编译错误或运行时恐慌）。
						// 这是一个近似判断，更精确的编译错误判断应该在执行前进行编译。
						var exitErr *exec.ExitError
						if errors.As(cmdErr, &exitErr) {
							if exitErr.ExitCode() == 2 { // Go 编译失败通常是 exit status 2
								log.Printf("Compile Error: %s, Stderr: %s", cmdErr, stderr.String())
								msg = stderr.String() // 将编译错误信息作为提示
								select {
								case CE <- struct{}{}:
								default:
								}
								return
							}
						}
						// 其他运行时错误。
						log.Printf("Command Run Error: %v, Stderr: %s", cmdErr, stderr.String())
						// 视为运行时错误，可以归类为 WA 或其他自定义错误状态。
						select {
						case WA <- struct{}{}: // 运行时错误可能导致输出不正确，暂时归类为WA
						default:
						}
						return
					}

					// 答案错误。
					// 注意：out.String() 可能包含额外的换行符，需要进行修剪。
					actualOutput := out.String()
					expectedOutput := tc.Output
					// 移除末尾的换行符和空格进行比较，以应对不同系统或语言习惯。
					//actualOutput = utils.TrimSpaceAndNewlines(actualOutput)
					//expectedOutput = utils.TrimSpaceAndNewlines(expectedOutput)

					if expectedOutput != actualOutput {
						log.Printf("Wrong Answer for test case %s. Expected: '%s', Actual: '%s'", tc.Input, expectedOutput, actualOutput)
						select {
						case WA <- struct{}{}:
						default:
						}
						return
					}

					// 运行超内存。
					// 这里的内存检查是基于 Go runtime 的统计，对于实际判题可能不够精确。
					// 更好的方式是使用 Linux cgroup 等系统级工具进行内存限制和监控。
					memoryUsed := em.Alloc/1024 - bm.Alloc/1024 // 计算增量内存使用 (KB)
					if memoryUsed > uint64(pb.MaxMem) {
						log.Printf("Out Of Memory for test case %s. Used: %dKB, Max: %dKB", tc.Input, memoryUsed, pb.MaxMem)
						select {
						case OOM <- struct{}{}:
						default:
						}
						return
					}

					// 所有检查通过，表示该测试用例通过。
					lock.Lock() // 加锁，保护 passCount 的并发修改。
					passCount++
					lock.Unlock() // 解锁。

					// 注意：这里不能在每个 goroutine 中发送 AC 信号。
					// AC 应该在所有测试用例都通过后，统一判断发送。
					// 如果提前发送 AC，可能导致 select 提前结束，而其他 goroutine 还在运行。
				}()
			}
		}
	}

	// 启动一个 goroutine 来等待所有判题 goroutine 完成，然后发送最终的 AC 或 TLE 信号。
	done := make(chan struct{})
	go func() {
		wg.Wait()   // 等待所有判题 goroutine 完成。
		close(done) // 关闭 done channel 表示所有判题任务已完成。
	}()

	// 监听不同状态的 channel，根据收到的信号更新提交记录的状态和提示信息。
	// 这里使用 select 块来决定最终的判题结果。优先级：非法代码 > 编译错误 > 运行超时 > 答案错误 > 超内存 > 答案正确。
	select {
	case <-EC: // 如果收到非法代码信号
		msg = "无效代码：包含非法操作或关键字"
		submitStatus = define.SubmitStatusInvalidCode // 状态设置为非法代码
	case <-CE: // 如果收到编译错误信号
		// msg 在编译错误时可能由 stderr 赋值，如果没有则使用默认提示
		if msg == "" {
			msg = "编译错误"
		}
		submitStatus = define.SubmitStatusCompileError // 状态设置为编译错误
	case <-TLE: // 如果收到运行超时信号
		msg = "运行超时"
		submitStatus = define.SubmitStatusTimeLimitExceeded // 状态设置为运行超时
	case <-WA: // 如果收到答案错误信号
		msg = "答案错误"
		submitStatus = define.SubmitStatusWrongAnswer // 状态设置为答案错误
	case <-OOM: // 如果收到超内存信号
		msg = "运行超内存"
		submitStatus = define.SubmitStatusMemoryLimitExceeded // 状态设置为超内存
	case <-done: // 所有测试用例执行完毕
		lock.Lock() // 加锁以安全读取 passCount
		currentPassCount := passCount
		lock.Unlock() // 解锁

		if currentPassCount == len(pb.TestCases) && len(pb.TestCases) > 0 { // 确保所有测试用例都通过且至少有一个测试用例
			msg = "答案正确"
			submitStatus = define.SubmitStatusAccepted // 状态设置为答案正确
		} else if len(pb.TestCases) == 0 && submitStatus == define.SubmitStatusPending {
			// 如果没有测试用例且之前没有其他错误，默认为通过
			msg = "无测试用例，默认正确"
			submitStatus = define.SubmitStatusAccepted
		} else {
			// 如果 some goroutine finished but not all passed, and no other error, means WA or TLE
			// This branch handles cases where 'done' is received, but not all test cases passed,
			// and no explicit WA, OOM, TLE, CE, EC signal was sent.
			// It might indicate a hidden bug or a state not explicitly handled by the channels.
			// For robustness, consider adding a 'default' case or clearer logic for partial passes.
			// For simplicity here, if not all passed and no other specific error, it implies WA or TLE.
			// Given the select order, if done is reached first and passCount is not full, it's ambiguous.
			// It's better to ensure one of the error channels is fired if any test case fails.
			// This default should ideally not be reached if error channels are properly managed.
			if submitStatus == define.SubmitStatusPending { // If no other specific error was set
				msg = "部分测试用例未通过或未知错误"
				submitStatus = define.SubmitStatusWrongAnswer // 暂时归类为答案错误
			}
		}
	case <-time.After(time.Duration(pb.MaxRuntime) * time.Millisecond * 2): // 增加一个总超时，防止判题进程卡死
		// 如果在所有测试用例执行时间的两倍内（或一个更合理的值）仍未完成，则强制判定为超时。
		// 这作为一种兜底机制，以防 `wg.Wait()` 永远不会返回。
		// 结合之前的 ctx.WithTimeout，这个总超时主要是为了处理 goroutine 调度或意外卡死的情况。
		log.Printf("Total Judging Process Timed Out after %d ms", pb.MaxRuntime*2)
		if submitStatus == define.SubmitStatusPending { // 只有在没有其他明确状态时才设置为超时
			msg = "判题系统总超时"
			submitStatus = define.SubmitStatusTimeLimitExceeded
		}
	}

	// 最终设置提交状态和消息。
	sb.Status = submitStatus
	//sb.Msg = msg // 将最终消息保存到数据库

	// 开启数据库事务，更新提交记录、用户信息和问题信息。
	// 使用事务确保数据一致性。
	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		// 保存提交记录。
		err = tx.Create(sb).Error
		if err != nil {
			log.Printf("SubmitBasic Save Error: %v", err)
			return errors.New("提交记录保存失败：" + err.Error())
		}

		// 定义要更新的字段，使用 map 方便批量更新。
		m := make(map[string]interface{})
		// 提交总数加 1。
		m["submit_num"] = gorm.Expr("submit_num + ?", 1)
		// 如果判题结果是正确，则通过数加 1。
		if sb.Status == define.SubmitStatusAccepted {
			m["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}

		// 更新 user_basic 表中用户的提交数和通过数。
		err = tx.Model(new(models.UserBasic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			log.Printf("UserBasic Modify Error: %v", err)
			return errors.New("用户数据更新失败：" + err.Error())
		}

		// 更新 problem_basic 表中问题的提交数和通过数。
		err = tx.Model(new(models.ProblemBasic)).Where("identity = ?", problemIdentity).Updates(m).Error
		if err != nil {
			log.Printf("ProblemBasic Modify Error: %v", err)
			return errors.New("问题数据更新失败：" + err.Error())
		}
		// 事务成功，返回 nil。
		return nil
	}); err != nil {
		// 若事务执行出错，返回错误信息。
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "提交判题失败：" + err.Error(),
		})
		return
	}

	// 提交成功，返回提交结果。
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"status": sb.Status, // 返回最终的判题状态码
			"msg":    msg,       // 返回最终的判题结果消息
		},
		"msg": "代码提交成功，等待判题结果", // 额外的成功提示
	})
}
