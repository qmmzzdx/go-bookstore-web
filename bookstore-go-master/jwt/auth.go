package jwt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bookstore/global"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret JWT密钥（建议通过配置文件或环境变量设置）
var jwtSecret = []byte("bookstore_secret_key")

// Token expiration constants
const (
	// AccessTokenExpire 访问token过期时间 (2小时)
	AccessTokenExpire = 2 * time.Hour
	// RefreshTokenExpire 刷新token过期时间 (7天)
	RefreshTokenExpire = 7 * 24 * time.Hour
)

// Claims JWT声明结构体，包含用户信息和标准声明
type Claims struct {
	UserID    uint   `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	TokenType string `json:"token_type"` // token类型："access" 或 "refresh"
	jwt.RegisteredClaims
}

// TokenResponse token响应结构体，返回给客户端
type TokenResponse struct {
	AccessToken  string `json:"access_token"`  // 访问token
	RefreshToken string `json:"refresh_token"` // 刷新token
	ExpiresIn    int64  `json:"expires_in"`    // 过期时间（秒）
}

// GenerateTokenPair 生成访问token和刷新token对
// 参数:
//   - userID: 用户ID
//   - username: 用户名
//
// 返回:
//   - *TokenResponse: 包含access token和refresh token的响应
//   - error: 生成过程中遇到的错误
func GenerateTokenPair(userID uint, username string) (*TokenResponse, error) {
	// 生成访问token
	accessClaims := Claims{
		UserID:    userID,
		Username:  username,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("生成访问token失败: %v", err)
	}

	// 生成刷新token
	refreshClaims := Claims{
		UserID:    userID,
		Username:  username,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("生成刷新token失败: %v", err)
	}

	// 将token存储到Redis
	if err := StoreTokenInRedis(userID, accessTokenString, refreshTokenString); err != nil {
		return nil, fmt.Errorf("存储token到Redis失败: %v", err)
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(AccessTokenExpire.Seconds()),
	}, nil
}

// GenerateToken 兼容旧的接口，只生成访问token
// 参数:
//   - userID: 用户ID
//   - username: 用户名
//
// 返回:
//   - string: 访问token
//   - error: 生成过程中遇到的错误
func GenerateToken(userID uint, username string) (string, error) {
	tokenResponse, err := GenerateTokenPair(userID, username)
	if err != nil {
		return "", err
	}
	return tokenResponse.AccessToken, nil
}

// ParseToken 解析和校验JWT Token
// 参数:
//   - tokenString: 要解析的token字符串
//
// 返回:
//   - *Claims: 解析出的声明信息
//   - error: 解析或验证过程中遇到的错误
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token解析失败: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// 检查token是否在Redis中被撤销
		if !IsTokenValidInRedis(claims.UserID, tokenString, claims.TokenType) {
			return nil, errors.New("token已被撤销")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// StoreTokenInRedis 将token存储到Redis
// 参数:
//   - userID: 用户ID
//   - accessToken: 访问token
//   - refreshToken: 刷新token
//
// 返回:
//   - error: 存储过程中遇到的错误
func StoreTokenInRedis(userID uint, accessToken, refreshToken string) error {
	ctx := context.Background()
	userKey := fmt.Sprintf("user_tokens:%d", userID)

	// 使用hash存储用户的token信息
	err := global.RedisClient.HMSet(ctx, userKey,
		"access_token", accessToken,
		"refresh_token", refreshToken,
		"created_at", time.Now().Unix(),
	).Err()
	if err != nil {
		return fmt.Errorf("存储token到Redis失败: %v", err)
	}

	// 设置过期时间为刷新token的过期时间
	return global.RedisClient.Expire(ctx, userKey, RefreshTokenExpire).Err()
}

// IsTokenValidInRedis 检查token是否在Redis中有效
// 参数:
//   - userID: 用户ID
//   - token: 要检查的token
//   - tokenType: token类型 ("access" 或 "refresh")
//
// 返回:
//   - bool: token是否有效
func IsTokenValidInRedis(userID uint, token string, tokenType string) bool {
	ctx := context.Background()
	userKey := fmt.Sprintf("user_tokens:%d", userID)

	var redisToken string
	var err error

	if tokenType == "access" {
		redisToken, err = global.RedisClient.HGet(ctx, userKey, "access_token").Result()
	} else {
		redisToken, err = global.RedisClient.HGet(ctx, userKey, "refresh_token").Result()
	}

	if err != nil {
		return false
	}

	return redisToken == token
}

// RefreshAccessToken 使用刷新token生成新的访问token
// 参数:
//   - refreshToken: 刷新token
//
// 返回:
//   - *TokenResponse: 包含新token的响应
//   - error: 刷新过程中遇到的错误
func RefreshAccessToken(refreshToken string) (*TokenResponse, error) {
	// 解析刷新token
	claims, err := ParseToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("刷新token无效: %v", err)
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("无效的刷新token")
	}

	// 生成新的token对
	return GenerateTokenPair(claims.UserID, claims.Username)
}

// RevokeToken 撤销用户的所有token
// 参数:
//   - userID: 用户ID
//
// 返回:
//   - error: 撤销过程中遇到的错误
func RevokeToken(userID uint) error {
	ctx := context.Background()
	userKey := fmt.Sprintf("user_tokens:%d", userID)
	return global.RedisClient.Del(ctx, userKey).Err()
}

// RevokeAllUserTokens 撤销所有用户的token（用于安全事件）
// 返回:
//   - error: 撤销过程中遇到的错误
func RevokeAllUserTokens() error {
	ctx := context.Background()
	// 获取所有用户token的key
	keys, err := global.RedisClient.Keys(ctx, "user_tokens:*").Result()
	if err != nil {
		return fmt.Errorf("获取用户token列表失败: %v", err)
	}

	// 如果有token存在，则全部删除
	if len(keys) > 0 {
		return global.RedisClient.Del(ctx, keys...).Err()
	}
	return nil
}
