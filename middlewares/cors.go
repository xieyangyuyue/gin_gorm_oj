package middlewares

import (
	"github.com/gin-gonic/gin" // 导入 Gin 框架，用于构建 Web 应用和中间件。
	"net/http"                 // 导入 net/http 包，包含了 HTTP 客户端和服务器的实现，这里主要用其状态码常量。
)

// Cors 函数返回一个配置好的 Gin 中间件，用于处理跨域资源共享 (CORS) 请求。
// 这种工厂模式允许我们在未来轻松地为中间件添加配置参数。
func Cors() gin.HandlerFunc {
	// 返回一个闭包，这才是真正的中间件处理函数。
	return func(c *gin.Context) {
		// 从请求头中获取 "Origin" 字段，它标识了发起请求的源（协议、域名、端口）。
		origin := c.Request.Header.Get("Origin")
		// 从请求头中获取请求方法，例如 GET, POST, OPTIONS 等。
		method := c.Request.Method

		// 设置 Access-Control-Allow-Origin 响应头。
		// 这是 CORS 的核心。我们不再使用 `*`，而是直接将请求的 Origin 返回。
		// 这样可以支持携带凭证的请求。在生产环境中，为了安全，应该添加一个白名单校验。
		// 例如: if origin in whitelist { c.Header("Access-Control-Allow-Origin", origin) }
		c.Header("Access-Control-Allow-Origin", origin)

		// 设置 Access-Control-Allow-Headers 响应头。
		// 表明服务器允许在实际请求中携带哪些自定义请求头。
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, X-Requested-With, Accept")

		// 设置 Access-Control-Allow-Methods 响应头。
		// 表明服务器允许的跨域请求方法。
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")

		// 设置 Access-Control-Expose-Headers 响应头。
		// 表明允许浏览器（客户端）访问的响应头，默认只能访问一些基本头部。
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")

		// 设置 Access-Control-Allow-Credentials 响应头。
		// 值为 "true" 表示允许浏览器发送包含凭证（如 Cookies）的请求。
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求 (Preflight Request)。
		// 当请求方法为 OPTIONS 时，这是一个预检请求，浏览器用它来询问服务器支持哪些方法和头部。
		// 我们应直接返回 204 No Content 状态码，并终止后续处理。
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 使用 204 状态码响应，表示请求成功但无返回内容。
			return                                  // 终止中间件链的执行。
		}

		// 如果不是预检请求，则继续处理后续的中间件或路由处理器。
		c.Next()
	}
}

//func Cors() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		method := c.Request.Method
//		c.Header("Access-Control-Allow-Origin", "*")
//		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
//		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, AccessToken")
//		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
//		c.Header("Access-Control-Allow-Credentials", "true")
//		if method == "OPTIONS" {
//			c.AbortWithStatus(http.StatusNoContent)
//		}
//		c.Next()
//	}
//}
