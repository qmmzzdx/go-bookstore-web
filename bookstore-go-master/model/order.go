package model

import "time"

// Order 订单模型
type Order struct {
	ID          int        `json:"id" gorm:"primaryKey"`            // 订单ID
	UserID      int        `json:"user_id" gorm:"not null"`         // 用户ID
	OrderNo     string     `json:"order_no" gorm:"not null;unique"` // 订单号
	TotalAmount int        `json:"total_amount" gorm:"not null"`    // 订单总金额（分）
	Status      int        `json:"status" gorm:"default:0"`         // 订单状态：0-待支付，1-已支付，2-已取消
	IsPaid      bool       `json:"is_paid" gorm:"default:false"`    // 是否已支付
	PaymentTime *time.Time `json:"payment_time"`                    // 支付时间
	CreatedAt   time.Time  `json:"created_at"`                      // 创建时间
	UpdatedAt   time.Time  `json:"updated_at"`                      // 更新时间

	// 关联字段
	User       *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`         // 关联的用户信息
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"` // 订单项列表
}

// TableName 指定Order模型对应的数据库表名
func (o *Order) TableName() string {
	return "orders"
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        int       `json:"id" gorm:"primaryKey"`     // 订单项ID
	OrderID   int       `json:"order_id" gorm:"not null"` // 订单ID
	BookID    int       `json:"book_id" gorm:"not null"`  // 图书ID
	Quantity  int       `json:"quantity" gorm:"not null"` // 购买数量
	Price     int       `json:"price" gorm:"not null"`    // 单价（分）
	Subtotal  int       `json:"subtotal" gorm:"not null"` // 小计金额（分）
	CreatedAt time.Time `json:"created_at"`               // 创建时间
	UpdatedAt time.Time `json:"updated_at"`               // 更新时间

	// 关联字段
	Book *Book `gorm:"foreignKey:BookID" json:"book,omitempty"` // 关联的图书信息
}

// TableName 指定OrderItem模型对应的数据库表名
func (oi *OrderItem) TableName() string {
	return "order_items"
}
