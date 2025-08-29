package service

import (
	"bookstore/model"
	"bookstore/repository"
)

// FavoriteService 收藏服务
// 负责用户收藏相关的业务逻辑处理
type FavoriteService struct {
	favoriteDAO *repository.FavoriteDAO // 收藏数据访问对象
}

// NewFavoriteService 创建新的收藏服务实例
// 返回:
//
//	*FavoriteService - 初始化好的收藏服务
func NewFavoriteService() *FavoriteService {
	return &FavoriteService{
		favoriteDAO: repository.NewFavoriteDAO(),
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
//	error - 错误信息
func (f *FavoriteService) AddFavorite(userID, bookID int) error {
	return f.favoriteDAO.AddFavorite(userID, bookID)
}

// RemoveFavorite 取消收藏
// 参数:
//
//	userID - 用户ID
//	bookID - 图书ID
//
// 返回:
//
//	error - 错误信息
func (f *FavoriteService) RemoveFavorite(userID, bookID int) error {
	return f.favoriteDAO.RemoveFavorite(userID, bookID)
}

// IsFavorited 检查是否已收藏
// 参数:
//
//	userID - 用户ID
//	bookID - 图书ID
//
// 返回:
//
//	bool - 是否已收藏(true:已收藏 false:未收藏)
//	error - 错误信息
func (f *FavoriteService) IsFavorited(userID, bookID int) (bool, error) {
	return f.favoriteDAO.CheckFavorite(userID, bookID)
}

// GetUserFavorites 获取用户收藏列表
// 参数:
//
//	userID - 用户ID
//	page - 页码
//	pageSize - 每页数量
//	timeFilter - 时间筛选条件
//
// 返回:
//
//	[]*model.Favorite - 收藏列表
//	int64 - 总记录数
//	error - 错误信息
func (f *FavoriteService) GetUserFavorites(userID int, page, pageSize int, timeFilter string) ([]*model.Favorite, int64, error) {
	favorites, err := f.favoriteDAO.GetUserFavorites(userID)
	if err != nil {
		return nil, 0, err
	}

	// 简单的分页实现
	total := int64(len(favorites))
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(favorites) {
		return []*model.Favorite{}, total, nil
	}

	if end > len(favorites) {
		end = len(favorites)
	}

	return favorites[start:end], total, nil
}

// GetUserFavoriteCount 获取用户收藏数量
// 参数:
//
//	userID - 用户ID
//
// 返回:
//
//	int64 - 收藏数量
//	error - 错误信息
func (f *FavoriteService) GetUserFavoriteCount(userID int) (int64, error) {
	return f.favoriteDAO.GetUserFavoriteCount(userID)
}
