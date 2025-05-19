package utils

import (
	"crypto/tls"
	"github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"time"
)

// SendCode
// 发送验证码
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "xieyang <xieyangyuyue@foxmail.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送，请查收"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")
	// QQ邮箱SMTP配置
	return e.SendWithTLS(
		"smtp.qq.com:465", // QQ邮箱SSL端口
		smtp.PlainAuth("", "1483618794@qq.com", MailPasswd, "smtp.qq.com"), // 密码需使用QQ邮箱的授权码
		&tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证（根据环境需要调整）
			ServerName:         "smtp.qq.com",
		},
	)
}

var (
	digits   = "0123456789"                          // 预定义数字字符集
	src      = rand.NewSource(time.Now().UnixNano()) // 独立随机源，避免全局锁[6,7](@ref)
	randPool = rand.New(src)                         // 创建独立的随机数生成器
)

// GetRand 生成6位纯数字验证码
func GetRand() string {
	b := make([]byte, 6) // 预分配字节切片，避免动态扩容
	for i := range b {
		b[i] = digits[randPool.Intn(10)] // 直接访问预定义字符集[6](@ref)
	}
	return string(b)
}
