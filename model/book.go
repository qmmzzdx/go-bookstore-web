package model

import "time"

// Book 图书模型
type Book struct {
	ID          int       `json:"id" gorm:"primaryKey"`  // 图书ID
	Title       string    `json:"title" gorm:"not null"` // 图书标题
	Author      string    `json:"author"`                // 作者
	Price       int       `json:"price"`                 // 价格（元）
	Discount    int       `json:"discount"`              // 折扣（百分比，100表示无折扣）
	Type        string    `json:"type"`                  // 图书类型
	Stock       int       `json:"stock"`                 // 库存数量
	Status      int       `json:"status"`                // 图书状态：0-下架，1-上架
	Description string    `json:"description"`           // 图书描述
	CoverURL    string    `json:"cover_url"`             // 封面图片URL
	ISBN        string    `json:"isbn"`                  // ISBN号
	Publisher   string    `json:"publisher"`             // 出版社
	PublishDate string    `json:"publish_date"`          // 出版日期
	Pages       int       `json:"pages"`                 // 页数
	Language    string    `json:"language"`              // 语言
	Format      string    `json:"format"`                // 装帧格式
	CategoryID  uint      `json:"category_id"`           // 分类ID
	Sale        int       `json:"sale"`                  // 销售量
	CreatedAt   time.Time `json:"created_at"`            // 创建时间
	UpdatedAt   time.Time `json:"updated_at"`            // 更新时间
}

// TableName 指定Book模型对应的数据库表名
func (b *Book) TableName() string {
	return "books"
}

// BookCreateRequest 创建图书请求
type BookCreateRequest struct {
	Title       string `json:"title" binding:"required"`                  // 图书标题
	Author      string `json:"author" binding:"required"`                 // 作者
	Price       int    `json:"price" binding:"required,min=0"`            // 价格（元），不能小于0
	Discount    int    `json:"discount" binding:"required,min=0,max=100"` // 折扣（百分比），0-100之间
	Type        string `json:"type" binding:"required,min=1"`             // 图书类型
	Stock       int    `json:"stock" binding:"required,min=0"`            // 库存数量，不能小于0
	Status      int    `json:"status" binding:"min=0,max=1"`              // 图书状态：0-下架，1-上架
	CoverURL    string `json:"cover_url"`                                 // 封面图片URL
	Description string `json:"description"`                               // 图书描述
	ISBN        string `json:"isbn"`                                      // ISBN号
	Publisher   string `json:"publisher"`                                 // 出版社
	PublishDate string `json:"publish_date"`                              // 出版日期
	Pages       int    `json:"pages"`                                     // 页数
	Language    string `json:"language"`                                  // 语言
	Format      string `json:"format"`                                    // 装帧格式
	CategoryID  uint   `json:"category_id"`                               // 分类ID
	Sale        int    `json:"sale" binding:"min=0"`                      // 销售量，不能小于0
}

// BookUpdateRequest 更新图书请求
type BookUpdateRequest struct {
	Title       string `json:"title"`                            // 图书标题
	Author      string `json:"author"`                           // 作者
	Price       int    `json:"price" binding:"min=0"`            // 价格（元），不能小于0
	Discount    int    `json:"discount" binding:"min=0,max=100"` // 折扣（百分比），0-100之间
	Type        string `json:"type"`                             // 图书类型
	Stock       int    `json:"stock" binding:"min=0"`            // 库存数量，不能小于0
	CoverURL    string `json:"cover_url"`                        // 封面图片URL
	Description string `json:"description"`                      // 图书描述
	ISBN        string `json:"isbn"`                             // ISBN号
	Publisher   string `json:"publisher"`                        // 出版社
	PublishDate string `json:"publish_date"`                     // 出版日期
	Pages       int    `json:"pages"`                            // 页数
	Language    string `json:"language"`                         // 语言
	Format      string `json:"format"`                           // 装帧格式
	Status      int    `json:"status" binding:"min=0,max=1"`     // 图书状态：0-下架，1-上架
	CategoryID  uint   `json:"category_id"`                      // 分类ID
	Sale        int    `json:"sale" binding:"min=0"`             // 销售量，不能小于0
}

// BookListRequest 图书列表请求
type BookListRequest struct {
	Page       int    `form:"page" binding:"min=1"`              // 页码，从1开始
	PageSize   int    `form:"page_size" binding:"min=1,max=100"` // 每页数量，1-100之间
	Title      string `form:"title"`                             // 按标题搜索
	Author     string `form:"author"`                            // 按作者搜索
	Type       string `form:"type"`                              // 按类型搜索
	Status     *int   `form:"status"`                            // 按状态搜索，nil表示不过滤状态
	CategoryID uint   `form:"category_id"`                       // 按分类ID搜索
}

// BookListResponse 图书列表响应
type BookListResponse struct {
	Books       []Book `json:"books"`        // 图书列表
	Total       int64  `json:"total"`        // 总数
	TotalPage   int    `json:"total_page"`   // 总页数
	CurrentPage int    `json:"current_page"` // 当前页码
}
