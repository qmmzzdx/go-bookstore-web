package router

import (
	"bookstore/web/controller"
	"bookstore/web/middleware"

	"github.com/gin-gonic/gin"
)

// InitAdminRouter 初始化管理员后台路由
// 返回配置好的gin.Engine实例，包含管理员后台所有API路由定义
func InitAdminRouter() *gin.Engine {
	// 创建默认gin实例（已包含Logger和Recovery中间件）
	r := gin.Default()

	// 添加CORS跨域中间件
	// 处理浏览器跨域请求和OPTIONS预检请求
	r.Use(func(c *gin.Context) {
		// 允许所有源访问（生产环境应改为具体域名）
		c.Header("Access-Control-Allow-Origin", "*")
		// 允许的HTTP方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// 暴露给客户端的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length")
		// 允许携带凭证（如cookies）
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理OPTIONS预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 204 No Content
			return
		}
		c.Next() // 继续处理请求
	})

	// ========== 管理员认证路由（无需认证） ========== //

	// 管理员认证路由组（前缀为/api/v1/admin/auth）
	auth := r.Group("/api/v1/admin/auth")
	{
		auth.POST("/login", controller.NewAdminAuthController().Login) // 管理员登录
	}

	// ========== 管理员API路由组（需要认证） ========== //

	// 管理员API路由组（前缀为/api/v1/admin）
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AdminAuthMiddleware()) // 添加管理员认证中间件
	{
		// ----- 仪表盘统计 ----- //
		admin.GET("/dashboard/stats", controller.NewAdminDashboardController().GetDashboardStats) // 获取仪表盘统计数据

		// ----- 图书管理 ----- //
		books := admin.Group("/books")
		{
			books.GET("/list", controller.NewAdminBookController().GetBookList)            // 获取图书列表
			books.GET("/:id", controller.NewAdminBookController().GetBookByID)             // 获取图书详情
			books.POST("/create", controller.NewAdminBookController().CreateBook)          // 创建图书
			books.PUT("/:id", controller.NewAdminBookController().UpdateBook)              // 更新图书信息
			books.DELETE("/:id", controller.NewAdminBookController().DeleteBook)           // 删除图书
			books.PUT("/:id/status", controller.NewAdminBookController().UpdateBookStatus) // 更新图书状态
		}

		// ----- 分类管理 ----- //
		categories := admin.Group("/categories")
		{
			categories.GET("/list", controller.NewAdminBookController().GetCategories)     // 获取分类列表
			categories.POST("/create", controller.NewAdminBookController().CreateCategory) // 创建分类
			categories.PUT("/:id", controller.NewAdminBookController().UpdateCategory)     // 更新分类信息
			categories.DELETE("/:id", controller.NewAdminBookController().DeleteCategory)  // 删除分类
		}

		// ----- 用户管理 ----- //
		users := admin.Group("/users")
		{
			users.GET("/list", controller.NewAdminUserController().GetUserList)            // 获取用户列表
			users.GET("/:id", controller.NewAdminUserController().GetUserByID)             // 获取用户详情
			users.POST("/create", controller.NewAdminUserController().CreateUser)          // 创建用户
			users.PUT("/:id", controller.NewAdminUserController().UpdateUser)              // 更新用户信息
			users.DELETE("/:id", controller.NewAdminUserController().DeleteUser)           // 删除用户
			users.PUT("/:id/status", controller.NewAdminUserController().UpdateUserStatus) // 更新用户状态
		}
	}
	return r
}
