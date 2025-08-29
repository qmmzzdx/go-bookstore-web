package model

import "time"

// Favorite 表示用户的书籍收藏关系模型
// 用于记录用户收藏的书籍信息
type Favorite struct {
	ID        int       `json:"id" gorm:"primaryKey"` // 收藏记录ID，主键
	UserID    int       `json:"user_id"`              // 用户ID，关联用户表
	BookID    int       `json:"book_id"`              // 书籍ID，关联书籍表
	CreatedAt time.Time `json:"created_at"`           // 收藏创建时间

	// Book 关联的书籍详细信息
	// 使用指针和omitempty标签，当为空时不序列化到JSON
	// gorm标签指定外键关联关系
	Book *Book `json:"book,omitempty" gorm:"foreignKey:BookID"`
}

// TableName 指定GORM使用的表名
// 实现GORM的Tabler接口，明确指定表名而不是自动推断
// 返回: 数据库表名字符串
func (f *Favorite) TableName() string {
	return "favorites" // 返回数据库中的实际表名
}
