package middlewares // Package middlewares 定义包名为 middlewares，表明该文件包含 Gin 框架的中间件功能。

import ( // 导入所需的 Go 语言标准库和第三方库。
	"github.com/gin-gonic/gin" // 导入 Gin Web 框架，用于构建 HTTP 服务和处理路由。
	"net/http"                 // 导入 net/http 包，提供了 HTTP 客户端和服务端的实现，用于处理 HTTP 状态码等。
)

// AuthAdminCheck is a middleware function that checks if the user is authenticated with admin role.
// AuthAdminCheck 是一个 Gin 中间件函数，用于检查用户是否以管理员角色进行认证。
func AuthAdminCheck() gin.HandlerFunc { // 定义 AuthAdminCheck 函数，返回一个 Gin 的 HandlerFunc 类型，即一个中间件函数。
	return func(c *gin.Context) { // 返回实际的中间件处理函数，该函数接收一个 Gin 上下文对象。
		auth := c.GetHeader("Authorization") // 从 HTTP 请求头中获取名为 "Authorization" 的值，通常包含 JWT token。
		userClaim, err := AnalyseToken(auth) // 调用 AnalyseToken 函数（假设此函数在其他地方定义），解析 JWT token 并返回用户声明信息和潜在错误。
		if err != nil {                      // 检查解析 token 是否发生错误。
			c.Abort()                    // 如果发生错误，终止当前请求的后续处理。
			c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
				"code": http.StatusUnauthorized,      // 设置响应状态码为 401 Unauthorized（未授权）。
				"msg":  "Unauthorized Authorization", // 设置错误信息为“未授权的 Authorization”。
			})
			return // 终止函数执行。
		}
		if userClaim == nil || userClaim.IsAdmin != 1 { // 检查 userClaim 是否为空，或者 userClaim 中的 IsAdmin 字段是否不等于 1（表示非管理员）。
			c.Abort()                    // 如果用户声明为空或者不是管理员，终止当前请求的后续处理。
			c.JSON(http.StatusOK, gin.H{ // 返回 JSON 格式的响应。
				"code": http.StatusUnauthorized, // 设置响应状态码为 401 Unauthorized（未授权）。
				"msg":  "Unauthorized Admin",    // 设置错误信息为“未授权的管理员”。
			})
			return // 终止函数执行。
		}
		c.Next() // 如果认证和管理员检查都通过，则继续执行请求链中的下一个处理函数（或下一个中间件）。
	}
}
