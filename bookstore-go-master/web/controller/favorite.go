package controller

import (
	"bookstore/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteController 收藏控制器
// 负责处理与用户收藏相关的HTTP请求，包括添加收藏、取消收藏、检查收藏状态、获取收藏列表和收藏数量
type FavoriteController struct {
	favoriteService *service.FavoriteService // 收藏服务
}

// NewFavoriteController 创建新的收藏控制器实例
// 返回:
//
//	*FavoriteController - 初始化好的收藏控制器
func NewFavoriteController() *FavoriteController {
	return &FavoriteController{
		favoriteService: service.NewFavoriteService(),
	}
}

// getUserID 获取用户ID的辅助函数
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 返回:
//
//	int - 用户ID，如果未获取到则返回0
func getUserID(c *gin.Context) int {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(int)
}

// AddFavorite 添加收藏
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理添加收藏请求，验证用户登录状态和参数，并调用服务层添加收藏
func (f *FavoriteController) AddFavorite(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    1,
			"message": "请先登录",
		})
		return
	}

	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "无效的书籍ID",
		})
		return
	}

	err = f.favoriteService.AddFavorite(userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "添加收藏失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "添加收藏成功",
	})
}

// RemoveFavorite 取消收藏
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理取消收藏请求，验证用户登录状态和参数，并调用服务层取消收藏
func (f *FavoriteController) RemoveFavorite(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    1,
			"message": "请先登录",
		})
		return
	}

	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "无效的书籍ID",
		})
		return
	}

	err = f.favoriteService.RemoveFavorite(userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "取消收藏失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消收藏成功",
	})
}

// CheckFavorite 检查是否已收藏
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理检查收藏状态请求，验证用户登录状态和参数，并返回指定书籍的收藏状态
func (f *FavoriteController) CheckFavorite(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    1,
			"message": "请先登录",
		})
		return
	}

	bookIDStr := c.Param("id")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "无效的书籍ID",
		})
		return
	}

	isFavorited, err := f.favoriteService.IsFavorited(userID, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "检查收藏状态失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"is_favorited": isFavorited,
		},
	})
}

// GetUserFavorites 获取用户收藏列表
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取用户收藏列表请求，验证用户登录状态，支持分页和时间筛选，返回用户的收藏列表
func (f *FavoriteController) GetUserFavorites(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    1,
			"message": "请先登录",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))
	timeFilter := c.DefaultQuery("time_filter", "all")

	favorites, total, err := f.favoriteService.GetUserFavorites(userID, page, pageSize, timeFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "获取收藏列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"favorites":    favorites,
			"total":        total,
			"total_pages":  (int(total) + pageSize - 1) / pageSize,
			"current_page": page,
		},
	})
}

// GetUserFavoriteCount 获取用户收藏数量
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取用户收藏数量请求，验证用户登录状态，返回用户的收藏总数
func (f *FavoriteController) GetUserFavoriteCount(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    1,
			"message": "请先登录",
		})
		return
	}

	count, err := f.favoriteService.GetUserFavoriteCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "获取收藏数量失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"count": count,
		},
	})
}
