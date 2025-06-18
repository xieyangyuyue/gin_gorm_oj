package test

//
//import (
//	"crypto/tls"
//	"gin_gorm_oj/utils"
//	"github.com/jordan-wright/email"
//	"net/smtp"
//	"testing"
//)
//
//// 使用 QQ 邮箱发送测试
//func TestSendEmailByQQ(t *testing.T) {
//	e := email.NewEmail()
//	e.From = "xieyang<1483618794@qq.com>"       // 替换为你的QQ邮箱
//	e.To = []string{"xieyangyuyue@foxmail.com"} // 接收方邮箱
//	e.Subject = "验证码发送测试1（QQ）"
//	e.HTML = []byte("您的验证码：<b>123456</b>")
//
//	// QQ邮箱SMTP配置
//	err := e.SendWithTLS(
//		"smtp.qq.com:465", // QQ邮箱SSL端口
//		smtp.PlainAuth("", "1483618794@qq.com", utils.MailPasswd, "smtp.qq.com"), // 密码需使用QQ邮箱的授权码
//		&tls.Config{
//			InsecureSkipVerify: true, // 跳过证书验证（根据环境需要调整）
//			ServerName:         "smtp.qq.com",
//		},
//	)
//	if err != nil {
//		t.Fatal("QQ邮箱发送失败:", err)
//	}
//}
//
////// 使用 Google 邮箱发送测试
////func TestSendEmailByGoogle(t *testing.T) {
////	e := email.NewEmail()
////	e.From = "Your Name <your-gmail@gmail.com>" // 替换为你的Gmail邮箱
////	e.To = []string{"recipient@example.com"}    // 接收方邮箱
////	e.Subject = "验证码发送测试（Google）"
////	e.HTML = []byte("您的验证码：<b>123456</b>")
////
////	// Google邮箱SMTP配置
////	err := e.SendWithTLS(
////		"smtp.gmail.com:465", // Google邮箱SSL端口
////		smtp.PlainAuth("", "your-gmail@gmail.com", "your-app-specific-password", "smtp.gmail.com"), // 使用应用专用密码
////		&tls.Config{
////			InsecureSkipVerify: false, // 生产环境应设为false并配置证书
////			ServerName:         "smtp.gmail.com",
////		},
////	)
////	if err != nil {
////		t.Fatal("Google邮箱发送失败:", err)
////	}
////}
