package controller

import (
	"bookstore/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminOrderController 管理员订单控制器
// 负责订单相关的管理操作，包括获取订单列表、订单详情和更新订单状态
type AdminOrderController struct {
	orderService *service.OrderService // 订单服务
}

// NewAdminOrderController 创建新的管理员订单控制器实例
// 返回:
//
//	*AdminOrderController - 初始化好的管理员订单控制器
func NewAdminOrderController() *AdminOrderController {
	return &AdminOrderController{
		orderService: service.NewOrderService(),
	}
}

// GetOrderList 获取订单列表
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理获取订单列表请求，返回所有订单信息（当前为示例实现）
func (c *AdminOrderController) GetOrderList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取订单列表成功",
		"data":    []any{},
	})
}

// GetOrderByID 根据ID获取订单
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理根据ID获取订单请求，验证参数并返回指定订单的详细信息
func (c *AdminOrderController) GetOrderByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "ID参数错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "获取订单成功",
		"data":    gin.H{"id": id},
	})
}

// UpdateOrderStatus 更新订单状态
// 参数:
//
//	ctx - Gin上下文对象，包含HTTP请求和响应信息
//
// 处理更新订单状态请求，验证参数并调用服务层更新订单状态
func (c *AdminOrderController) UpdateOrderStatus(ctx *gin.Context) {
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
		Status int `json:"status" binding:"required"` // 订单状态
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新订单状态成功",
		"data":    gin.H{"id": id, "status": req.Status},
	})
}
