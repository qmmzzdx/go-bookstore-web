package controller

import (
	"bookstore/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// BookController 书籍控制器
// 负责处理与书籍相关的HTTP请求，包括获取书籍列表、详情、热销书籍、新书、搜索和分类查询
type BookController struct {
	BookService *service.BookService // 书籍服务
}

// NewBookController 创建新的书籍控制器实例
// 返回:
//
//	*BookController - 初始化好的书籍控制器
func NewBookController() *BookController {
	return &BookController{
		BookService: service.NewBookService(),
	}
}

// GetBookList 获取书籍列表
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取书籍列表请求，支持分页参数，返回分页后的书籍列表
func (b *BookController) GetBookList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))

	books, total, err := b.BookService.GetBooksByPage(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取书籍列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"books":      books,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
		"message": "获取书籍列表成功",
	})
}

// GetBookDetail 获取书籍详情
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理根据ID获取书籍详情请求，验证参数并返回指定书籍的详细信息
func (b *BookController) GetBookDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的书籍ID",
		})
		return
	}

	book, err := b.BookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "书籍不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    book,
		"message": "获取书籍详情成功",
	})
}

// GetHotBooks 获取热销书籍
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取热销书籍请求，支持限制返回数量，返回销量最高的书籍列表
func (b *BookController) GetHotBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	books, err := b.BookService.GetHotBooks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取热销书籍失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    books,
		"message": "获取热销书籍成功",
	})
}

// GetNewBooks 获取新书
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取新书请求，支持限制返回数量，返回最新上架的书籍列表
func (b *BookController) GetNewBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))

	books, err := b.BookService.GetNewBooks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取新书失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    books,
		"message": "获取新书成功",
	})
}

// SearchBooks 搜索书籍
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理搜索书籍请求，支持关键词搜索和分页参数，返回匹配的书籍列表
func (b *BookController) SearchBooks(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		keyword = c.Query("keyword") // 兼容旧版本
	}

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "搜索关键词不能为空",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "12"))

	books, total, err := b.BookService.SearchBooksWithPagination(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "搜索书籍失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"books":      books,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
		"message": "搜索书籍成功",
	})
}

// GetBooksByCategory 根据分类获取书籍
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理根据分类获取书籍请求，返回指定分类下的所有书籍
func (b *BookController) GetBooksByCategory(c *gin.Context) {
	category := c.Param("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "分类不能为空",
		})
		return
	}

	books, err := b.BookService.GetBooksByType(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取分类书籍失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    books,
		"message": "获取分类书籍成功",
	})
}
