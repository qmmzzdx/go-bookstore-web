package controller

import (
	"bookstore/model"
	"bookstore/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryController 分类控制器
// 负责处理与图书分类相关的HTTP请求，包括获取分类列表、分类详情、创建、更新和删除分类
type CategoryController struct {
	CategoryService *service.CategoryService // 分类服务
}

// NewCategoryController 创建新的分类控制器实例
// 返回:
//
//	*CategoryController - 初始化好的分类控制器
func NewCategoryController() *CategoryController {
	return &CategoryController{
		CategoryService: service.NewCategoryService(),
	}
}

// GetCategories 获取所有分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取所有分类请求，返回系统中的全部分类列表
func (c *CategoryController) GetCategories(ctx *gin.Context) {
	categories, err := c.CategoryService.GetAllCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取分类列表失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    categories,
		"message": "获取分类列表成功",
	})
}

// GetCategoryByID 根据ID获取分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理根据ID获取分类请求，验证参数并返回指定分类的详细信息
func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的分类ID",
		})
		return
	}

	category, err := c.CategoryService.GetCategoryByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "分类不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    category,
		"message": "获取分类成功",
	})
}

// CreateCategory 创建分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理创建分类请求，验证参数并创建新的分类
func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	err := c.CategoryService.CreateCategory(&category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "创建分类失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    category,
		"message": "创建分类成功",
	})
}

// UpdateCategory 更新分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理更新分类请求，验证参数并更新指定分类的信息
func (c *CategoryController) UpdateCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的分类ID",
		})
		return
	}

	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	category.ID = id
	err = c.CategoryService.UpdateCategory(&category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "更新分类失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    category,
		"message": "更新分类成功",
	})
}

// DeleteCategory 删除分类
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理删除分类请求，验证参数并删除指定分类
func (c *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的分类ID",
		})
		return
	}

	err = c.CategoryService.DeleteCategory(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "删除分类失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除分类成功",
	})
}
