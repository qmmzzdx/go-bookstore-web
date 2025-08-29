package middleware

import (
	"bookstore/global"
	"bookstore/jwt"
	"bookstore/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminAuthMiddleware 管理员认证中间件
// 该中间件用于验证请求是否来自有效的管理员用户
// 会检查JWT令牌的有效性、用户权限等信息
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "未提供认证令牌",
			})
			c.Abort() // 终止请求处理链
			return
		}

		// 检查Authorization头是否以Bearer开头
		// JWT标准格式要求Authorization头格式为: Bearer <token>
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "认证令牌格式错误",
			})
			c.Abort()
			return
		}

		// 提取真正的token字符串（去掉"Bearer "前缀）
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析JWT令牌
		// 这里会验证令牌的签名、有效期等基本信息
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "认证令牌无效",
			})
			c.Abort()
			return
		}

		// 检查令牌类型是否为access token
		// 防止使用refresh token进行接口访问
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "token类型错误，请使用access token",
			})
			c.Abort()
			return
		}

		// 根据令牌中的用户ID查询数据库
		var user model.User
		if err := global.DBClient.First(&user, int(claims.UserID)).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    -1,
				"message": "用户不存在",
			})
			c.Abort()
			return
		}

		// 验证用户是否为管理员
		// IsAdmin字段为true表示是管理员用户
		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    -1,
				"message": "权限不足，需要管理员权限",
			})
			c.Abort()
			return
		}

		// 认证通过，将用户信息存入gin上下文
		// 后续处理函数可以通过c.Get("admin_user")获取用户信息
		c.Set("admin_user", user)       // 存储完整的用户对象
		c.Set("admin_user_id", user.ID) // 存储用户ID

		// 继续处理后续中间件或路由处理函数
		c.Next()
	}
}
