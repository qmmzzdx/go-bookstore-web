package middleware

import (
	"bookstore/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件（强制认证）
// 作用：验证请求中的JWT token，保护需要认证的接口
// 返回值：gin.HandlerFunc 中间件函数
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "请求头中缺少Authorization字段",
			})
			c.Abort() // 终止请求处理链
			return
		}

		// 检查Bearer token格式
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "Authorization格式错误，应为：Bearer {token}",
			})
			c.Abort()
			return
		}

		tokenString := tokenParts[1] // 提取实际的token字符串

		// 解析并验证JWT token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "无效的token",
				"error":   err.Error(), // 返回具体错误信息方便调试
			})
			c.Abort()
			return
		}

		// 检查token类型（只允许access token访问API）
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "token类型错误，请使用access token",
			})
			c.Abort()
			return
		}

		// 认证通过，将用户信息存储到gin上下文中
		c.Set("userID", int(claims.UserID)) // 设置用户ID
		c.Set("username", claims.Username)  // 设置用户名
		c.Set("authenticated", true)        // 标记已认证

		// 继续处理后续中间件或路由处理函数
		c.Next()
	}
}

// OptionalAuthMiddleware 可选认证中间件
// 作用：尝试认证但不强制，用于可选登录的接口（如获取用户信息接口）
// 返回值：gin.HandlerFunc 中间件函数
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 没有token，继续处理请求（视为未登录状态）
			c.Next()
			return
		}

		// 检查Bearer token格式
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			// token格式错误，但不中断请求（视为未登录状态）
			c.Next()
			return
		}

		tokenString := tokenParts[1]

		// 尝试解析token
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			// token无效，但不中断请求（视为未登录状态）
			c.Next()
			return
		}

		// 检查token类型（只处理access token）
		if claims.TokenType == "access" {
			// 认证成功，设置用户信息
			c.Set("userID", int(claims.UserID))
			c.Set("username", claims.Username)
			c.Set("authenticated", true) // 标记已认证
		}

		// 继续处理请求（可能是已登录或未登录状态）
		c.Next()
	}
}
