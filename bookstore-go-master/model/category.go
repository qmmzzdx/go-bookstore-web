package model

import "time"

// Category 图书分类模型
// 用于管理图书分类信息，包含分类的基本属性和展示特性
type Category struct {
	ID          int       `json:"id" gorm:"primaryKey"`          // 分类ID，主键
	Name        string    `json:"name" gorm:"not null;unique"`   // 分类名称(非空且唯一)
	Description string    `json:"description"`                   // 分类描述信息
	Icon        string    `json:"icon"`                          // 分类图标URL或图标标识
	Color       string    `json:"color"`                         // 分类主色(十六进制颜色码)
	Gradient    string    `json:"gradient"`                      // 渐变色配置(用于UI展示)
	Sort        int       `json:"sort" gorm:"default:0"`         // 排序权重(数字越大排序越靠前)
	IsActive    bool      `json:"is_active" gorm:"default:true"` // 是否启用该分类(默认true)
	BookCount   int       `json:"book_count" gorm:"default:0"`   // 关联图书数量(统计值)
	CreatedAt   time.Time `json:"created_at"`                    // 创建时间(GORM自动维护)
	UpdatedAt   time.Time `json:"updated_at"`                    // 更新时间(GORM自动维护)
}

// TableName 指定GORM使用的表名
// 实现GORM的Tabler接口，明确指定表名而不是自动推断
// 返回值: 数据库中的实际表名"categories"
func (c *Category) TableName() string {
	return "categories"
}
