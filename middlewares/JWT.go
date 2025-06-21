package middlewares // Package middlewares 声明包名为 middlewares，通常用于存放 Gin 框架的中间件函数

import (
	"errors"
	"fmt"                          // 导入 fmt 包，用于格式化输出错误信息
	"gin_gorm_oj/utils"            // 导入自定义的 utils 包，用于获取 JWT 密钥。请确保该路径正确
	"github.com/golang-jwt/jwt/v5" // 导入 go-jwt/jwt v5 库，用于处理 JSON Web Tokens (JWT)
	"time"                         // 导入 time 包，用于处理时间相关的操作，如设置 token 过期时间
)

// UserClaims 结构体定义了 JWT Payload 中的自定义声明（claims）。
// 它包含了用户的身份信息，以及 JWT 标准注册声明。
type UserClaims struct {
	Identity             string `json:"identity"` // 用户的唯一标识符，例如用户ID或用户名
	Name                 string `json:"name"`     // 用户的名称
	IsAdmin              int    `json:"is_admin"` // 用户是否是管理员，0表示否，1表示是
	jwt.RegisteredClaims        // 嵌入 jwt.RegisteredClaims，包含了 JWT 标准定义的声明，如 Issuer, Subject, Audience, ExpiresAt 等
}

// myKey 是用于签名和验证 JWT 的密钥。
// 在生产环境中，这个密钥应该是一个强随机字符串，并且从安全的环境变量、配置文件或密钥管理服务中加载，
// 避免硬编码在代码中。
var myKey = []byte(utils.JwtKey) // 从 utils 包中获取 JWT 密钥并转换为字节切片

// GenerateToken 函数用于根据用户身份信息生成一个签名的 JWT。
//
// 参数:
//
//	identity string: 用户的唯一标识。
//	name string: 用户的名称。
//	isAdmin int: 表示用户是否是管理员的标志。
//
// 返回值:
//
//	string: 生成的 JWT 字符串。
//	error: 如果生成过程中发生错误，则返回错误信息。
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	// 创建 UserClaims 实例，填充自定义声明和标准注册声明。
	UserClaim := &UserClaims{
		Identity: identity, // 设置用户唯一标识
		Name:     name,     // 设置用户名称
		IsAdmin:  isAdmin,  // 设置管理员标志
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gin-gorm-oj",                                                // Token 的发行者，推荐设置为应用名称
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(7 * 24 * time.Hour)), // 设置 Token 过期时间为当前时间起7天后（推荐使用 UTC 时间）
			//JTI: "unique_id", // 可以在这里添加 JWT ID (JTI)，用于防止重放攻击，确保Token的唯一性
			NotBefore: jwt.NewNumericDate(time.Now().UTC()), // Token 生效时间，可在此时间之前拒绝 Token
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()), // Token 的发行时间
		},
	}
	// 使用 HMAC SHA256 签名方法和自定义声明创建新的 JWT Token 对象。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	// 使用预定义的密钥 myKey 对 Token 进行签名，并返回签名后的字符串。
	return token.SignedString(myKey)
}

// AnalyseToken 函数用于解析和验证给定的 JWT 字符串。
//
// 参数:
//
//	tokenString string: 待解析和验证的 JWT 字符串。
//
// 返回值:
//
//	*UserClaims: 如果解析和验证成功，返回包含用户声明的 UserClaims 指针。
//	error: 如果解析或验证失败，则返回相应的错误信息。
func AnalyseToken(tokenString string) (*UserClaims, error) {
	// 使用 jwt.ParseWithClaims 解析 Token 字符串。
	// 它需要 Token 字符串、用于存储解析后声明的结构体实例（&UserClaims{}）
	// 以及一个用于提供签名密钥的回调函数。
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{}, // 用于将 Token 中的 claims 解析到 UserClaims 结构体中
		func(token *jwt.Token) (interface{}, error) {
			// 在此回调函数中返回用于验证签名的密钥。
			// 通常需要检查 token.Method 是否是你预期的签名方法（例如 jwt.SigningMethodHS256）。
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return myKey, nil // 返回签名密钥
		},
	)
	//log.Printf("格式化日志: %s", token)
	//log.Printf("格式化日志: %s", myKey)
	// 检查解析过程中是否发生错误。
	if err != nil {
		// 区分不同类型的解析错误，提供更详细的错误信息
		// 在 jwt v5 中，可以直接比较错误类型或者使用 errors.Is 来检查
		if errors.Is(err, jwt.ErrTokenMalformed) { // Token 格式不正确
			return nil, fmt.Errorf("token is malformed: %w", err) // Token 格式不正确
		} else if errors.Is(err, jwt.ErrTokenExpired) { // Token 已过期
			return nil, fmt.Errorf("token is expired: %w", err) // Token 已过期
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) { // Token 尚未生效
			return nil, fmt.Errorf("token not active yet: %w", err) // Token 尚未生效
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) { // Token 签名无效
			return nil, fmt.Errorf("token signature is invalid: %w", err) // Token 签名无效
		} else {
			// 对于其他未明确处理的错误，可能是通用的解析失败或其他内部错误
			return nil, fmt.Errorf("failed to parse token: %w", err) // 通用解析失败
		}
	}

	// 检查解析后的 Token 是否有效，并且 Claims 是否可以断言为 *UserClaims 类型。
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		// 添加自定义字段验证
		if claims.Identity == "" { // 假设UserID是必要字段
			return nil, fmt.Errorf("invalid token or claims type mismatch")
		}
		return claims, nil // Token 有效且声明类型正确，返回 claims
	}
	// 如果 Token 无效或 claims 类型不匹配，则返回错误。
	return nil, fmt.Errorf("invalid token or claims type mismatch")
}
