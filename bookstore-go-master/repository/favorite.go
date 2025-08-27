package repository

import (
	"bookstore/global"
	"bookstore/model"

	"gorm.io/gorm"
)

// FavoriteDAO 收藏数据访问对象
// 封装了所有与收藏相关的数据库操作
type FavoriteDAO struct {
	db *gorm.DB // GORM数据库连接实例
}

// NewFavoriteDAO 创建新的收藏DAO实例
// 返回:
//
//	*FavoriteDAO - 初始化后的收藏数据访问对象
func NewFavoriteDAO() *FavoriteDAO {
	return &FavoriteDAO{
		db: global.GetDB(), // 从全局变量获取数据库连接
	}
}

// AddFavorite 添加收藏
// 参数:
//
//	userID - 用户ID
//	bookID - 图书ID
//
// 返回:
//
//	error - 如果添加过程中出现错误则返回错误
func (f *FavoriteDAO) AddFavorite(userID int, bookID int) error {
	favorite := &model.Favorite{
		UserID: userID,
		BookID: bookID,
	}
	// 对应SQL: INSERT INTO favorites (user_id, book_id) VALUES (userID, bookID);
	err := f.db.Create(favorite).Error
	return err
}

// RemoveFavorite 移除收藏
// 参数:
//
//	userID - 用户ID
//	bookID - 图书ID
//
// 返回:
//
//	error - 如果移除过程中出现错误则返回错误
func (f *FavoriteDAO) RemoveFavorite(userID int, bookID int) error {
	// 对应SQL: DELETE FROM favorites WHERE user_id = userID AND book_id = bookID;
	err := f.db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&model.Favorite{}).Error
	return err
}

// CheckFavorite 检查是否已收藏
// 参数:
//
//	userID - 用户ID
//	bookID - 图书ID
//
// 返回:
//
//	bool - 如果已收藏返回true，否则返回false
//	error - 如果查询过程中出现错误则返回错误
func (f *FavoriteDAO) CheckFavorite(userID int, bookID int) (bool, error) {
	var count int64
	// 对应SQL: SELECT COUNT(*) FROM favorites WHERE user_id = userID AND book_id = bookID;
	err := f.db.Model(&model.Favorite{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserFavorites 获取用户的收藏列表
// 参数:
//
//	userID - 用户ID
//
// 返回:
//
//	[]*model.Favorite - 收藏对象切片（包含关联的Book信息）
//	error - 如果查询过程中出现错误则返回错误
func (f *FavoriteDAO) GetUserFavorites(userID int) ([]*model.Favorite, error) {
	var favorites []*model.Favorite
	// 对应SQL: SELECT * FROM favorites WHERE user_id = userID;
	// 同时预加载关联的图书信息: SELECT * FROM books WHERE id IN (SELECT book_id FROM favorites WHERE user_id = userID);
	err := f.db.Preload("Book").Where("user_id = ?", userID).Find(&favorites).Error
	return favorites, err
}

// GetUserFavoriteCount 获取用户的收藏数量
// 参数:
//
//	userID - 用户ID
//
// 返回:
//
//	int64 - 用户的收藏数量
//	error - 如果查询过程中出现错误则返回错误
func (f *FavoriteDAO) GetUserFavoriteCount(userID int) (int64, error) {
	var count int64
	// 对应SQL: SELECT COUNT(*) FROM favorites WHERE user_id = userID;
	err := f.db.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
