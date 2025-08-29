package controller

import (
	"bookstore/jwt"
	"bookstore/model"
	"bookstore/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器，处理所有用户相关的HTTP请求
type UserController struct {
	UserService *service.UserService // 用户服务层实例
}

// NewUserController 创建用户控制器实例
// 返回:
//
//	*UserController - 初始化好的用户控制器
func NewUserController() *UserController {
	return &UserController{
		UserService: service.NewUserService(), // 初始化用户服务
	}
}

// RegisterRequest 注册请求数据结构
type RegisterRequest struct {
	Username        string `json:"username" binding:"required"`         // 用户名（必填）
	Password        string `json:"password" binding:"required"`         // 密码（必填）
	ConfirmPassword string `json:"confirm_password" binding:"required"` // 确认密码（必填）
	Email           string `json:"email" binding:"required,email"`      // 邮箱（必填且需符合email格式）
	Phone           string `json:"phone" binding:"required"`            // 手机号（必填）
	CaptchaID       string `json:"captcha_id" binding:"required"`       // 验证码ID（必填）
	CaptchaValue    string `json:"captcha_value" binding:"required"`    // 验证码值（必填）
}

// LoginRequest 登录请求数据结构
type LoginRequest struct {
	Username     string `json:"username" binding:"required"`      // 用户名（必填）
	Password     string `json:"password" binding:"required"`      // 密码（必填）
	CaptchaID    string `json:"captcha_id" binding:"required"`    // 验证码ID（必填）
	CaptchaValue string `json:"captcha_value" binding:"required"` // 验证码值（必填）
}

// Register 处理用户注册请求
// 路由: POST /user/register
func (u *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证验证码有效性
	captchaService := service.NewCaptchaService()
	if !captchaService.VerifyCaptcha(req.CaptchaID, req.CaptchaValue) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "验证码错误",
		})
		return
	}

	// 检查两次密码是否一致
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "两次密码不一致",
		})
		return
	}

	// 调用服务层注册用户
	err := u.UserService.UserRegister(req.Username, req.Password, req.Phone, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "注册成功",
	})
}

// Login 处理用户登录请求
// 路由: POST /user/login
func (u *UserController) Login(c *gin.Context) {
	var req LoginRequest
	// 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证验证码有效性
	captchaService := service.NewCaptchaService()
	if !captchaService.VerifyCaptcha(req.CaptchaID, req.CaptchaValue) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "验证码错误",
		})
		return
	}

	// 调用服务层进行登录验证
	loginResponse, err := u.UserService.UserLogin(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	// 返回成功响应（包含token等登录信息）
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    loginResponse,
		"message": "登录成功",
	})
}

// GetUserProfile 获取用户个人信息
// 路由: GET /user/profile (需认证)
func (u *UserController) GetUserProfile(c *gin.Context) {
	// 从JWT中间件设置的上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户未登录",
		})
		return
	}

	// 调用服务层获取用户信息
	user, err := u.UserService.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	// 构建响应数据（过滤敏感信息）
	response := gin.H{
		"id":         uint(user.ID),
		"username":   user.Username,
		"email":      user.Email,
		"phone":      user.Phone,
		"avatar":     user.Avatar,
		"created_at": user.CreatedAt.Format("1970-01-01 00:00:00"), // 格式化时间
		"updated_at": user.UpdatedAt.Format("1970-01-01 00:00:00"),
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    response,
		"message": "获取用户信息成功",
	})
}

// Logout 处理用户退出登录
// 路由: DELETE /user/logout (需认证)
func (u *UserController) Logout(c *gin.Context) {
	// 从JWT中间件设置的上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户未登录",
		})
		return
	}

	// 调用JWT服务撤销token
	err := jwt.RevokeToken(uint(userID.(int)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "退出登录失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "退出登录成功",
	})
}

// UpdateUserProfile 更新用户个人信息
// 路由: PUT /user/profile (需认证)
func (u *UserController) UpdateUserProfile(c *gin.Context) {
	// 从JWT中间件设置的上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户未登录",
		})
		return
	}

	// 定义更新数据结构
	var updateData struct {
		Username string `json:"username"` // 新用户名
		Email    string `json:"email"`    // 新邮箱
		Phone    string `json:"phone"`    // 新手机号
		Avatar   string `json:"avatar"`   // 新头像URL
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 构建用户对象
	user := &model.User{
		ID:       userID.(int),
		Username: updateData.Username,
		Email:    updateData.Email,
		Phone:    updateData.Phone,
		Avatar:   updateData.Avatar,
	}

	// 调用服务层更新数据
	err := u.UserService.UpdateUserInfo(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	// 获取更新后的用户信息
	updatedUser, err := u.UserService.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取用户信息失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "用户信息更新成功",
		"data":    updatedUser,
	})
}

// ChangePassword 修改用户密码
// 路由: PUT /user/password (需认证)
func (u *UserController) ChangePassword(c *gin.Context) {
	// 从JWT中间件设置的上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户未登录",
		})
		return
	}

	// 定义密码修改数据结构
	var passwordData struct {
		OldPassword string `json:"old_password" binding:"required"` // 旧密码（必填）
		NewPassword string `json:"new_password" binding:"required"` // 新密码（必填）
	}

	// 绑定请求数据
	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证新密码长度
	if len(passwordData.NewPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "新密码至少6位",
		})
		return
	}

	// 调用服务层修改密码
	err := u.UserService.ChangePassword(userID.(int), passwordData.OldPassword, passwordData.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "密码修改成功",
	})
}
