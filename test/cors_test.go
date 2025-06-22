package test

import (
	"gin_gorm_oj/middlewares"
	"github.com/gin-gonic/gin" // 导入 Gin 框架
	"net/http"                 // 导入 http 包，用于 http 客户端、服务器和状态码
	"net/http/httptest"        // 导入 httptest 包，它提供了用于 HTTP 测试的实用工具
	"testing"                  // 导入 testing 包，用于编写 Go 语言的测试
)

// TestCorsMiddleware 是针对我们 Cors 中间件的测试套件。
func TestCorsMiddleware(t *testing.T) {
	// 设置 Gin 为测试模式，这样可以减少日志输出。
	gin.SetMode(gin.TestMode)

	// --- 测试用例 1: 预检请求 (OPTIONS) ---
	t.Run("PreflightRequest", func(t *testing.T) {
		// 创建一个 Gin 路由器实例。
		router := gin.New()
		// 将我们的 Cors 中间件应用到路由器。
		router.Use(middlewares.Cors())
		// 为 OPTIONS 方法注册一个测试路由。
		router.OPTIONS("/test", func(c *gin.Context) {
			c.Status(http.StatusOK) // 在实际情况中，预检请求不应到达这里。
		})

		// 创建一个响应记录器，用于捕获服务器的响应。
		w := httptest.NewRecorder()
		// 创建一个新的 HTTP 请求，模拟浏览器的预检请求。
		req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
		// 设置关键的预检请求头。
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", "POST")

		// 让路由器处理这个模拟请求。
		router.ServeHTTP(w, req)

		// --- 断言 ---
		// 1. 检查响应状态码是否为 204 No Content。
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
		}
		// 2. 检查 Access-Control-Allow-Origin 响应头是否正确设置。
		if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
			t.Errorf("Incorrect Access-Control-Allow-Origin header. Got %s", w.Header().Get("Access-Control-Allow-Origin"))
		}
		// 3. 检查 Access-Control-Allow-Credentials 响应头是否为 true。
		if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
			t.Errorf("Incorrect Access-Control-Allow-Credentials header. Got %s", w.Header().Get("Access-Control-Allow-Credentials"))
		}
	})

	// --- 测试用例 2: 实际请求 (GET) ---
	t.Run("ActualRequest", func(t *testing.T) {
		// 创建一个新的路由器和中间件实例。
		router := gin.New()
		router.Use(middlewares.Cors())
		// 注册一个 GET 路由，用于测试实际的数据请求。
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"}) // 返回成功的 JSON 响应。
		})

		// 创建响应记录器和模拟请求。
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		// 实际请求只需要设置 Origin 头。
		req.Header.Set("Origin", "http://example.com")

		// 处理请求。
		router.ServeHTTP(w, req)

		// --- 断言 ---
		// 1. 检查状态码是否为 200 OK，这是由我们的路由处理器设置的。
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}
		// 2. 检查响应体是否正确。
		if w.Body.String() != `{"message":"success"}` {
			t.Errorf("Incorrect response body. Got %s", w.Body.String())
		}
		// 3. 检查 Access-Control-Allow-Origin 头是否与请求的 Origin 一致。
		if w.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
			t.Errorf("Incorrect Access-Control-Allow-Origin header. Got %s", w.Header().Get("Access-Control-Allow-Origin"))
		}
	})
}
