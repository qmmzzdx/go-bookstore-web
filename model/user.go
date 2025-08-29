package model

import "time"

// User 用户模型
type User struct {
	ID        int       `json:"id" `                             // 用户ID
	Username  string    `json:"username" gorm:"unique;not null"` // 用户名，唯一且不能为空
	Password  string    `json:"-" gorm:"not null"`               // 密码，不返回给前端
	Email     string    `json:"email" gorm:"unique;not null"`    // 邮箱，唯一且不能为空
	Phone     string    `json:"phone"`                           // 手机号码
	Avatar    string    `json:"avatar"`                          // 头像URL
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`   // 是否为管理员
	CreatedAt time.Time `json:"created_at"`                      // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                      // 更新时间
}

// TableName 指定User模型对应的数据库表名
func (u *User) TableName() string {
	return "users"
}
