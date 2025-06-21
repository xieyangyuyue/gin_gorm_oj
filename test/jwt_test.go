package test

import (
	"errors"
	"fmt"                          // 导入 fmt 包，用于格式化输出和错误信息
	"gin_gorm_oj/middlewares"      // 导入你的 middlewares 包，包含 JWT 相关的函数。请确保路径正确
	"gin_gorm_oj/utils"            // 导入你的 utils 包，用于获取 JWT 密钥。请确保路径正确
	"github.com/golang-jwt/jwt/v5" // 导入 go-jwt/jwt v5 库，用于 JWT 操作
	"strings"
	"sync"    // 导入 sync 包，用于并发测试中的同步原语，如 WaitGroup
	"testing" // 导入 Go 语言内置的 testing 包，用于编写测试函数
	"time"    // 导入 time 包，用于处理时间相关的操作
)

// TestGenerateAndAnalyseToken_Unit 是一个单元测试函数，用于测试 JWT 的生成和解析的端到端流程。
// 覆盖了基本的生成、解析、过期时间、以及核心声明的验证。
func TestGenerateAndAnalyseToken_Unit(t *testing.T) { // 定义一个名为 TestGenerateAndAnalyseToken_Unit 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 定义测试用例
	tests := []struct {
		name          string        // 测试用例名称
		identity      string        // 用户身份
		userName      string        // 用户名
		isAdmin       int           // 是否管理员
		expiresIn     time.Duration // Token 过期时间（正值表示未来过期，负值表示已过期）
		expectedError error         // 预期 AnalyseToken 返回的错误类型 (nil 表示无错误)
	}{
		{
			name:          "Valid Token - Admin User",
			identity:      "user_admin_1",
			userName:      "AdminUser",
			isAdmin:       1,
			expiresIn:     1 * time.Hour, // Token 在未来1小时后过期
			expectedError: nil,
		},
		{
			name:          "Valid Token - Regular User",
			identity:      "user_normal_2",
			userName:      "NormalUser",
			isAdmin:       0,
			expiresIn:     2 * time.Hour, // Token 在未来2小时后过期
			expectedError: nil,
		},
		{
			name:          "Expired Token",
			identity:      "user_expired",
			userName:      "ExpiredUser",
			isAdmin:       0,
			expiresIn:     -5 * time.Minute,    // Token 在5分钟前已过期
			expectedError: jwt.ErrTokenExpired, // 预期返回 Token 已过期的错误
		},
		//	//可以添加更多测试用例，例如空身份、空用户名等（根据业务需求决定是否允许）
		//	{
		//		name:     "Token Not Valid Yet",
		//		identity: "user_future",
		//		userName: "FutureUser",
		//		isAdmin:  0,
		//		// 为了测试 NotValidYet，我们需要在 GenerateToken 中设置 NotBefore 字段，这里假设 GenerateToken 已经有此功能
		//		// 假设 Token 5分钟后才生效
		//		expiresIn: 1 * time.Hour, // 总体过期时间
		//		notBeforeOffset: 5 * time.Minute, // 5分钟后才生效
		//		expectedError: jwt.ErrTokenNotValidYet,
		//	},
	}

	for _, tt := range tests { // 遍历所有测试用例
		t.Run(tt.name, func(t *testing.T) { // 为每个测试用例创建一个子测试
			// **关键修改：在生成 Token 时直接控制其过期状态**
			// 构建 UserClaims，根据 expiresIn 设置 ExpiresAt
			claim := &middlewares.UserClaims{
				Identity: tt.identity,
				Name:     tt.userName,
				IsAdmin:  tt.isAdmin,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer: "gin-gorm-oj",
					//NewNumericDate(time.Now().UTC().Add(tt.expiresIn)) 直接设置未来的或过去的时间
					ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(tt.expiresIn)),
					// 如果需要测试 jwt.ErrTokenNotValidYet，可以在这里添加 NotBefore
					//NotBefore: jwt.NewNumericDate(time.Now().UTC().Add(tt.notBeforeOffset)),
				},
			}
			// 使用 HMAC SHA256 签名方法和自定义声明创建新的 JWT Token 对象。
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			// 使用预定义的密钥 myKey 对 Token 进行签名，并返回签名后的字符串。
			tokenString, err := token.SignedString([]byte(utils.JwtKey)) // 直接使用 utils.JwtKey，避免依赖 middlewares.myKey
			if err != nil {                                              // 检查生成 Token 是否出错
				t.Fatalf("GenerateToken unexpected error: %v", err) // 如果出错，则致命错误并退出测试
			}
			if tokenString == "" { // 检查生成的 Token 字符串是否为空
				t.Fatal("GenerateToken returned empty token string") // 如果为空，则致命错误
			}
			t.Logf("Generated Token for %s: %s", tt.userName, tokenString) // 打印生成的 Token

			// 解析 Token
			actualClaims, actualErr := middlewares.AnalyseToken(tokenString) // 调用 middlewares 包的 AnalyseToken 函数

			// **核心验证逻辑修改：根据 expectedError 来断言**
			if tt.expectedError != nil { // 如果预期有错误
				if actualErr == nil { // 实际没有错误
					t.Errorf("AnalyseToken expected error '%v' but got none", tt.expectedError) // 报错
				} else if !errors.Is(actualErr, tt.expectedError) && !strings.Contains(actualErr.Error(), tt.expectedError.Error()) {
					// 由于 AnalyseToken 中我们用了 fmt.Errorf("token is expired: %w", err) 封装了错误
					// 所以这里需要检查实际错误是否包含预期的错误信息
					t.Errorf("AnalyseToken expected error containing '%v', but got '%v'", tt.expectedError, actualErr) // 报错
				}
				return // 结束当前子测试
			}

			// 如果不预期错误 (即 expectedError 为 nil)
			if actualErr != nil { // 但实际发生了错误
				t.Fatalf("AnalyseToken unexpected error: %v", actualErr) // 致命错误并退出测试
			}
			if actualClaims == nil { // 如果解析结果为空
				t.Fatal("AnalyseToken returned nil claims") // 致命错误
			}

			// 验证 Token 中的声明是否正确（仅在预期无错误时进行）
			if actualClaims.Identity != tt.identity { // 验证 Identity 字段
				t.Errorf("Expected identity %q, got %q", tt.identity, actualClaims.Identity) // 如果不匹配，则报错
			}
			if actualClaims.Name != tt.userName { // 验证 Name 字段
				t.Errorf("Expected name %q, got %q", tt.userName, actualClaims.Name) // 如果不匹配，则报错
			}
			if actualClaims.IsAdmin != tt.isAdmin { // 验证 IsAdmin 字段
				t.Errorf("Expected isAdmin %d, got %d", tt.isAdmin, actualClaims.IsAdmin) // 如果不匹配，则报错
			}
			// 对于未过期的 Token，可以额外验证 ExpiresAt 是否在未来
			if time.Now().UTC().After(actualClaims.ExpiresAt.Time) {
				t.Errorf("Token should not be expired, but its ExpiresAt (%v) is in the past.", actualClaims.ExpiresAt.Time)
			}

			t.Logf("Successfully parsed and validated token for %s", tt.userName) // 记录成功信息
		})
	}
}

// TestAnalyseToken_InvalidSignature 是一个单元测试函数，用于测试解析带有无效签名的 JWT。
func TestAnalyseToken_InvalidSignature(t *testing.T) { // 定义一个名为 TestAnalyseToken_InvalidSignature 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 生成一个使用正确密钥签名的 Token
	_, err := middlewares.GenerateToken("test_user", "TestUser", 0) // 生成一个有效 Token
	if err != nil {
		t.Fatalf("Failed to generate valid token: %v", err) // 如果生成失败，则致命错误
	}

	// 篡改密钥，创建一个无效签名的 Token 字符串 (这里只是一个概念，实际中难以直接篡改签名)
	// 最直接的方式是使用一个不同的密钥生成另一个 Token 来模拟“无效签名”
	maliciousKey := []byte("a_different_secret_key_that_is_wrong") // 定义一个错误的密钥
	// 使用错误的密钥对相同的 claims 进行签名，模拟无效签名
	fakeClaims := &middlewares.UserClaims{
		Identity: "test_user",
		Name:     "TestUser",
		IsAdmin:  0,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		},
	}
	fakeToken := jwt.NewWithClaims(jwt.SigningMethodHS256, fakeClaims) // 使用错误的密钥生成一个“伪造”的 Token
	fakeTokenString, err := fakeToken.SignedString(maliciousKey)       // 用错误的密钥签名
	if err != nil {
		t.Fatalf("Failed to generate fake token: %v", err) // 如果生成失败，则致命错误
	}

	// 尝试解析这个无效签名的 Token
	_, err = middlewares.AnalyseToken(fakeTokenString) // 尝试解析使用错误密钥签名的 Token
	if err == nil {                                    // 如果没有报错
		t.Error("AnalyseToken expected an error for invalid signature, but got none") // 报错，因为预期应该报错
	} else if err.Error() != "failed to parse token: crypto/ecdsa: verification error" && // 兼容 go-jwt v5 的不同错误信息
		err.Error() != "failed to parse token: signature verification failed" && // 早期版本或不同错误信息
		err.Error() != "failed to parse token: token is malformed: signature verification failed" &&
		err.Error() != "token signature is invalid: token signature is invalid: signature is invalid" { // 兼容更具体的错误
		t.Errorf("AnalyseToken got unexpected error for invalid signature: %v", err) // 报错，因为错误信息不符预期
	}
	t.Logf("Successfully rejected token with invalid signature: %v", err) // 记录成功拒绝
}

// TestAnalyseToken_ExpiredToken 是一个单元测试函数，用于测试解析已过期的 JWT。
func TestAnalyseToken_ExpiredToken(t *testing.T) { // 定义一个名为 TestAnalyseToken_ExpiredToken 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 生成一个立即过期的 Token
	_, err := middlewares.GenerateToken("expired_user", "ExpiredUser", 0) // 生成一个 Token
	if err != nil {
		t.Fatalf("Failed to generate token for expiration test: %v", err) // 如果生成失败，则致命错误
	}

	// 模拟 Token 过期
	// 实际情况中，由于 GenerateToken 会设置过期时间，这里需要修改 JWT.go 让过期时间可控，
	// 或者直接构造一个过期时间的 Token 字符串。
	// 为了简化，我们直接调用 GenerateToken 时传入一个极短的过期时间，然后在测试中等待。

	// 直接构造一个即将过期的 Token（这里为了不修改 GenerateToken 函数，模拟一下）
	claims := &middlewares.UserClaims{
		Identity: "expired_user",
		Name:     "ExpiredUser",
		IsAdmin:  0,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-1 * time.Minute)), // 设置为一分钟前过期
			Issuer:    "gin-gorm-oj",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)   // 创建 Token
	tokenString, err := token.SignedString([]byte(utils.JwtKey)) // 签名 Token
	if err != nil {
		t.Fatalf("Failed to sign expired token: %v", err) // 如果签名失败，则致命错误
	}

	// 尝试解析已过期的 Token
	_, err = middlewares.AnalyseToken(tokenString) // 尝试解析过期 Token
	if err == nil {                                // 如果没有报错
		t.Error("AnalyseToken expected an error for expired token, but got none") // 报错，因为预期应该报错
	} else if err.Error() != "token is expired: token has invalid claims: token is expired" { // 检查错误信息是否符合预期
		t.Errorf("AnalyseToken got unexpected error for expired token: %v", err) // 报错，因为错误信息不符预期
	}
	t.Logf("Successfully rejected expired token: %v", err) // 记录成功拒绝
}

// BenchmarkGenerateToken 是基准测试函数，用于评估 GenerateToken 的性能。
func BenchmarkGenerateToken(b *testing.B) { // 定义一个名为 BenchmarkGenerateToken 的基准测试函数
	b.ResetTimer() // 重置计时器，排除 Setup 代码的执行时间

	b.Run("Sequential Generation", func(b *testing.B) { // 运行子基准测试，顺序生成 Token
		for i := 0; i < b.N; i++ { // 循环 b.N 次，b.N 会根据测试运行时间自动调整
			_, _ = middlewares.GenerateToken("benchmark_user", "BenchmarkUser", 0) // 生成 Token，忽略返回值
		}
	})

	b.Run("Concurrent Generation (10 Goroutines)", func(b *testing.B) { // 运行子基准测试，并发生成 Token
		var wg sync.WaitGroup // 声明一个 WaitGroup 用于等待所有 goroutine 完成
		numGoroutines := 10   // 设定并发的 goroutine 数量
		if b.N < numGoroutines {
			b.N = numGoroutines // 确保至少每个 goroutine 执行一次
		}
		iterationsPerGoroutine := b.N / numGoroutines // 每个 goroutine 需要生成的 Token 数量

		b.ResetTimer()                       // 在并发测试的内部再次重置计时器
		for i := 0; i < numGoroutines; i++ { // 启动指定数量的 goroutine
			wg.Add(1)   // 增加 WaitGroup 计数
			go func() { // 启动一个新的 goroutine
				defer wg.Done()                               // goroutine 结束后调用 Done() 减少 WaitGroup 计数
				for j := 0; j < iterationsPerGoroutine; j++ { // 每个 goroutine 内部循环生成 Token
					_, _ = middlewares.GenerateToken("benchmark_user", "BenchmarkUser", 0) // 生成 Token
				}
			}()
		}
		wg.Wait() // 等待所有 goroutine 完成
	})
}

// BenchmarkAnalyseToken 是基准测试函数，用于评估 AnalyseToken 的性能。
func BenchmarkAnalyseToken(b *testing.B) {
	// 定义一个名为 BenchmarkAnalyseToken 的基准测试函数
	// 提前生成一个 Token 供解析使用，避免在循环内重复生成影响解析性能测试
	tokenString, err := middlewares.GenerateToken("benchmark_user", "BenchmarkUser", 0) // 生成一个 Token
	if err != nil {
		b.Fatalf("Failed to generate token for benchmark: %v", err) // 如果生成失败，则致命错误
	}

	b.ResetTimer() // 重置计时器

	b.Run("Sequential Analysis", func(b *testing.B) { // 运行子基准测试，顺序解析 Token
		for i := 0; i < b.N; i++ { // 循环 b.N 次
			_, _ = middlewares.AnalyseToken(tokenString) // 解析 Token，忽略返回值
		}
	})

	b.Run("Concurrent Analysis (10 Goroutines)", func(b *testing.B) { // 运行子基准测试，并发解析 Token
		var wg sync.WaitGroup // 声明一个 WaitGroup
		numGoroutines := 10   // 设定并发 goroutine 数量
		if b.N < numGoroutines {
			b.N = numGoroutines
		}
		iterationsPerGoroutine := b.N / numGoroutines // 每个 goroutine 需要解析的 Token 数量

		b.ResetTimer()                       // 在并发测试的内部再次重置计时器
		for i := 0; i < numGoroutines; i++ { // 启动指定数量的 goroutine
			wg.Add(1)   // 增加 WaitGroup 计数
			go func() { // 启动新的 goroutine
				defer wg.Done()                               // goroutine 结束后调用 Done()
				for j := 0; j < iterationsPerGoroutine; j++ { // 每个 goroutine 内部循环解析 Token
					_, _ = middlewares.AnalyseToken(tokenString) // 解析 Token
				}
			}()
		}
		wg.Wait() // 等待所有 goroutine 完成
	})
}

// TestJWT_Integration 是一个集成测试函数，用于验证 JWT 功能与应用其他部分的集成（例如用户服务）。
// 在此示例中，我们模拟与一个简单的“用户查找”逻辑的集成。
func TestJWT_Integration(t *testing.T) { // 定义一个名为 TestJWT_Integration 的测试函数
	t.Log("开始集成测试：JWT与模拟用户服务的集成") // 记录测试开始信息

	// 模拟用户数据（通常来自数据库或用户服务）
	mockUsers := map[string]struct {
		Name    string
		IsAdmin int
	}{
		"user_id_123":  {"Alice", 0},
		"admin_id_456": {"Bob", 1},
	}

	// 模拟用户查找函数
	findUserByIdentity := func(identity string) (string, int, error) {
		if user, ok := mockUsers[identity]; ok {
			return user.Name, user.IsAdmin, nil
		}
		return "", 0, fmt.Errorf("user not found: %s", identity) // 模拟用户不存在
	}

	// 集成测试用例1: 生成并解析一个合法用户的 Token，然后尝试“查找”该用户
	t.Run("Generate and Validate Real User Token", func(t *testing.T) { // 运行子测试
		identity := "user_id_123"
		userName := "Alice"
		isAdmin := 0

		tokenString, err := middlewares.GenerateToken(identity, userName, isAdmin) // 生成 Token
		if err != nil {
			t.Fatalf("Integration test failed: GenerateToken error: %v", err) // 如果生成失败，则致命错误
		}

		claims, err := middlewares.AnalyseToken(tokenString) // 解析 Token
		if err != nil {
			t.Fatalf("Integration test failed: AnalyseToken error: %v", err) // 如果解析失败，则致命错误
		}

		// 验证解析出的身份是否可以在模拟用户服务中找到
		foundName, foundIsAdmin, err := findUserByIdentity(claims.Identity) // 使用解析出的 Identity 查找用户
		if err != nil {
			t.Errorf("Integration test failed: User %s not found in mock service: %v", claims.Identity, err) // 如果未找到用户，则报错
		}
		if foundName != claims.Name || foundIsAdmin != claims.IsAdmin { // 验证查找出的用户与 Token 声明是否一致
			t.Errorf("Integration test failed: User data mismatch. Expected name %s, got %s; Expected admin %d, got %d",
				claims.Name, foundName, claims.IsAdmin, foundIsAdmin) // 如果不一致，则报错
		}
		t.Logf("Integration test for %s successful.", identity) // 记录成功信息
	})

	// 集成测试用例2: 尝试解析一个不存在用户的 Token (但 Token 本身是有效的)
	t.Run("Analyse Token for Non-existent User", func(t *testing.T) { // 运行子测试
		nonExistentIdentity := "non_existent_user_999"
		nonExistentName := "GhostUser"
		nonExistentAdmin := 0

		tokenString, err := middlewares.GenerateToken(nonExistentIdentity, nonExistentName, nonExistentAdmin) // 生成 Token
		if err != nil {
			t.Fatalf("Integration test failed: GenerateToken for non-existent user error: %v", err) // 如果生成失败，则致命错误
		}

		claims, err := middlewares.AnalyseToken(tokenString) // 解析 Token
		if err != nil {
			t.Fatalf("Integration test failed: AnalyseToken for non-existent user error: %v", err) // 如果解析失败，则致命错误
		}

		// 验证查找不存在的用户
		_, _, err = findUserByIdentity(claims.Identity) // 尝试查找用户
		if err == nil {
			t.Errorf("Integration test failed: Expected user %s not to be found, but it was.", claims.Identity) // 如果找到用户，则报错
		} else {
			t.Logf("Integration test: Successfully identified non-existent user '%s' (%v)", claims.Identity, err) // 记录成功识别
		}
	})
	t.Log("集成测试完成。") // 记录测试完成信息
}

// TestJWT_SystemScenario 是一个系统测试函数，模拟高并发环境下 JWT 的生成和解析。
// 这通常涉及更复杂的系统设置，例如 Web 服务器、数据库等，这里仅为概念性模拟。
func TestJWT_SystemScenario(t *testing.T) { // 定义一个名为 TestJWT_SystemScenario 的测试函数
	t.Log("开始系统测试：高并发 JWT 生成与解析场景") // 记录测试开始信息

	numRequests := 1000 // 模拟的请求数量
	numWorkers := 100   // 并发工作者（goroutine）的数量

	var wg sync.WaitGroup                    // 声明 WaitGroup
	errChan := make(chan error, numRequests) // 用于收集并发操作中的错误

	testUser := struct {
		identity string
		name     string
		isAdmin  int
	}{
		identity: "sys_test_user",
		name:     "SystemUser",
		isAdmin:  0,
	}

	for i := 0; i < numWorkers; i++ { // 启动多个并发工作者
		wg.Add(1)               // 增加 WaitGroup 计数
		go func(workerID int) { // 启动新的 goroutine
			defer wg.Done() // goroutine 结束后调用 Done()

			for j := 0; j < numRequests/numWorkers; j++ { // 每个工作者处理一部分请求
				// 步骤1: 生成 Token
				tokenString, err := middlewares.GenerateToken(testUser.identity, testUser.name, testUser.isAdmin) // 生成 Token
				if err != nil {
					errChan <- fmt.Errorf("worker %d: generate token error: %w", workerID, err) // 将错误发送到错误通道
					continue
				}

				// 步骤2: 解析 Token
				claims, err := middlewares.AnalyseToken(tokenString) // 解析 Token
				if err != nil {
					errChan <- fmt.Errorf("worker %d: analyse token error: %w", workerID, err) // 将错误发送到错误通道
					continue
				}

				// 步骤3: 验证解析结果
				if claims.Identity != testUser.identity || claims.Name != testUser.name || claims.IsAdmin != testUser.isAdmin { // 验证声明
					errChan <- fmt.Errorf("worker %d: claims mismatch. Expected %v, got %v", workerID, testUser, claims) // 如果不匹配，则发送错误
					continue
				}
			}
			t.Logf("工作者 %d 完成 %d 个请求", workerID, numRequests/numWorkers) // 记录工作者完成信息
		}(i) // 传入 workerID
	}

	wg.Wait()      // 等待所有工作者完成
	close(errChan) // 关闭错误通道

	// 检查是否有错误发生
	for err := range errChan { // 遍历错误通道中的所有错误
		t.Error(err) // 记录错误
	}
	if t.Failed() { // 检查是否有测试失败
		t.Error("系统测试失败：在高并发场景下发现错误。") // 如果有失败，则报错
	} else {
		t.Log("系统测试完成：高并发 JWT 生成与解析通过。") // 记录测试通过
	}
}

// TestJWT_Validation 是一个验证测试函数，用于严格验证 JWT 的各个方面。
func TestJWT_Validation(t *testing.T) { // 定义一个名为 TestJWT_Validation 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 验证 Token 的结构和格式
	t.Run("Token Format Validation", func(t *testing.T) { // 运行子测试，验证 Token 格式
		tokenString, err := middlewares.GenerateToken("validation_user", "ValidationUser", 0) // 生成 Token
		if err != nil {                                                                       // 检查生成 Token 是否出错
			t.Fatalf("Failed to generate token: %v", err) // 如果生成失败，则致命错误并退出测试
		}
		// 使用 strings.Split 函数将 Token 字符串按 "." 分割成三部分：Header, Payload, Signature
		parts := strings.Split(tokenString, ".") // 分割 Token 字符串为三部分
		if len(parts) != 3 {                     // 检查分割后的部分数量是否为 3
			t.Errorf("Invalid token format: expected 3 parts, got %d", len(parts)) // 如果不为 3，则报错
		}
		t.Logf("Token format validation successful. Token: %s", tokenString) // 记录 Token 格式验证成功信息
	})

	// 验证不同签名算法的兼容性 (AnalyseToken 目前只支持 HS256, 但如果未来支持其他算法需要扩展)
	t.Run("Unsupported Signing Method Rejection", func(t *testing.T) { // 运行子测试，验证不支持的签名方法是否被拒绝
		// 构造一个使用不同签名方法的 Token (例如 RS256)，期望被拒绝
		// 注意：这里需要使用一个 RS256 的私钥来签名，为了简化，我们仅模拟其 Token 结构。
		// 实际中需要生成 RSA 密钥对。
		// 这里只是为了触发 AnalyseToken 中的 `unexpected signing method` 错误

		// 为了测试方便，我们直接构造一个带有不同 'alg' Header 的 Token 字符串来测试
		// 例如：eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVuc3VwcG9ydGVkX2FsZ191c2VyIiwibmFtZSI6IlVuc3VwcG9ydGVkQWxnVXNlciIsImlzX2FkbWluIjowLCJpZXNzIjoiaW4tZ29ybS1vaiIsImV4cCI6MTY3ODgyNzYwMH0.signature
		testTokenString := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVuc3VwcG9ydGVkX2FsZ191c2VyIiwibmFtZSI6IlVuc3VwcG9ydGVkQWxnVXNlciIsImlzX2FkbWluIjowLCJpc3MiOiJnaW4tZ29ybS1vaiIsImV4cCI6MTY3ODgyNzYwMH0.some_fake_signature"

		_, err := middlewares.AnalyseToken(testTokenString) // 尝试解析这个带有不支持签名方法的 Token
		if err == nil {                                     // 如果没有报错
			t.Error("Validation test failed: Expected error for unsupported signing method, but got none.") // 报错，因为预期应该报错
		} else if err.Error() != "failed to parse token: token is unverifiable: error while executing keyfunc: unexpected signing method: RS256" { // 检查错误信息是否符合预期
			t.Errorf("Validation test failed: Unexpected error message for unsupported signing method: %v", err) // 报错，因为错误信息不符预期
		}
		t.Log("Successfully rejected unsupported signing method.") // 记录成功拒绝
	})

	// 验证非法 Claims 类型是否被拒绝
	t.Run("Invalid Claims Type Rejection", func(t *testing.T) { // 运行子测试，验证非法 Claims 类型是否被拒绝
		// 构造一个 JWT，但其 claims 不是 UserClaims 类型
		type AnotherClaims struct { // 定义一个不同的 Claims 结构体
			SomeField            string `json:"some_field"` // 包含一个自定义字段
			jwt.RegisteredClaims        // 嵌入 JWT 标准注册声明
		}
		anotherClaims := &AnotherClaims{ // 创建 AnotherClaims 实例
			SomeField: "test", // 填充自定义字段
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)), // 设置过期时间
				Issuer:    "test-app",                                          // 设置发行者
			},
		}
		// 使用 HMAC SHA256 签名方法和 AnotherClaims 创建新的 JWT Token 对象
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, anotherClaims)
		// 使用密钥对 Token 进行签名，并返回签名后的字符串
		tokenString, err := token.SignedString([]byte(utils.JwtKey))
		if err != nil { // 检查签名是否出错
			t.Fatalf("Failed to sign token with another claims type: %v", err) // 如果出错，则致命错误
		}

		_, err = middlewares.AnalyseToken(tokenString) // 尝试解析这个带有非法 Claims 类型的 Token
		if err == nil {                                // 如果没有报错
			t.Error("Validation test failed: Expected error for invalid claims type, but got none.") // 报错，因为预期应该报错
		} else if err.Error() != "invalid token or claims type mismatch" { // 检查错误信息是否符合预期
			t.Errorf("Validation test failed: Unexpected error message for invalid claims type: %v", err) // 报错，因为错误信息不符预期
		}
		t.Log("Successfully rejected invalid claims type.") // 记录成功拒绝
	})
}

// TestJWT_Regression 是一个回归测试函数，确保在未来的代码变更后，JWT 生成和解析的行为仍然一致且正确。
func TestJWT_Regression(t *testing.T) { // 定义一个名为 TestJWT_Regression 的测试函数
	t.Parallel() // 标记该测试可以并行运行

	// 回归测试通常会重复运行现有的关键单元测试和集成测试，以确保没有引入新的缺陷。
	// 这里我们重复调用一些核心的验证逻辑。

	t.Run("Consistent Token Generation for Known Input", func(t *testing.T) { // 运行子测试，验证已知输入的 Token 生成一致性
		identity := "reg_user_1"
		name := "RegressionUser"
		isAdmin := 0
		// 注意：由于过期时间是动态的，这里不能直接断言完整的 token 字符串。
		// 但我们可以验证其核心信息和结构。
		tokenString, err := middlewares.GenerateToken(identity, name, isAdmin) // 生成 Token
		if err != nil {
			t.Fatalf("Regression test failed: GenerateToken error: %v", err) // 如果生成失败，则致命错误
		}

		claims, err := middlewares.AnalyseToken(tokenString) // 解析 Token
		if err != nil {
			t.Fatalf("Regression test failed: AnalyseToken error: %v", err) // 如果解析失败，则致命错误
		}
		if claims.Identity != identity || claims.Name != name || claims.IsAdmin != isAdmin { // 验证声明
			t.Errorf("Regression test failed: Claims mismatch after regression. Expected %s/%s/%d, got %s/%s/%d",
				identity, name, isAdmin, claims.Identity, claims.Name, claims.IsAdmin) // 如果不匹配，则报错
		}
		t.Log("Regression test: Token generation and basic claims parsing are consistent.") // 记录一致性
	})

	t.Run("AnalyseToken Error Consistency (Expired Token)", func(t *testing.T) { // 运行子测试，验证过期 Token 错误的一致性
		// 构造一个过期 Token
		claims := &middlewares.UserClaims{
			Identity: "expired_reg_user",
			Name:     "ExpiredRegUser",
			IsAdmin:  0,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-5 * time.Minute)), // 5分钟前过期
				Issuer:    "gin-gorm-oj",
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)   // 创建 Token
		tokenString, err := token.SignedString([]byte(utils.JwtKey)) // 签名 Token
		if err != nil {
			t.Fatalf("Failed to sign expired token for regression: %v", err) // 如果签名失败，则致命错误
		}

		_, err = middlewares.AnalyseToken(tokenString) // 尝试解析
		if err == nil {
			t.Error("Regression test failed: Expected error for expired token, but got none.") // 如果没有报错，则报错
		} else if err.Error() != "token is expired: token has invalid claims: token is expired" { // 检查错误信息是否保持一致
			t.Errorf("Regression test failed: Expired token error message changed. Expected 'token is expired: Token has expired', got '%v'", err) // 如果不一致，则报错
		}
		t.Log("Regression test: Expired token error message is consistent.") // 记录一致性
	})
}
