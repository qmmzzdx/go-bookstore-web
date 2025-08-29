package repository

import (
	"bookstore/global"
	"bookstore/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// OrderDAO 订单数据访问对象
// 封装了所有与订单相关的数据库操作
type OrderDAO struct {
	db *gorm.DB // GORM数据库连接实例
}

// NewOrderDAO 创建新的订单DAO实例
// 返回:
//
//	*OrderDAO - 初始化后的订单数据访问对象
func NewOrderDAO() *OrderDAO {
	return &OrderDAO{
		db: global.GetDB(), // 从全局变量获取数据库连接
	}
}

// CreateOrder 创建订单
// 参数:
//
//	order - 订单对象指针
//
// 返回:
//
//	error - 如果创建过程中出现错误则返回错误
func (o *OrderDAO) CreateOrder(order *model.Order) error {
	// 对应SQL: INSERT INTO orders (order_no, user_id, total_amount, status, is_paid, ...)
	// VALUES (order.OrderNo, order.UserID, order.TotalAmount, order.Status, order.IsPaid, ...);
	err := o.db.Create(order).Error
	return err
}

// GetOrderByID 根据ID获取订单
// 参数:
//
//	id - 订单ID
//
// 返回:
//
//	*model.Order - 订单对象指针（包含订单项和图书信息）
//	error - 如果查询过程中出现错误则返回错误
func (o *OrderDAO) GetOrderByID(id int) (*model.Order, error) {
	var order model.Order
	// 对应SQL:
	// 1. SELECT * FROM orders WHERE id = id LIMIT 1;
	// 2. SELECT * FROM order_items WHERE order_id = id;
	// 3. SELECT * FROM books WHERE id IN (SELECT book_id FROM order_items WHERE order_id = id);
	err := o.db.Preload("OrderItems.Book").First(&order, id).Error
	return &order, err
}

// GetOrderByOrderNo 根据订单号获取订单
// 参数:
//
//	orderNo - 订单号
//
// 返回:
//
//	*model.Order - 订单对象指针（包含订单项和图书信息）
//	error - 如果查询过程中出现错误则返回错误
func (o *OrderDAO) GetOrderByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	// 对应SQL:
	// 1. SELECT * FROM orders WHERE order_no = orderNo LIMIT 1;
	// 2. SELECT * FROM order_items WHERE order_id = (SELECT id FROM orders WHERE order_no = orderNo);
	// 3. SELECT * FROM books WHERE id IN (SELECT book_id FROM order_items WHERE order_id = (SELECT id FROM orders WHERE order_no = orderNo));
	err := o.db.Preload("OrderItems.Book").Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

// GetUserOrders 获取用户的订单列表（分页）
// 参数:
//
//	userID - 用户ID
//	page - 页码
//	pageSize - 每页大小
//
// 返回:
//
//	[]*model.Order - 订单对象切片
//	int64 - 订单总数
//	error - 如果查询过程中出现错误则返回错误
func (o *OrderDAO) GetUserOrders(userID int, page, pageSize int) ([]*model.Order, int64, error) {
	var orders []*model.Order
	var total int64

	// 获取总数
	// 对应SQL: SELECT COUNT(*) FROM orders WHERE user_id = userID;
	err := o.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	// 对应SQL:
	// 1. SELECT * FROM orders WHERE user_id = userID ORDER BY created_at DESC LIMIT pageSize OFFSET offset;
	// 2. 对于每个订单: SELECT * FROM order_items WHERE order_id = orderID;
	// 3. 对于每个订单项: SELECT * FROM books WHERE id IN (SELECT book_id FROM order_items WHERE order_id = orderID);
	offset := (page - 1) * pageSize
	err = o.db.Preload("OrderItems.Book").Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders).Error
	return orders, total, err
}

// UpdateOrderStatus 更新订单状态
// 参数:
//
//	orderID - 订单ID
//	status - 订单状态
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (o *OrderDAO) UpdateOrderStatus(orderID int, status int) error {
	// 对应SQL: UPDATE orders SET status = status WHERE id = orderID;
	err := o.db.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status).Error
	return err
}

// MarkOrderAsPaid 标记订单为已支付
// 参数:
//
//	orderID - 订单ID
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (o *OrderDAO) MarkOrderAsPaid(orderID int) error {
	now := time.Now()
	// 对应SQL: UPDATE orders SET status = 1, is_paid = true, payment_time = now WHERE id = orderID;
	err := o.db.Model(&model.Order{}).Where("id = ?", orderID).Updates(map[string]any{
		"status":       1,
		"is_paid":      true,
		"payment_time": &now,
	}).Error
	return err
}

// GenerateOrderNo 生成订单号
// 返回:
//
//	string - 生成的订单号
func (o *OrderDAO) GenerateOrderNo() string {
	// 使用时间戳生成唯一订单号
	orderNo := fmt.Sprintf("ORD%d", time.Now().UnixNano())
	return orderNo
}

// CreateOrderWithItems 创建订单和订单项（事务操作）
// 参数:
//
//	order - 订单对象指针
//	items - 订单项对象切片
//
// 返回:
//
//	error - 如果创建过程中出现错误则返回错误
func (o *OrderDAO) CreateOrderWithItems(order *model.Order, items []*model.OrderItem) error {
	// 对应SQL:
	// BEGIN TRANSACTION;
	// INSERT INTO orders (order_no, user_id, total_amount, ...) VALUES (order.OrderNo, order.UserID, order.TotalAmount, ...);
	// 对于每个订单项: INSERT INTO order_items (order_id, book_id, quantity, price, ...) VALUES (order.ID, item.BookID, item.Quantity, item.Price, ...);
	// COMMIT;
	err := o.db.Transaction(func(tx *gorm.DB) error {
		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		// 创建订单项
		for _, item := range items {
			item.OrderID = order.ID
			if err := tx.Create(item).Error; err != nil {
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
//	map[string]any - 包含订单统计信息的映射
//	error - 如果查询过程中出现错误则返回错误
func (o *OrderDAO) GetOrderStatistics(userID int) (map[string]any, error) {
	var stats struct {
		TotalOrders   int64   `json:"total_orders"`
		TotalAmount   float64 `json:"total_amount"`
		PaidOrders    int64   `json:"paid_orders"`
		PendingOrders int64   `json:"pending_orders"`
	}

	// 对应SQL: SELECT COUNT(*) as total_orders, SUM(total_amount) as total_amount FROM orders WHERE user_id = userID;
	err := o.db.Model(&model.Order{}).
		Select("COUNT(*) as total_orders, SUM(total_amount) as total_amount").
		Where("user_id = ?", userID).
		Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	// 对应SQL: SELECT COUNT(*) FROM orders WHERE user_id = userID AND is_paid = true;
	err = o.db.Model(&model.Order{}).
		Where("user_id = ? AND is_paid = ?", userID, true).
		Count(&stats.PaidOrders).Error
	if err != nil {
		return nil, err
	}

	// 对应SQL: SELECT COUNT(*) FROM orders WHERE user_id = userID AND is_paid = false;
	err = o.db.Model(&model.Order{}).
		Where("user_id = ? AND is_paid = ?", userID, false).
		Count(&stats.PendingOrders).Error
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"total_orders":   stats.TotalOrders,
		"total_amount":   stats.TotalAmount,
		"paid_orders":    stats.PaidOrders,
		"pending_orders": stats.PendingOrders,
	}, nil
}
