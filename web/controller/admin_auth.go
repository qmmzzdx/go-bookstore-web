package controller

import (
	"bookstore/global"
	"bookstore/jwt"
	"bookstore/model"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminAuthController 管理员认证控制器
// 负责管理员登录和获取管理员信息等认证相关的请求处理
type AdminAuthController struct{}

// NewAdminAuthController 创建新的管理员认证控制器实例
// 返回:
//
//	*AdminAuthController - 初始化好的管理员认证控制器
func NewAdminAuthController() *AdminAuthController {
	return &AdminAuthController{}
}

// AdminLoginRequest 管理员登录请求
// 包含管理员登录所需的用户名和密码字段
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名，必填
	Password string `json:"password" binding:"required"` // 密码，必填
}

// LoginResponse 管理员登录响应
// 包含登录成功后返回的令牌和用户信息
type LoginResponse struct {
	Token string     `json:"token"` // JWT令牌
	User  model.User `json:"user"`  // 用户信息
}

// Login 管理员登录
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理管理员登录请求，验证用户名密码，生成JWT令牌并返回
func (c *AdminAuthController) Login(ctx *gin.Context) {
	var req AdminLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 查询用户
	var user model.User
	if err := global.DBClient.Where("username = ?", req.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户名或密码错误",
		})
		return
	}

	// 验证密码
	// 将输入的密码进行base64编码，然后与数据库中的密码比较
	encodedInputPassword := base64.StdEncoding.EncodeToString([]byte(req.Password))
	if encodedInputPassword != user.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户名或密码错误",
		})
		return
	}

	// 检查是否为管理员
	if !user.IsAdmin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    -1,
			"message": "权限不足，需要管理员权限",
		})
		return
	}

	// 生成JWT token
	token, err := jwt.GenerateToken(uint(user.ID), user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "生成token失败",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": LoginResponse{
			Token: token,
			User:  user,
		},
	})
}

// GetProfile 获取管理员信息
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 从JWT令牌中解析出用户ID，查询并返回管理员信息
func (c *AdminAuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("admin_user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "未登录",
		})
		return
	}

	var user model.User
	if err := global.DBClient.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "用户不存在",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取用户信息成功",
		"data":    user,
	})
}
