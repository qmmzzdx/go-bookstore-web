package service

import (
	"encoding/base64"
	"errors"

	"bookstore/jwt"
	"bookstore/model"
	"bookstore/repository"
)

// UserService 用户服务层
// 处理用户相关的业务逻辑，包括注册、登录、信息修改等
type UserService struct {
	UserDB *repository.UserDAO // 用户数据访问对象
}

// LoginResponse 登录响应数据结构
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`  // 访问令牌
	RefreshToken string    `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int64     `json:"expires_in"`    // 过期时间(秒)
	UserInfo     *UserInfo `json:"user_info"`     // 用户基本信息
}

// UserInfo 用户基本信息结构
type UserInfo struct {
	ID       int    `json:"id"`       // 用户ID
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱
	Phone    string `json:"phone"`    // 手机号
}

// NewUserService 创建新的用户服务实例
// 返回:
//
//	*UserService - 初始化好的用户服务实例
func NewUserService() *UserService {
	return &UserService{
		UserDB: repository.NewUserDAO(), // 初始化用户DAO
	}
}

// UserRegister 用户注册服务
// 参数:
//
//	username - 用户名
//	password - 密码(明文)
//	phone - 手机号
//	email - 邮箱
//
// 返回:
//
//	error - 错误信息
func (u *UserService) UserRegister(username, password, phone, email string) error {
	// 1. 检查用户名、邮箱、手机号唯一性
	exists, err := u.checkUserExists(username, phone, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名、邮箱或手机号已存在")
	}

	// 2. 密码编码（使用base64）
	encodedPwd := u.encodePassword(password)

	// 3. 调用 DAO 层插入用户
	err = u.createUser(username, encodedPwd, phone, email)
	if err != nil {
		return err
	}
	return nil
}

// encodePassword 密码base64编码
// 参数:
//
//	password - 明文密码
//
// 返回:
//
//	string - base64编码后的密码
func (u *UserService) encodePassword(password string) string {
	// 使用 base64 编码密码
	encodedPwd := base64.StdEncoding.EncodeToString([]byte(password))
	return encodedPwd
}

// checkUserExists 检查用户是否存在
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
func (u *UserService) checkUserExists(username, phone, email string) (bool, error) {
	return u.UserDB.CheckUserExists(username, phone, email)
}

// createUser 创建新用户
// 参数:
//
//	username - 用户名
//	passwordHash - 加密后的密码
//	phone - 手机号
//	email - 邮箱
//
// 返回:
//
//	error - 错误信息
func (u *UserService) createUser(username, passwordHash, phone, email string) error {
	user := &model.User{
		Username: username,
		Password: passwordHash,
		Phone:    phone,
		Email:    email,
	}
	return u.UserDB.CreateUser(user)
}

// UserLogin 用户登录（不带验证码）
// 参数:
//
//	username - 用户名
//	password - 密码(明文)
//
// 返回:
//
//	*LoginResponse - 登录响应数据
//	error - 错误信息
func (u *UserService) UserLogin(username, password string) (*LoginResponse, error) {
	// 1. 获取用户信息
	user, err := u.UserDB.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 2. 验证密码
	if !u.verifyPassword(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	// 3. 生成 JWT Token 对
	tokenResponse, err := jwt.GenerateTokenPair(uint(user.ID), user.Username)
	if err != nil {
		return nil, errors.New("生成 token 失败")
	}

	// 4. 构建登录响应
	response := &LoginResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
		UserInfo: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
		},
	}
	return response, nil
}

// verifyPassword 验证密码
// 参数:
//
//	inputPassword - 用户输入的密码(明文)
//	storedPassword - 存储的加密密码
//
// 返回:
//
//	bool - 验证结果
func (u *UserService) verifyPassword(inputPassword, storedPassword string) bool {
	// 进行base64编码比较
	return u.encodePassword(inputPassword) == storedPassword
}

// UserLoginWithCaptcha 带验证码的登录
// 参数:
//
//	username - 用户名
//	password - 密码(明文)
//	captchaId - 验证码ID
//	captcha - 验证码值
//
// 返回:
//
//	*LoginResponse - 登录响应数据
//	error - 错误信息
func (u *UserService) UserLoginWithCaptcha(username, password, captchaId, captcha string) (*LoginResponse, error) {
	// 1. 验证验证码
	captchaService := NewCaptchaService()
	if !captchaService.VerifyCaptcha(captchaId, captcha) {
		return nil, errors.New("验证码错误")
	}

	// 2. 获取用户信息
	user, err := u.UserDB.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 3. 验证密码
	if !u.verifyPassword(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	// 4. 生成 JWT Token 对
	tokenResponse, err := jwt.GenerateTokenPair(uint(user.ID), user.Username)
	if err != nil {
		return nil, errors.New("生成 token 失败")
	}

	// 5. 构建登录响应
	response := &LoginResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresIn:    tokenResponse.ExpiresIn,
		UserInfo: &UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
		},
	}
	return response, nil
}

// GetUserByID 根据用户ID获取用户信息
// 参数:
//
//	userID - 用户ID
//
// 返回:
//
//	*model.User - 用户对象
//	error - 错误信息
func (u *UserService) GetUserByID(userID int) (*model.User, error) {
	user, err := u.UserDB.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return user, nil
}

// UpdateUserInfo 更新用户信息
// 参数:
//
//	user - 用户对象(包含更新信息)
//
// 返回:
//
//	error - 错误信息
func (u *UserService) UpdateUserInfo(user *model.User) error {
	// 1. 检查用户是否存在
	existingUser, err := repository.NewUserDAO().GetUserByID(user.ID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 2. 更新用户信息
	existingUser.Phone = user.Phone
	existingUser.Email = user.Email
	existingUser.Avatar = user.Avatar

	// 3. 调用 DAO 层更新用户信息
	err = u.UserDB.UpdateUser(existingUser)
	if err != nil {
		return errors.New("更新用户信息失败")
	}
	return nil
}

// ChangePassword 修改密码
// 参数:
//
//	userID - 用户ID
//	oldPassword - 旧密码(明文)
//	newPassword - 新密码(明文)
//
// 返回:
//
//	error - 错误信息
func (u *UserService) ChangePassword(userID int, oldPassword, newPassword string) error {
	// 1. 获取用户信息
	user, err := u.UserDB.GetUserByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 2. 验证旧密码
	if !u.verifyPassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 3. 编码新密码
	encodedNewPassword := u.encodePassword(newPassword)

	// 4. 更新密码
	user.Password = encodedNewPassword
	err = u.UserDB.UpdateUser(user)
	if err != nil {
		return errors.New("密码修改失败")
	}
	return nil
}
