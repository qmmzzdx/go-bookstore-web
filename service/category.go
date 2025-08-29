package service

import (
	"bookstore/model"
	"bookstore/repository"
)

// CategoryService 分类服务
// 负责分类相关的业务逻辑处理
type CategoryService struct {
	CategoryDAO *repository.CategoryDAO // 分类数据访问对象
}

// NewCategoryService 创建新的分类服务实例
// 返回:
//
//	*CategoryService - 初始化好的分类服务
func NewCategoryService() *CategoryService {
	return &CategoryService{
		CategoryDAO: repository.NewCategoryDAO(),
	}
}

// GetAllCategories 获取所有分类
// 返回:
//
//	[]*model.Category - 所有分类列表
//	error - 错误信息
func (c *CategoryService) GetAllCategories() ([]*model.Category, error) {
	return c.CategoryDAO.GetAllCategories()
}

// GetCategoryByID 根据ID获取分类
// 参数:
//
//	id - 分类ID
//
// 返回:
//
//	*model.Category - 分类对象指针
//	error - 错误信息
func (c *CategoryService) GetCategoryByID(id int) (*model.Category, error) {
	return c.CategoryDAO.GetCategoryByID(id)
}

// CreateCategory 创建分类
// 参数:
//
//	category - 分类对象指针
//
// 返回:
//
//	error - 错误信息
func (c *CategoryService) CreateCategory(category *model.Category) error {
	return c.CategoryDAO.CreateCategory(category)
}

// UpdateCategory 更新分类
// 参数:
//
//	category - 分类对象指针
//
// 返回:
//
//	error - 错误信息
func (c *CategoryService) UpdateCategory(category *model.Category) error {
	return c.CategoryDAO.UpdateCategory(category)
}

// DeleteCategory 删除分类
// 参数:
//
//	id - 分类ID
//
// 返回:
//
//	error - 错误信息
func (c *CategoryService) DeleteCategory(id int) error {
	return c.CategoryDAO.DeleteCategory(id)
}
