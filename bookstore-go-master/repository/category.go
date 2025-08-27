package repository

import (
	"bookstore/global"
	"bookstore/model"

	"gorm.io/gorm"
)

// CategoryDAO 分类数据访问对象
// 封装了所有与分类相关的数据库操作
type CategoryDAO struct {
	db *gorm.DB // GORM数据库连接实例
}

// NewCategoryDAO 创建新的分类DAO实例
// 返回:
//
//	*CategoryDAO - 初始化后的分类数据访问对象
func NewCategoryDAO() *CategoryDAO {
	return &CategoryDAO{
		db: global.GetDB(), // 从全局变量获取数据库连接
	}
}

// GetAllCategories 获取所有分类
// 返回:
//
//	[]*model.Category - 分类对象切片
//	error - 如果查询过程中出现错误则返回错误
func (c *CategoryDAO) GetAllCategories() ([]*model.Category, error) {
	var categories []*model.Category
	// 对应SQL: SELECT * FROM categories;
	err := c.db.Find(&categories).Error
	return categories, err
}

// GetCategoryByID 根据ID获取分类
// 参数:
//
//	id - 分类ID
//
// 返回:
//
//	*model.Category - 分类对象指针
//	error - 如果查询过程中出现错误则返回错误
func (c *CategoryDAO) GetCategoryByID(id int) (*model.Category, error) {
	var category model.Category
	// 对应SQL: SELECT * FROM categories WHERE id = id LIMIT 1;
	err := c.db.First(&category, id).Error
	return &category, err
}

// CreateCategory 创建分类
// 参数:
//
//	category - 分类对象指针
//
// 返回:
//
//	error - 如果创建过程中出现错误则返回错误
func (c *CategoryDAO) CreateCategory(category *model.Category) error {
	// 对应SQL: INSERT INTO categories (name, description, ...) VALUES (category.Name, category.Description, ...);
	err := c.db.Create(category).Error
	return err
}

// UpdateCategory 更新分类
// 参数:
//
//	category - 分类对象指针
//
// 返回:
//
//	error - 如果更新过程中出现错误则返回错误
func (c *CategoryDAO) UpdateCategory(category *model.Category) error {
	// 对应SQL: UPDATE categories SET name = category.Name, description = category.Description, ... WHERE id = category.ID;
	err := c.db.Save(category).Error
	return err
}

// DeleteCategory 删除分类
// 参数:
//
//	id - 分类ID
//
// 返回:
//
//	error - 如果删除过程中出现错误则返回错误
func (c *CategoryDAO) DeleteCategory(id int) error {
	// 对应SQL: DELETE FROM categories WHERE id = id;
	err := c.db.Delete(&model.Category{}, id).Error
	return err
}
