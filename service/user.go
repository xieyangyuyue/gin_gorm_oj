package service

import (
	"errors"
	"gin_gorm_oj/define" // 引入自定义常量
	"gin_gorm_oj/middlewares"
	"gin_gorm_oj/models"
	"gin_gorm_oj/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log" // 引入log包用于日志记录
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "用户唯一标识"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	// 获取URL查询参数中的identity
	identity := c.Query("identity")
	// 检查identity是否为空
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}

	// 创建UserBasic结构体实例用于存储查询结果
	data := new(models.UserBasic)
	// 从数据库查询用户信息，使用Omit("password")排除密码字段，防止敏感信息泄露
	err := models.DB.Omit("password").Where("identity = ? ", identity).Find(&data).Error
	// 检查数据库查询是否出错
	if err != nil {
		// 记录详细错误日志
		log.Printf("根据标识查询用户失败: %v, identity: %s\n", err, identity)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "根据标识查询用户失败：" + err.Error(),
		})
		return
	}

	// 返回查询到的用户详情
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "用户名"
// @Param password formData string false "密码"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	// 获取POST表单中的用户名
	username := c.PostForm("username")
	// 获取POST表单中的密码
	password := c.PostForm("password")
	// 检查用户名或密码是否为空
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "必填信息为空",
		})
		return
	}

	// 对用户输入的密码进行MD5加密，与数据库中存储的加密密码进行比对
	password = utils.GetMd5(password)
	// 创建UserBasic结构体实例用于存储查询结果
	data := new(models.UserBasic)
	// 查询用户是否存在且密码匹配
	err := models.DB.Where("name = ? AND password = ? ", username, password).First(&data).Error
	// 检查数据库查询是否出错
	if err != nil {
		// 如果记录未找到，则表示用户名或密码错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "用户名或密码错误",
			})
			return
		}
		// 其他数据库查询错误
		log.Printf("查询用户失败: %v, username: %s\n", err, username)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询用户失败：" + err.Error(),
		})
		return
	}

	// 用户验证成功，生成JWT Token
	token, err := middlewares.GenerateToken(data.Identity, data.Name, data.IsAdmin)
	// 检查Token生成是否出错
	if err != nil {
		log.Printf("生成Token失败: %v, identity: %s, name: %s\n", err, data.Identity, data.Name)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "生成Token失败：" + err.Error(),
		})
		return
	}

	// 返回生成的Token和用户是否为管理员的信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token":    token,
			"is_admin": data.IsAdmin,
		},
	})
}

// SendCode
// @Tags 公共方法
// @Summary 发送验证码
// @Param email formData string true "email"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	// 获取POST表单中的邮箱
	email := c.PostForm("email")
	// 检查邮箱是否为空
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	// 生成一个随机验证码
	code := utils.GetRand()
	// 将验证码存储到Redis，以邮箱为key，验证码为value，设置5分钟过期时间
	// time.Second * 300 等于 5分钟
	models.RDB.Set(c, email, code, time.Second*300)
	// 调用工具函数发送验证码到指定邮箱
	err := utils.SendCode(email, code)
	// 检查发送验证码是否出错
	if err != nil {
		log.Printf("发送验证码失败: %v, email: %s\n", err, email)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Send Code Error:" + err.Error(),
		})
		return
	}

	// 返回验证码发送成功信息
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "验证码发送成功",
	})
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param mail formData string true "邮箱"
// @Param code formData string true "验证码"
// @Param name formData string true "用户名"
// @Param password formData string true "密码"
// @Param phone formData string false "手机号"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]
func Register(c *gin.Context) {
	// 获取POST表单中的邮箱
	mail := c.PostForm("mail")
	// 获取POST表单中的用户输入的验证码
	userCode := c.PostForm("code")
	// 获取POST表单中的用户名
	name := c.PostForm("name")
	// 获取POST表单中的密码
	password := c.PostForm("password")
	// 获取POST表单中的手机号（可选）
	phone := c.PostForm("phone")

	// 检查必填参数是否为空
	if mail == "" || userCode == "" || name == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}

	// 从Redis获取系统生成的验证码
	sysCode, err := models.RDB.Get(c, mail).Result()
	// 检查从Redis获取验证码是否出错
	if err != nil {
		log.Printf("获取验证码失败：%v, mail: %s \n", err, mail)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确，请重新获取",
		})
		return
	}
	// 校验用户输入的验证码与系统生成的验证码是否一致
	if sysCode != userCode {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "验证码不正确",
		})
		return
	}

	// 检查邮箱是否已注册
	var cnt int64
	// 在UserBasic表中查询是否存在该邮箱的用户
	err = models.DB.Where("mail = ?", mail).Model(new(models.UserBasic)).Count(&cnt).Error
	// 检查数据库查询是否出错
	if err != nil {
		log.Printf("查询用户失败: %v, mail: %s\n", err, mail)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询用户失败：" + err.Error(),
		})
		return
	}
	// 如果查询到已存在该邮箱的用户，则返回邮箱已注册信息
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该邮箱已被注册",
		})
		return
	}

	// 生成用户唯一标识
	userIdentity := utils.GetUUID()
	// 创建UserBasic结构体实例，并填充用户数据
	data := &models.UserBasic{
		Identity:  userIdentity,
		Name:      name,
		Password:  utils.GetMd5(password), // 对密码进行MD5加密
		Phone:     phone,
		Mail:      mail,
		CreatedAt: models.MyTime(time.Now()), // 设置创建时间
		UpdatedAt: models.MyTime(time.Now()), // 设置更新时间
	}
	// 将用户数据插入数据库
	err = models.DB.Create(data).Error
	// 检查数据库插入是否出错
	if err != nil {
		log.Printf("创建用户失败: %v, name: %s, mail: %s\n", err, name, mail)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "创建用户失败：" + err.Error(),
		})
		return
	}

	// 用户注册成功，生成JWT Token
	token, err := middlewares.GenerateToken(userIdentity, name, data.IsAdmin)
	// 检查Token生成是否出错
	if err != nil {
		log.Printf("生成Token失败: %v, identity: %s, name: %s\n", err, userIdentity, name)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "生成Token失败：" + err.Error(),
		})
		return
	}
	// 返回生成的Token
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": token,
		},
	})
}

// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
	// 获取URL查询参数中的size，如果未提供则使用define.DefaultSize
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	// 获取URL查询参数中的page，如果未提供则使用define.DefaultPage
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	// 检查分页参数转换是否出错
	if err != nil {
		log.Println("分页参数转换错误:", err)
		c.JSON(http.StatusOK, gin.H{ // 统一错误返回格式
			"code": -1,
			"msg":  "分页参数错误：" + err.Error(),
		})
		return
	}
	// 计算数据库查询的偏移量
	page = (page - 1) * size

	// 用于存储用户总数
	var count int64
	// 用于存储查询到的用户列表
	list := make([]*models.UserBasic, 0)
	// 查询UserBasic表的总记录数，并按通过数降序、提交数升序排序，然后进行分页查询
	err = models.DB.Model(new(models.UserBasic)).Count(&count).Order("pass_num DESC, submit_num ASC").
		Offset(page).Limit(size).Find(&list).Error
	// 检查数据库查询是否出错
	if err != nil {
		log.Printf("获取排行榜失败: %v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "获取排行榜失败：" + err.Error(),
		})
		return
	}

	// 对邮箱进行脱敏处理（例如：a**@example.com）
	for _, v := range list {
		// 将邮箱地址按"@"分割
		mail := strings.Split(v.Mail, "@")
		// 如果分割后至少有两部分（即包含@符号和域名）
		if len(mail) >= 2 {
			// 将邮箱名部分的首字符保留，其余替换为"**"，然后拼接域名
			v.Mail = string(mail[0][0]) + "**@" + mail[1]
		}
	}

	// 返回排行榜列表和总用户数
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"list":  list,  // 排行榜列表
			"count": count, // 总用户数
		},
	})
}
