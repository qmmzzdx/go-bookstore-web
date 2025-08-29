package controller

import (
	"bookstore/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderController 订单控制器
// 负责处理与订单相关的HTTP请求，包括创建订单、获取订单详情、获取用户订单列表、支付订单和获取订单统计信息
type OrderController struct {
	OrderService *service.OrderService // 订单服务
}

// NewOrderController 创建新的订单控制器实例
// 返回:
//
//	*OrderController - 初始化好的订单控制器
func NewOrderController() *OrderController {
	return &OrderController{
		OrderService: service.NewOrderService(),
	}
}

// CreateOrder 创建订单
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理创建订单请求，验证参数并从上下文中获取用户ID，然后调用服务层创建订单
func (o *OrderController) CreateOrder(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户未登录",
		})
		return
	}
	req.UserID = userID.(int)

	order, err := o.OrderService.CreateOrder(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "创建订单失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    order,
		"message": "创建订单成功",
	})
}

// GetOrderByID 根据ID获取订单
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理根据ID获取订单请求，验证参数并调用服务层获取订单详情
func (o *OrderController) GetOrderByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的订单ID",
		})
		return
	}

	order, err := o.OrderService.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "订单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    order,
		"message": "获取订单成功",
	})
}

// GetUserOrders 获取用户的订单列表
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取用户订单列表请求，从上下文中获取用户ID，支持分页参数，返回用户的订单列表
func (o *OrderController) GetUserOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户未登录",
		})
		return
	}

	orders, total, err := o.OrderService.GetUserOrders(userID.(int), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取订单列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"orders":      orders,
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
		"message": "获取订单列表成功",
	})
}

// PayOrder 支付订单
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理支付订单请求，验证参数并调用服务层完成订单支付
func (o *OrderController) PayOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "无效的订单ID",
		})
		return
	}

	err = o.OrderService.PayOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "支付失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "支付成功",
	})
}

// GetOrderStatistics 获取订单统计信息
// 参数:
//
//	c - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取订单统计信息请求，从上下文中获取用户ID，返回用户的订单统计信息
func (o *OrderController) GetOrderStatistics(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    -1,
			"message": "用户未登录",
		})
		return
	}

	stats, err := o.OrderService.GetOrderStatistics(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": "获取统计信息失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    stats,
		"message": "获取统计信息成功",
	})
}
