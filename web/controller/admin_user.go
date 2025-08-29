package controller

import (
	"bookstore/global"
	"bookstore/model"
	"bookstore/service"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminUserController 管理员用户控制器
// 负责用户相关的管理操作，包括用户列表查询、用户详情、创建用户、更新用户和删除用户等
type AdminUserController struct {
	userService *service.UserService // 用户服务
}

// NewAdminUserController 创建新的管理员用户控制器实例
// 返回:
//
//	*AdminUserController - 初始化好的管理员用户控制器
func NewAdminUserController() *AdminUserController {
	return &AdminUserController{
		userService: service.NewUserService(),
	}
}

// GetUserList 获取用户列表
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取用户列表请求，支持分页、用户名搜索、邮箱搜索和管理员状态筛选
func (c *AdminUserController) GetUserList(ctx *gin.Context) {
	var users []model.User
	var total int64

	// 获取查询参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	searchUsername := ctx.Query("username")
	searchEmail := ctx.Query("email")
	isAdmin := ctx.Query("is_admin")

	query := global.DBClient.Model(&model.User{})

	// 添加搜索条件
	if searchUsername != "" {
		query = query.Where("username LIKE ?", "%"+searchUsername+"%")
	}
	if searchEmail != "" {
		query = query.Where("email LIKE ?", "%"+searchEmail+"%")
	}
	if isAdmin != "" {
		isAdminBool := isAdmin == "true"
		query = query.Where("is_admin = ?", isAdminBool)
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取用户列表失败",
		})
		return
	}

	// 清除密码字段
	for i := range users {
		users[i].Password = ""
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取用户列表成功",
		"data": gin.H{
			"users":        users,
			"total":        total,
			"current_page": page,
			"page_size":    pageSize,
		},
	})
}

// GetUserByID 根据ID获取用户
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理根据ID获取用户请求，验证参数并返回指定用户的详细信息
func (c *AdminUserController) GetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	var user model.User
	if err := global.DBClient.First(&user, id).Error; err != nil {
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
		"message": "获取用户成功",
		"data":    user,
	})
}

// CreateUser 创建用户
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理创建用户请求，验证参数并创建新用户
func (c *AdminUserController) CreateUser(ctx *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`    // 用户名，必填
		Password string `json:"password" binding:"required"`    // 密码，必填
		Email    string `json:"email" binding:"required,email"` // 邮箱，必填且需符合邮箱格式
		Phone    string `json:"phone"`                          // 电话，可选
		IsAdmin  bool   `json:"is_admin"`                       // 是否为管理员，可选
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := global.DBClient.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "用户名已存在",
		})
		return
	}

	// 检查邮箱是否已存在
	if err := global.DBClient.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "邮箱已存在",
		})
		return
	}

	// Base64编码密码
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(req.Password))

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: encodedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		IsAdmin:  req.IsAdmin,
	}

	if err := global.DBClient.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "创建用户失败",
		})
		return
	}

	// 清除密码字段
	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建用户成功",
		"data":    user,
	})
}

// UpdateUser 更新用户
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理更新用户请求，验证参数并更新用户信息
func (c *AdminUserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	var req struct {
		Username string `json:"username"` // 用户名，可选
		Email    string `json:"email"`    // 邮箱，可选
		Phone    string `json:"phone"`    // 电话，可选
		IsAdmin  *bool  `json:"is_admin"` // 是否为管理员，可选
		Password string `json:"password"` // 密码，可选
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	var user model.User
	if err := global.DBClient.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "用户不存在",
		})
		return
	}

	// 更新字段
	updates := make(map[string]any)
	if req.Username != "" {
		// 检查用户名是否已被其他用户使用
		var existingUser model.User
		if err := global.DBClient.Where("username = ? AND id != ?", req.Username, id).First(&existingUser).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": "用户名已存在",
			})
			return
		}
		updates["username"] = req.Username
	}

	if req.Email != "" {
		// 检查邮箱是否已被其他用户使用
		var existingUser model.User
		if err := global.DBClient.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    -1,
				"message": "邮箱已存在",
			})
			return
		}
		updates["email"] = req.Email
	}

	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if req.IsAdmin != nil {
		updates["is_admin"] = *req.IsAdmin
	}

	if req.Password != "" {
		encodedPassword := base64.StdEncoding.EncodeToString([]byte(req.Password))
		updates["password"] = encodedPassword
	}

	if err := global.DBClient.Model(&user).Updates(updates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "更新用户失败",
		})
		return
	}

	// 重新获取用户信息
	global.DBClient.First(&user, id)
	user.Password = ""

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新用户成功",
		"data":    user,
	})
}

// DeleteUser 删除用户
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理删除用户请求，验证参数并删除指定用户
func (c *AdminUserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	var user model.User
	if err := global.DBClient.First(&user, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "用户不存在",
		})
		return
	}

	if err := global.DBClient.Delete(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "删除用户失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除用户成功",
	})
}

// UpdateUserStatus 更新用户状态
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理更新用户状态请求，验证参数并更新用户的管理员状态
func (c *AdminUserController) UpdateUserStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	var req struct {
		IsAdmin bool `json:"is_admin" binding:"required"` // 是否为管理员，必填
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := global.DBClient.Model(&model.User{}).Where("id = ?", id).Update("is_admin", req.IsAdmin).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "更新用户状态失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新用户状态成功",
	})
}
