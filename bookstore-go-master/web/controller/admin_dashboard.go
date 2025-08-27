package controller

import (
	"bookstore/global"
	"bookstore/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminDashboardController 管理员仪表盘控制器
// 负责提供管理员仪表盘相关的数据统计和展示功能
type AdminDashboardController struct{}

// NewAdminDashboardController 创建新的管理员仪表盘控制器实例
// 返回:
//
//	*AdminDashboardController - 初始化好的管理员仪表盘控制器
func NewAdminDashboardController() *AdminDashboardController {
	return &AdminDashboardController{}
}

// DashboardStats 仪表盘统计数据
// 包含系统核心指标的统计信息，用于展示在管理员仪表盘上
type DashboardStats struct {
	TotalBooks   int64   `json:"total_books"`   // 图书总数
	TotalOrders  int64   `json:"total_orders"`  // 订单总数
	TotalUsers   int64   `json:"total_users"`   // 用户总数
	TotalRevenue float64 `json:"total_revenue"` // 总收入
	RecentBooks  []Book  `json:"recent_books"`  // 最近添加的图书列表
}

// Book 简化的图书信息
// 用于展示在仪表盘上的图书基本信息
type Book struct {
	ID        int    `json:"id"`         // 图书ID
	Title     string `json:"title"`      // 图书标题
	Author    string `json:"author"`     // 图书作者
	Price     int    `json:"price"`      // 图书价格
	CoverURL  string `json:"cover_url"`  // 图书封面URL
	CreatedAt string `json:"created_at"` // 创建时间
}

// GetDashboardStats 获取仪表盘统计数据
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 返回系统核心指标的统计数据，包括图书总数、订单总数、用户总数、总收入和最近添加的图书
func (c *AdminDashboardController) GetDashboardStats(ctx *gin.Context) {
	var stats DashboardStats

	// 获取图书总数
	global.DBClient.Model(&model.Book{}).Count(&stats.TotalBooks)

	// 获取订单总数
	global.DBClient.Model(&model.Order{}).Count(&stats.TotalOrders)

	// 获取用户总数
	global.DBClient.Model(&model.User{}).Count(&stats.TotalUsers)

	// 获取总收入（已支付的订单）
	var totalRevenue int64
	global.DBClient.Model(&model.Order{}).
		Where("is_paid = ?", true).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalRevenue)
	stats.TotalRevenue = float64(totalRevenue) // 已经是元为单位

	// 获取最近3天添加的图书
	var recentBooks []model.Book
	threeDaysAgo := time.Now().AddDate(0, 0, -3)
	global.DBClient.Where("created_at >= ?", threeDaysAgo).
		Order("created_at DESC").
		Limit(10).
		Find(&recentBooks)

	// 转换为简化的图书信息
	for _, book := range recentBooks {
		stats.RecentBooks = append(stats.RecentBooks, Book{
			ID:        book.ID,
			Title:     book.Title,
			Author:    book.Author,
			Price:     book.Price,
			CoverURL:  book.CoverURL,
			CreatedAt: book.CreatedAt.Format("1970-01-01 00:00:00"),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取统计数据成功",
		"data":    stats,
	})
}
