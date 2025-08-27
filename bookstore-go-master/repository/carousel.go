package repository

import (
	"bookstore/global"
	"bookstore/model"

	"gorm.io/gorm"
)

// CarouselDAO 轮播图数据访问对象
// 封装了所有与轮播图相关的数据库操作
type CarouselDAO struct {
	db *gorm.DB // GORM数据库连接实例
}

// NewCarouselDAO 创建新的轮播图DAO实例
// 返回:
//
//	*CarouselDAO - 初始化后的轮播图数据访问对象
func NewCarouselDAO() *CarouselDAO {
	return &CarouselDAO{
		db: global.GetDB(), // 从全局变量获取数据库连接
	}
}

// GetAllCarousels 获取所有轮播图
// 返回:
//
//	[]*model.Carousel - 轮播图对象切片
//	error - 如果查询过程中出现错误则返回错误
func (c *CarouselDAO) GetAllCarousels() ([]*model.Carousel, error) {
	var carousels []*model.Carousel
	// 对应SQL: SELECT * FROM carousels;
	err := c.db.Find(&carousels).Error
	return carousels, err
}

// GetActiveCarousels 获取活跃的轮播图
// 返回:
//
//	[]*model.Carousel - 活跃的轮播图对象切片
//	error - 如果查询过程中出现错误则返回错误
func (c *CarouselDAO) GetActiveCarousels() ([]*model.Carousel, error) {
	var carousels []*model.Carousel
	// 对应SQL: SELECT * FROM carousels WHERE is_active = true ORDER BY sort_order ASC;
	err := c.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&carousels).Error
	return carousels, err
}

// GetCarouselByID 根据ID获取轮播图
// 参数:
//
//	id - 轮播图ID
//
// 返回:
//
//	*model.Carousel - 轮播图对象指针
//	error - 如果查询过程中出现错误则返回错误
func (c *CarouselDAO) GetCarouselByID(id int) (*model.Carousel, error) {
	var carousel model.Carousel
	// 对应SQL: SELECT * FROM carousels WHERE id = id LIMIT 1;
	err := c.db.First(&carousel, id).Error
	return &carousel, err
}

// CreateCarousel 创建轮播图
// 参数:
//
//	carousel - 轮播图对象指针
//
// 返回:
//
//	error - 如果创建过程中出现错误则返回错误
func (c *CarouselDAO) CreateCarousel(carousel *model.Carousel) error {
	// 对应SQL: INSERT INTO carousels (image_url, link, is_active, sort_order, ...) VALUES (carousel.ImageURL, carousel.Link, carousel.IsActive, carousel.SortOrder, ...);
	err := c.db.Create(carousel).Error
	return err
}

// UpdateCarousel 更新轮播图
// 参数:
//
//	carousel - 轮播图对象指针
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (c *CarouselDAO) UpdateCarousel(carousel *model.Carousel) error {
	// 对应SQL: UPDATE carousels SET image_url = carousel.ImageURL, link = carousel.Link, is_active = carousel.IsActive, sort_order = carousel.SortOrder, ... WHERE id = carousel.ID;
	err := c.db.Save(carousel).Error
	return err
}

// DeleteCarousel 删除轮播图
// 参数:
//
//	id - 轮播图ID
//
// 返回:
//
//	error - 如果删除过程中出现错误则返回错误
func (c *CarouselDAO) DeleteCarousel(id int) error {
	// 对应SQL: DELETE FROM carousels WHERE id = id;
	err := c.db.Delete(&model.Carousel{}, id).Error
	return err
}
