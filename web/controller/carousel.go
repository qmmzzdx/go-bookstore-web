package controller

import (
	"bookstore/model"
	"bookstore/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CarouselController 轮播图控制器
// 负责处理与轮播图相关的HTTP请求，包括获取激活的轮播图、获取轮播图详情、创建、更新和删除轮播图
type CarouselController struct {
	CarouselService *service.CarouselService // 轮播图服务
}

// NewCarouselController 创建新的轮播图控制器实例
// 返回:
//
//	*CarouselController - 初始化好的轮播图控制器
func NewCarouselController() *CarouselController {
	return &CarouselController{
		CarouselService: service.NewCarouselService(),
	}
}

// GetActiveCarousels 获取激活的轮播图
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取激活轮播图请求，返回所有状态为激活的轮播图列表
func (c *CarouselController) GetActiveCarousels(ctx *gin.Context) {
	carousels, err := c.CarouselService.GetActiveCarousels()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "获取轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取轮播图成功",
		"data":    carousels,
	})
}

// GetCarouselByID 根据ID获取轮播图
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理根据ID获取轮播图请求，验证参数并返回指定轮播图的详细信息
func (c *CarouselController) GetCarouselByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "无效的轮播图ID",
		})
		return
	}

	carousel, err := c.CarouselService.GetCarouselByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    1,
			"message": "轮播图不存在",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取轮播图成功",
		"data":    carousel,
	})
}

// CreateCarousel 创建轮播图
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理创建轮播图请求，验证参数并创建新的轮播图
func (c *CarouselController) CreateCarousel(ctx *gin.Context) {
	var carousel model.Carousel
	if err := ctx.ShouldBindJSON(&carousel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	err := c.CarouselService.CreateCarousel(&carousel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "创建轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建轮播图成功",
		"data":    carousel,
	})
}

// UpdateCarousel 更新轮播图
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理更新轮播图请求，验证参数并更新指定轮播图的信息
func (c *CarouselController) UpdateCarousel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "无效的轮播图ID",
		})
		return
	}

	var carousel model.Carousel
	if err := ctx.ShouldBindJSON(&carousel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	carousel.ID = id
	err = c.CarouselService.UpdateCarousel(&carousel)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "更新轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新轮播图成功",
		"data":    carousel,
	})
}

// DeleteCarousel 删除轮播图
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理删除轮播图请求，验证参数并删除指定轮播图
func (c *CarouselController) DeleteCarousel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "无效的轮播图ID",
		})
		return
	}

	err = c.CarouselService.DeleteCarousel(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "删除轮播图失败",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除轮播图成功",
	})
}
