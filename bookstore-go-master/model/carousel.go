package model

import "time"

// Carousel 表示网站首页轮播图/幻灯片模型
// 用于管理首页轮播展示的图片及其相关信息
type Carousel struct {
	ID          int       `json:"id" gorm:"primaryKey"`          // 轮播图ID，主键
	Title       string    `json:"title" gorm:"not null"`         // 轮播图标题(必填)
	Description string    `json:"description" gorm:"type:text"`  // 详细描述，使用text类型存储长文本
	ImageURL    string    `json:"image_url" gorm:"not null"`     // 图片访问地址(必填)
	LinkURL     string    `json:"link_url"`                      // 点击轮播图跳转的目标URL
	SortOrder   int       `json:"sort_order" gorm:"default:0"`   // 排序权重(数字越大越靠前)
	IsActive    bool      `json:"is_active" gorm:"default:true"` // 是否显示该轮播图(默认true)
	CreatedAt   time.Time `json:"created_at"`                    // 创建时间(GORM自动维护)
	UpdatedAt   time.Time `json:"updated_at"`                    // 更新时间(GORM自动维护)
}

// TableName 指定GORM使用的表名
// 实现GORM的Tabler接口，明确指定表名而不是自动推断
// 返回值: 数据库中的实际表名"carousel"
func (Carousel) TableName() string {
	return "carousel"
}
