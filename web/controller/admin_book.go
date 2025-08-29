package controller

import (
	"bookstore/model"
	"bookstore/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminBookController 管理员图书控制器
// 负责图书和分类的增删改查等管理操作
type AdminBookController struct {
	bookService *service.BookService // 图书服务
}

// NewAdminBookController 创建新的管理员图书控制器实例
// 返回:
//
//	*AdminBookController - 初始化好的管理员图书控制器
func NewAdminBookController() *AdminBookController {
	return &AdminBookController{
		bookService: service.NewBookService(),
	}
}

// CreateBook 创建图书
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理创建图书请求，验证参数并调用服务层创建图书
func (c *AdminBookController) CreateBook(ctx *gin.Context) {
	var req model.BookCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := c.bookService.CreateBookFromRequest(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "创建图书失败: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建图书成功",
	})
}

// UpdateBook 更新图书
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理更新图书请求，验证参数并调用服务层更新图书信息
func (c *AdminBookController) UpdateBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	var req model.BookUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := c.bookService.UpdateBookFromRequest(uint(id), &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "更新图书失败: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新图书成功",
	})
}

// DeleteBook 删除图书
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理删除图书请求，验证参数并调用服务层删除图书
func (c *AdminBookController) DeleteBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	if err := c.bookService.DeleteBook(int(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "删除图书失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除图书成功",
	})
}

// GetBookByID 根据ID获取图书
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理根据ID获取图书请求，验证参数并调用服务层获取图书详情
func (c *AdminBookController) GetBookByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	book, err := c.bookService.GetBookByIDForAdmin(int(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "图书不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取图书成功",
		"data":    book,
	})
}

// GetBookList 获取图书列表
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取图书列表请求，支持分页和筛选条件
func (c *AdminBookController) GetBookList(ctx *gin.Context) {
	var req model.BookListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := c.bookService.GetBookList(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取图书列表失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取图书列表成功",
		"data":    result,
	})
}

// UpdateBookStatus 更新图书状态
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理更新图书状态请求，验证参数并调用服务层更新图书状态
func (c *AdminBookController) UpdateBookStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	// 从查询参数获取状态值，更简单直接
	statusStr := ctx.Query("status")
	if statusStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "缺少状态参数",
		})
		return
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil || (status != 0 && status != 1) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "状态参数错误，只能是0或1",
		})
		return
	}

	if err := c.bookService.UpdateBookStatus(uint(id), status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "更新图书状态失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新图书状态成功",
	})
}

// GetCategories 获取所有分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取所有分类请求，调用服务层获取分类列表
func (c *AdminBookController) GetCategories(ctx *gin.Context) {
	categories, err := c.bookService.GetCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取分类成功",
		"data":    categories,
	})
}

// CreateCategory 创建分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理创建分类请求，验证参数并调用服务层创建分类
func (c *AdminBookController) CreateCategory(ctx *gin.Context) {
	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := c.bookService.CreateCategory(&category); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "创建分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建分类成功",
	})
}

// UpdateCategory 更新分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理更新分类请求，验证参数并调用服务层更新分类信息
func (c *AdminBookController) UpdateCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	var updates map[string]any
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if err := c.bookService.UpdateCategory(uint(id), updates); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "更新分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新分类成功",
	})
}

// DeleteCategory 删除分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理删除分类请求，验证参数并调用服务层删除分类
func (c *AdminBookController) DeleteCategory(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	if err := c.bookService.DeleteCategory(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "删除分类失败: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除分类成功",
	})
}
