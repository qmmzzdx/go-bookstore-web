package router

import (
	"bookstore/web/controller"
	"bookstore/web/middleware"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化应用路由
// 返回配置好的gin.Engine实例，包含所有API路由定义
func InitRouter() *gin.Engine {
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

	// ========== 初始化依赖组件 ========== //

	// 创建Controller控制器实例
	userController := controller.NewUserController()         // 用户控制器
	captchaController := controller.NewCaptchaController()   // 验证码控制器
	bookController := controller.NewBookController()         // 书籍控制器
	categoryController := controller.NewCategoryController() // 分类控制器
	orderController := controller.NewOrderController()       // 订单控制器
	favoriteController := controller.NewFavoriteController() // 收藏控制器（注入服务）
	carouselController := controller.NewCarouselController() // 轮播图控制器（注入服务）

	// ========== 路由注册 ========== //

	// API v1 路由组（所有v1版本API前缀为/api/v1）
	v1 := r.Group("/api/v1")
	{
		// ----- 用户相关路由 ----- //
		user := v1.Group("/user")
		{
			// 无需认证的公共接口
			user.POST("/register", userController.Register) // 用户注册
			user.POST("/login", userController.Login)       // 用户登录

			// 需要JWT认证的私有接口
			auth := user.Group("")
			auth.Use(middleware.JWTAuthMiddleware()) // 添加JWT认证中间件
			{
				auth.GET("/profile", userController.GetUserProfile)    // 获取用户资料
				auth.PUT("/profile", userController.UpdateUserProfile) // 更新用户资料
				auth.PUT("/password", userController.ChangePassword)   // 修改密码
				auth.DELETE("/logout", userController.Logout)          // 用户登出
			}
		}

		// ----- 书籍相关路由 ----- //
		book := v1.Group("/book")
		{
			book.GET("/hot", bookController.GetHotBooks)                       // 获取热门书籍
			book.GET("/new", bookController.GetNewBooks)                       // 获取新书上架
			book.GET("/list", bookController.GetBookList)                      // 获取书籍列表
			book.GET("/search", bookController.SearchBooks)                    // 搜索书籍
			book.GET("/detail/:id", bookController.GetBookDetail)              // 获取书籍详情
			book.GET("/category/:category", bookController.GetBooksByCategory) // 按分类获取书籍
		}

		// ----- 分类相关路由 ----- //
		category := v1.Group("/category")
		{
			category.GET("/list", categoryController.GetCategories)     // 获取全部分类
			category.GET("/:id", categoryController.GetCategoryByID)    // 获取分类详情
			category.POST("/create", categoryController.CreateCategory) // 创建分类
			category.PUT("/:id", categoryController.UpdateCategory)     // 更新分类
			category.DELETE("/:id", categoryController.DeleteCategory)  // 删除分类
		}

		// ----- 订单相关路由 ----- //
		order := v1.Group("/order")
		order.Use(middleware.JWTAuthMiddleware()) // 需要登录
		{
			order.POST("/create", orderController.CreateOrder)           // 创建订单
			order.GET("/:id", orderController.GetOrderByID)              // 获取订单详情
			order.GET("/list", orderController.GetUserOrders)            // 获取用户订单列表
			order.POST("/:id/pay", orderController.PayOrder)             // 支付订单
			order.GET("/statistics", orderController.GetOrderStatistics) // 订单统计
		}

		// ----- 收藏相关路由 ----- //
		favorite := v1.Group("/favorite")
		favorite.Use(middleware.JWTAuthMiddleware()) // 需要登录
		{
			favorite.POST("/:id", favoriteController.AddFavorite)           // 添加收藏
			favorite.DELETE("/:id", favoriteController.RemoveFavorite)      // 移除收藏
			favorite.GET("/list", favoriteController.GetUserFavorites)      // 获取用户收藏列表
			favorite.GET("/:id/check", favoriteController.CheckFavorite)    // 检查是否收藏
			favorite.GET("/count", favoriteController.GetUserFavoriteCount) // 获取用户收藏数量
		}

		// ----- 验证码相关路由 ----- //
		captcha := v1.Group("/captcha")
		{
			captcha.GET("/generate", captchaController.GenerateCaptcha) // 生成验证码
			captcha.POST("/verify", captchaController.VerifyCaptcha)    // 验证验证码
		}

		// ----- 轮播图相关路由 ----- //
		carousel := v1.Group("/carousel")
		{
			carousel.GET("/list", carouselController.GetActiveCarousels) // 获取有效轮播图
			carousel.GET("/:id", carouselController.GetCarouselByID)     // 获取轮播图详情
			carousel.POST("/create", carouselController.CreateCarousel)  // 创建轮播图
			carousel.PUT("/:id", carouselController.UpdateCarousel)      // 更新轮播图
			carousel.DELETE("/:id", carouselController.DeleteCarousel)   // 删除轮播图
		}
	}
	return r
}
