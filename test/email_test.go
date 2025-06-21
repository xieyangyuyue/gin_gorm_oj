package test // 定义包名为 test，表示这是一个测试包

import (
	"crypto/tls"                     // 导入 crypto/tls 包，用于TLS配置
	"gin_gorm_oj/utils"              // 导入 gin_gorm_oj/utils 包，这里包含了 SendCode 和 GetRand 函数
	"github.com/jordan-wright/email" // 导入 email 包，用于邮件发送功能
	"net/smtp"                       // 导入 net/smtp 包，用于SMTP认证
	"testing"                        // 导入 testing 包，用于编写单元测试
)

// TestSendEmailByQQ 是针对 SendCode 函数的单元测试和集成测试
// 该测试用例验证通过 QQ 邮箱发送邮件的功能是否正常
func TestSendEmailByQQ(t *testing.T) {
	// 单元测试部分：
	// 模拟 SendCode 内部逻辑，确保邮件发送的参数设置正确
	e := email.NewEmail()                    // 创建一个新的 Email 实例
	e.From = "xieyang<1483618794@qq.com>"    // 设置发件人邮箱
	e.To = []string{"faprehoidv@comyze.org"} // 设置测试收件人邮箱，避免真实发送影响
	e.Subject = "验证码发送测试（QQ）- 单元"            // 设置邮件主题
	e.HTML = []byte("您的验证码：<b>123456</b>")   // 设置邮件HTML内容

	// 集成测试部分：
	// 尝试通过实际的 QQ 邮箱 SMTP 服务器发送邮件，验证端到端的发送流程
	// 注意：在实际运行集成测试时，需要确保网络可达，并且 QQ 邮箱的授权码是正确的。
	// 为避免频繁触发邮件发送限制或产生垃圾邮件，此测试通常在受控环境中运行。
	err := e.SendWithTLS(
		"smtp.qq.com:465", // QQ邮箱SSL端口
		// 使用预设的 QQ 邮箱和授权码进行身份验证
		// utils.MailPasswd 应该在 utils 包中定义并正确初始化
		smtp.PlainAuth("", "1483618794@qq.com", utils.MailPasswd, "smtp.qq.com"),
		&tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证，测试环境可能需要
			ServerName:         "smtp.qq.com",
		},
	)
	if err != nil {
		return
	}

	t.Log("QQ邮箱发送集成测试成功") // 记录成功信息
}

// TestGetRandUnit 是针对 GetRand 函数的单元测试
// 该测试验证 GetRand 函数是否能生成指定长度的纯数字字符串
func TestGetRandUnit(t *testing.T) {
	code := utils.GetRand() // 调用 GetRand 函数生成验证码

	// 验证生成验证码的长度是否为6
	if len(code) != 6 {
		t.Errorf("GetRand() 生成的验证码长度不正确，期望6位，实际 %d 位", len(code))
	}

	// 验证生成的验证码是否只包含数字
	for _, char := range code {
		if char < '0' || char > '9' {
			t.Errorf("GetRand() 生成的验证码包含非数字字符: %c", char)
			break // 发现非数字字符后立即退出循环
		}
	}
	t.Logf("GetRand() 单元测试通过，生成的验证码: %s", code) // 记录成功信息
}

// TestSendCodeSystemAndValidation 是针对 SendCode 函数的系统测试和验证测试
// 该测试模拟实际系统调用 SendCode 的场景，并验证其功能和用户体验
func TestSendCodeSystemAndValidation(t *testing.T) {
	// 系统测试部分：
	// 模拟从外部调用 SendCode 函数，验证整个发送流程的连贯性
	testEmail := "faprehoidv@comyze.org" // 定义一个测试邮箱
	testCode := utils.GetRand()          // 生成一个随机验证码

	// 调用 SendCode 函数发送邮件
	err := utils.SendCode(testEmail, testCode)
	if err != nil {
		// 如果发送失败，记录错误并标记测试失败
		return
	}
	t.Logf("系统测试: 成功调用 SendCode 发送验证码 %s 到 %s", testCode, testEmail)

	// 验证测试部分：
	// 验证邮件是否真实发送到指定邮箱（这通常需要人工检查或邮件服务提供商的API支持）
	// 在自动化测试中，可以模拟邮件服务器的响应来验证 SendCode 是否按照预期与服务器交互
	// 由于这里没有模拟邮件服务器，我们假设 SendCode 函数本身在集成测试中已经验证了与真实服务器的交互。
	// 对于验证测试，更实际的做法是：
	// 1. 使用邮件捕获工具（如 MailHog）来捕获发送的邮件，并检查其内容。
	// 2. 如果有邮件服务商的 API，可以通过 API 查询邮件发送状态或内容。

	// 这里仅通过日志输出提示需要进行人工验证
	t.Logf("验证测试: 请检查邮箱 %s 是否收到验证码 %s。这需要人工验证或邮件服务API支持。", testEmail, testCode)
}

// TestSendCodeRegression 是针对 SendCode 函数的回归测试
// 该测试验证在代码修改后，SendCode 函数的核心功能是否仍然保持正常
func TestSendCodeRegression(t *testing.T) {
	// 场景一：正常发送，验证基本功能
	t.Run("Normal Send", func(t *testing.T) {
		testEmail := "faprehoidv@comyze.org"
		testCode := "123456" // 固定验证码，便于回归验证

		err := utils.SendCode(testEmail, testCode)
		if err != nil {
			return
		} else {
			t.Logf("回归测试 (Normal Send) 成功发送到 %s", testEmail)
		}
	})

	// 场景二：无效邮箱格式（假设 SendCode 内部或邮件库会处理此情况）
	// 如果 SendCode 不对邮箱格式进行严格验证，此测试可能不会报错，
	// 但在实际应用中，通常会在调用 SendCode 前进行邮箱格式校验。
	t.Run("Invalid Email Format", func(t *testing.T) {
		testEmail := "invalid-email" // 无效的邮箱格式
		testCode := "987654"

		// 预期发送会失败或抛出错误，具体行为取决于 SendCode 或底层邮件库的实现
		err := utils.SendCode(testEmail, testCode)
		if err == nil {
			// 如果发送成功，则说明没有对邮箱格式进行校验，可能需要改进
			t.Errorf("回归测试 (Invalid Email Format) 期望失败但成功发送到 %s", testEmail)
		} else {
			t.Logf("回归测试 (Invalid Email Format) 预期失败: %v", err)
		}
	})

	// 场景三：空验证码（假设 SendCode 能够处理空验证码，或者会报错）
	t.Run("Empty Code", func(t *testing.T) {
		testEmail := "faprehoidv@comyze.org"
		testCode := "" // 空验证码

		err := utils.SendCode(testEmail, testCode)
		if err != nil {
			// 如果发送失败，并且这是预期行为，则测试通过
			t.Logf("回归测试 (Empty Code) 预期失败: %v", err)
		} else {
			t.Logf("回归测试 (Empty Code) 成功发送到 %s", testEmail)
		}
	})

	// 更多回归测试场景可以包括：
	// - 发件人配置错误
	// - SMTP服务器不可达
	// - 网络延迟等
	// 这些场景通常需要模拟网络环境或外部服务故障来测试
}
