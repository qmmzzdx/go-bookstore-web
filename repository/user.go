package repository

import (
	"bookstore/global"
	"bookstore/model"
	"fmt"

	"gorm.io/gorm"
)

// UserDAO 用户数据访问对象
// 封装所有与用户表相关的数据库操作
type UserDAO struct {
	db *gorm.DB // GORM数据库实例
}

// NewUserDAO 创建新的用户DAO实例
// 返回:
//
//	*UserDAO - 初始化好的用户数据访问对象
func NewUserDAO() *UserDAO {
	return &UserDAO{
		db: global.GetDB(), // 获取全局数据库连接
	}
}

// GetUserByID 根据ID获取用户信息
// 参数:
//
//	id - 用户ID
//
// 返回:
//
//	*model.User - 用户对象指针
//	error - 错误信息
func (u *UserDAO) GetUserByID(id int) (*model.User, error) {
	var user model.User
	// 对应SQL: SELECT * FROM users WHERE id = id LIMIT 1;
	err := u.db.First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
// 参数:
//
//	username - 用户名
//
// 返回:
//
//	*model.User - 用户对象指针
//	error - 错误信息
func (u *UserDAO) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	// 对应SQL: SELECT * FROM users WHERE username = username LIMIT 1;
	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户名不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return &user, nil
}

// CreateUser 创建新用户
// 参数:
//
//	user - 用户对象指针
//
// 返回:
//
//	error - 错误信息
func (u *UserDAO) CreateUser(user *model.User) error {
	// 对应SQL: INSERT INTO users (username, password, email, phone, ...)
	// VALUES (user.Username, user.Password, user.Email, user.Phone, ...);
	err := u.db.Create(user).Error
	if err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}
	return nil
}

// UpdateUser 更新用户信息
// 参数:
//
//	user - 用户对象指针
//
// 返回:
//
//	error - 错误信息
func (u *UserDAO) UpdateUser(user *model.User) error {
	// 对应SQL: UPDATE users SET username = user.Username, password = user.Password,
	// email = user.Email, phone = user.Phone, ... WHERE id = user.ID;
	err := u.db.Save(user).Error
	if err != nil {
		return fmt.Errorf("更新用户失败: %v", err)
	}
	return nil
}

// DeleteUser 删除用户
// 参数:
//
//	id - 用户ID
//
// 返回:
//
//	error - 错误信息
func (u *UserDAO) DeleteUser(id int) error {
	// 对应SQL: DELETE FROM users WHERE id = id;
	err := u.db.Delete(&model.User{}, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("用户不存在")
		}
		return fmt.Errorf("删除用户失败: %v", err)
	}
	return nil
}

// CheckUserExists 检查用户是否存在（通过用户名/手机号/邮箱）
// 参数:
//
//	username - 用户名
//	phone - 手机号
//	email - 邮箱
//
// 返回:
//
//	bool - 是否存在
//	error - 错误信息
func (u *UserDAO) CheckUserExists(username, phone, email string) (bool, error) {
	var count int64

	// 检查用户名是否存在
	// 对应SQL: SELECT COUNT(*) FROM users WHERE username = username;
	err := u.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查用户名失败: %v", err)
	}
	if count > 0 {
		return true, nil
	}

	// 检查手机号是否存在
	// 对应SQL: SELECT COUNT(*) FROM users WHERE phone = phone;
	err = u.db.Model(&model.User{}).Where("phone = ?", phone).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查手机号失败: %v", err)
	}
	if count > 0 {
		return true, nil
	}

	// 检查邮箱是否存在
	// 对应SQL: SELECT COUNT(*) FROM users WHERE email = email;
	err = u.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查邮箱失败: %v", err)
	}
	return count > 0, nil
}

// GetUserList 分页获取用户列表
// 参数:
//
//	page - 页码
//	pageSize - 每页数量
//
// 返回:
//
//	[]*model.User - 用户列表
//	int64 - 总记录数
//	error - 错误信息
func (u *UserDAO) GetUserList(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 获取总数
	// 对应SQL: SELECT COUNT(*) FROM users;
	err := u.db.Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取用户总数失败: %v", err)
	}

	// 分页查询
	// 对应SQL: SELECT * FROM users LIMIT pageSize OFFSET offset;
	offset := (page - 1) * pageSize
	err = u.db.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %v", err)
	}
	return users, total, nil
}
