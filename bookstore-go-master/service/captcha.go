package service

import (
	"context"
	"fmt"
	"time"

	"bookstore/global"

	"github.com/mojocn/base64Captcha"
)

// CaptchaService 验证码服务
// 负责验证码的生成、验证和清理工作
type CaptchaService struct {
	store base64Captcha.Store // 验证码存储接口
}

// NewCaptchaService 创建新的验证码服务实例
// 返回:
//
//	*CaptchaService - 初始化好的验证码服务
func NewCaptchaService() *CaptchaService {
	return &CaptchaService{
		// 使用内存存储作为默认存储（实际使用Redis存储）
		store: base64Captcha.DefaultMemStore,
	}
}

// CaptchaResponse 验证码响应数据结构
// 用于返回给客户端的验证码信息
type CaptchaResponse struct {
	CaptchaID     string `json:"captcha_id"`     // 验证码唯一标识
	CaptchaBase64 string `json:"captcha_base64"` // base64编码的验证码图片
}

// GenerateCaptcha 生成验证码
// 返回:
//
//	*CaptchaResponse - 生成的验证码数据
//	error - 错误信息
func (c *CaptchaService) GenerateCaptcha() (*CaptchaResponse, error) {
	// 配置数字验证码驱动参数
	driver := base64Captcha.NewDriverDigit(
		80,  // 图片高度(像素)
		240, // 图片宽度(像素)
		4,   // 验证码字符长度
		0.7, // 干扰线强度(0-1)
		80,  // 干扰点数量
	)

	// 创建验证码实例
	captcha := base64Captcha.NewCaptcha(driver, c.store)

	// 生成验证码(返回ID、base64图片、答案)
	id, b64s, answer, err := captcha.Generate()
	if err != nil {
		return nil, fmt.Errorf("生成验证码失败: %v", err)
	}

	// 将验证码答案存储到Redis，设置5分钟过期时间
	ctx := context.Background()
	redisKey := fmt.Sprintf("captcha:%s", id)
	err = global.RedisClient.Set(ctx, redisKey, answer, 5*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("存储验证码失败: %v", err)
	}

	// 返回生成的验证码数据
	return &CaptchaResponse{
		CaptchaID:     id,
		CaptchaBase64: b64s,
	}, nil
}

// VerifyCaptcha 验证验证码
// 参数:
//
//	captchaID - 验证码ID
//	captchaValue - 用户输入的验证码值
//
// 返回:
//
//	bool - 验证结果(true:验证成功 false:验证失败)
func (c *CaptchaService) VerifyCaptcha(captchaID, captchaValue string) bool {
	// 空值检查
	if captchaID == "" || captchaValue == "" {
		return false
	}

	// 从Redis获取存储的验证码答案
	ctx := context.Background()
	redisKey := fmt.Sprintf("captcha:%s", captchaID)
	storedAnswer, err := global.RedisClient.Get(ctx, redisKey).Result()
	if err != nil {
		return false
	}

	// 比较用户输入和存储的答案(不区分大小写)
	// 验证成功后删除Redis中的验证码(防止重复使用)
	if storedAnswer == captchaValue {
		global.RedisClient.Del(ctx, redisKey)
		return true
	}
	return false
}

// CleanExpiredCaptcha 清理过期的验证码
// 定期清理Redis中过期的验证码记录
func (c *CaptchaService) CleanExpiredCaptcha() {
	ctx := context.Background()

	// 获取所有captcha开头的key
	keys, err := global.RedisClient.Keys(ctx, "captcha:*").Result()
	if err != nil {
		return
	}

	// 检查每个key的TTL(生存时间)
	for _, key := range keys {
		ttl, err := global.RedisClient.TTL(ctx, key).Result()
		if err == nil && ttl <= 0 {
			// 如果已过期则删除
			global.RedisClient.Del(ctx, key)
		}
	}
}
