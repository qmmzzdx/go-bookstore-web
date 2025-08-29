package service

import (
	"bookstore/model"
	"bookstore/repository"
)

// CarouselService 轮播图服务
// 负责轮播图相关的业务逻辑处理
type CarouselService struct {
	CarouselDB *repository.CarouselDAO // 轮播图数据访问对象
}

// NewCarouselService 创建新的轮播图服务实例
// 返回:
//
//	*CarouselService - 初始化好的轮播图服务
func NewCarouselService() *CarouselService {
	return &CarouselService{
		CarouselDB: repository.NewCarouselDAO(),
	}
}

// GetActiveCarousels 获取激活的轮播图
// 返回:
//
//	[]*model.Carousel - 激活的轮播图列表
//	error - 错误信息
func (c *CarouselService) GetActiveCarousels() ([]*model.Carousel, error) {
	return c.CarouselDB.GetActiveCarousels()
}

// GetCarouselByID 根据ID获取轮播图
// 参数:
//
//	id - 轮播图ID
//
// 返回:
//
//	*model.Carousel - 轮播图对象指针
//	error - 错误信息
func (c *CarouselService) GetCarouselByID(id int) (*model.Carousel, error) {
	return c.CarouselDB.GetCarouselByID(id)
}

// CreateCarousel 创建轮播图
// 参数:
//
//	carousel - 轮播图对象指针
//
// 返回:
//
//	error - 错误信息
func (c *CarouselService) CreateCarousel(carousel *model.Carousel) error {
	return c.CarouselDB.CreateCarousel(carousel)
}

// UpdateCarousel 更新轮播图
// 参数:
//
//	carousel - 轮播图对象指针
//
// 返回:
//
//	error - 错误信息
func (c *CarouselService) UpdateCarousel(carousel *model.Carousel) error {
	return c.CarouselDB.UpdateCarousel(carousel)
}

// DeleteCarousel 删除轮播图
// 参数:
//
//	id - 轮播图ID
//
// 返回:
//
//	error - 错误信息
func (c *CarouselService) DeleteCarousel(id int) error {
	return c.CarouselDB.DeleteCarousel(id)
}

// GetAllCarousels 获取所有轮播图
// 返回:
//
//	[]*model.Carousel - 所有轮播图列表
//	error - 错误信息
func (c *CarouselService) GetAllCarousels() ([]*model.Carousel, error) {
	return c.CarouselDB.GetAllCarousels()
}
