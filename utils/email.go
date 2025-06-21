package utils

import (
	"crypto/tls"                     // 导入 crypto/tls 包，用于TLS配置
	"github.com/jordan-wright/email" // 导入 email 包，用于邮件发送功能
	"math/rand"                      // 导入 math/rand 包，用于生成随机数
	"net/smtp"                       // 导入 net/smtp 包，用于SMTP认证
	"time"                           // 导入 time 包，用于时间相关的操作
)

// SendCode 函数用于发送验证码
// toUserEmail: 接收验证码的用户邮箱地址
// code: 要发送的验证码
// 返回值: 如果发送成功返回 nil，否则返回错误信息
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail() // 创建一个新的 Email 实例

	// 设置发件人信息
	// 这里的 "xieyang <xieyangyuyue@foxmail.com>" 是发件人名称和邮箱地址
	e.From = "xieyang <xieyangyuyue@foxmail.com>"

	// 设置收件人邮箱地址
	e.To = []string{toUserEmail}

	// 设置邮件主题
	e.Subject = "验证码已发送，请查收"

	// 设置邮件HTML内容，将验证码嵌入到邮件正文中
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")

	// 使用 TLS 发送邮件
	// "smtp.qq.com:465" 是 QQ 邮箱的 SMTP 服务器地址和 SSL 端口
	return e.SendWithTLS(
		"smtp.qq.com:465",
		// 使用 PlainAuth 进行身份验证
		// 第一个参数通常为空字符串
		// 第二个参数是发送邮箱地址（QQ邮箱）
		// 第三个参数是邮箱的授权码 (MailPasswd，假设已在其他地方定义并初始化)
		// 第四个参数是 SMTP 服务器的域名
		smtp.PlainAuth("", "1483618794@qq.com", MailPasswd, "smtp.qq.com"),
		&tls.Config{
			// InsecureSkipVerify 设置为 true，表示跳过服务器证书验证
			// 在生产环境中，出于安全考虑，应谨慎使用此设置或根据需要配置证书
			InsecureSkipVerify: true,
			// ServerName 指定要验证的服务器名称
			ServerName: "smtp.qq.com",
		},
	)
}

var (
	digits   = "0123456789"                          // digits 存储所有可能的数字字符，用于生成验证码
	src      = rand.NewSource(time.Now().UnixNano()) // src 是一个独立的随机源，使用当前时间的纳秒作为种子，避免全局锁和重复的随机数序列
	randPool = rand.New(src)                         // randPool 是一个基于 src 的独立随机数生成器
)

// GetRand 函数用于生成一个6位的纯数字验证码
// 返回值: 生成的6位验证码字符串
func GetRand() string {
	b := make([]byte, 6) // 创建一个长度为6的字节切片，用于存储验证码的每个数字
	for i := range b {
		// 从 digits 字符串中随机选择一个数字字符，并赋值给字节切片的当前位置
		b[i] = digits[randPool.Intn(10)]
	}
	return string(b) // 将字节切片转换为字符串并返回
}
