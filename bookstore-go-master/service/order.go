package service

import (
	"bookstore/global"
	"bookstore/model"
	"bookstore/repository"
	"errors"

	"gorm.io/gorm"
)

// OrderService 订单服务
// 负责订单相关的业务逻辑处理，包括创建订单、支付订单、查询订单等
type OrderService struct {
	OrderDAO *repository.OrderDAO // 订单数据访问对象
	BookDAO  *repository.BookDAO  // 图书数据访问对象
}

// CreateOrderRequest 创建订单请求
// 用于接收创建订单的请求数据
type CreateOrderRequest struct {
	UserID int                      `json:"user_id"` // 用户ID
	Items  []CreateOrderItemRequest `json:"items"`   // 订单项列表
}

// CreateOrderItemRequest 创建订单项请求
// 用于接收创建订单项的请求数据
type CreateOrderItemRequest struct {
	BookID   int `json:"book_id"`  // 图书ID
	Quantity int `json:"quantity"` // 购买数量
	Price    int `json:"price"`    // 价格（元）
}

// NewOrderService 创建新的订单服务实例
// 返回:
//
//	*OrderService - 初始化好的订单服务
func NewOrderService() *OrderService {
	return &OrderService{
		OrderDAO: repository.NewOrderDAO(),
		BookDAO:  repository.NewBookDAO(),
	}
}

// CreateOrder 创建订单
// 参数:
//
//	req - 创建订单请求对象指针
//
// 返回:
//
//	*model.Order - 创建的订单对象指针
//	error - 错误信息
func (o *OrderService) CreateOrder(req *CreateOrderRequest) (*model.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("订单项不能为空")
	}

	// 检查库存
	err := o.checkStockAvailability(req.Items)
	if err != nil {
		return nil, err
	}

	// 生成订单号
	orderNo := o.OrderDAO.GenerateOrderNo()

	// 计算总金额
	var totalAmount int
	var orderItems []*model.OrderItem

	for _, item := range req.Items {
		subtotal := item.Price * item.Quantity
		totalAmount += subtotal

		orderItems = append(orderItems, &model.OrderItem{
			BookID:   item.BookID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: subtotal,
		})
	}

	// 创建订单
	order := &model.Order{
		UserID:      req.UserID,
		OrderNo:     orderNo,
		TotalAmount: totalAmount,
		Status:      0, // 待支付
		IsPaid:      false,
	}

	// 创建订单和订单项
	err = o.OrderDAO.CreateOrderWithItems(order, orderItems)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// checkStockAvailability 检查库存是否充足
// 参数:
//
//	items - 订单项请求列表
//
// 返回:
//
//	error - 错误信息
func (o *OrderService) checkStockAvailability(items []CreateOrderItemRequest) error {
	for _, item := range items {
		book, err := o.BookDAO.GetBookByIDForAdmin(item.BookID)
		if err != nil {
			return errors.New("图书不存在")
		}

		if book.Status != 1 {
			return errors.New("图书已下架")
		}

		if book.Stock < item.Quantity {
			return errors.New("库存不足")
		}
	}
	return nil
}

// GetOrderByID 根据ID获取订单
// 参数:
//
//	id - 订单ID
//
// 返回:
//
//	*model.Order - 订单对象指针
//	error - 错误信息
func (o *OrderService) GetOrderByID(id int) (*model.Order, error) {
	return o.OrderDAO.GetOrderByID(id)
}

// GetOrderByOrderNo 根据订单号获取订单
// 参数:
//
//	orderNo - 订单号
//
// 返回:
//
//	*model.Order - 订单对象指针
//	error - 错误信息
func (o *OrderService) GetOrderByOrderNo(orderNo string) (*model.Order, error) {
	return o.OrderDAO.GetOrderByOrderNo(orderNo)
}

// GetUserOrders 获取用户的订单列表
// 参数:
//
//	userID - 用户ID
//	page - 页码
//	pageSize - 每页数量
//
// 返回:
//
//	[]*model.Order - 订单列表
//	int64 - 总记录数
//	error - 错误信息
func (o *OrderService) GetUserOrders(userID int, page, pageSize int) ([]*model.Order, int64, error) {
	return o.OrderDAO.GetUserOrders(userID, page, pageSize)
}

// PayOrder 支付订单
// 参数:
//
//	orderID - 订单ID
//
// 返回:
//
//	error - 错误信息
func (o *OrderService) PayOrder(orderID int) error {
	// 检查订单是否存在
	order, err := o.OrderDAO.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	// 检查订单是否已支付
	if order.IsPaid {
		return errors.New("订单已支付")
	}

	// 使用事务处理支付和库存更新
	err = global.DBClient.Transaction(func(tx *gorm.DB) error {
		// 再次检查库存（防止并发问题）
		for _, item := range order.OrderItems {
			var book model.Book
			if err := tx.First(&book, item.BookID).Error; err != nil {
				return errors.New("图书不存在")
			}
			if book.Stock < item.Quantity {
				return errors.New("库存不足")
			}
		}

		// 标记订单为已支付
		if err := tx.Model(&model.Order{}).Where("id = ?", orderID).Updates(
			map[string]any{
				"status":       1,
				"is_paid":      true,
				"payment_time": gorm.Expr("NOW()"),
			}).Error; err != nil {
			return err
		}

		// 更新图书库存和销售量
		for _, item := range order.OrderItems {
			if err := tx.Model(&model.Book{}).
				Where("id = ?", item.BookID).
				Updates(map[string]any{
					"stock": gorm.Expr("stock - ?", item.Quantity),
					"sale":  gorm.Expr("sale + ?", item.Quantity),
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// GetOrderStatistics 获取订单统计信息
// 参数:
//
//	userID - 用户ID
//
// 返回:
//
//	map[string]any - 订单统计信息
//	error - 错误信息
func (o *OrderService) GetOrderStatistics(userID int) (map[string]any, error) {
	return o.OrderDAO.GetOrderStatistics(userID)
}
