package controller

import (
	"bookstore/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CaptchaController 验证码控制器
// 负责处理验证码相关的HTTP请求，包括生成和验证验证码
type CaptchaController struct {
	CaptchaService *service.CaptchaService // 验证码服务层实例
}

// NewCaptchaController 创建新的验证码控制器实例
// 返回:
//
//	*CaptchaController - 初始化好的验证码控制器
func NewCaptchaController() *CaptchaController {
	return &CaptchaController{
		CaptchaService: service.NewCaptchaService(), // 初始化验证码服务
	}
}

// GenerateCaptcha 生成验证码
// 路由: GET /captcha/generate
// 功能: 生成新的验证码图片和对应ID
func (c *CaptchaController) GenerateCaptcha(ctx *gin.Context) {
	// 调用服务层生成验证码
	captchaResponse, err := c.CaptchaService.GenerateCaptcha()
	if err != nil {
		// 生成失败返回500错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,          // 错误代码
			"message": "生成验证码失败",   // 错误信息
			"error":   err.Error(), // 详细错误
		})
		return
	}

	// 返回生成的验证码信息
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,               // 成功代码
		"data":    captchaResponse, // 验证码数据(包含ID和Base64图片)
		"message": "验证码生成成功",       // 成功信息
	})
}

// VerifyCaptcha 验证验证码
// 路由: POST /captcha/verify
// 功能: 验证用户输入的验证码是否正确
func (c *CaptchaController) VerifyCaptcha(ctx *gin.Context) {
	// 定义请求参数结构
	var req struct {
		CaptchaID    string `json:"captcha_id" binding:"required"`    // 验证码ID(必填)
		CaptchaValue string `json:"captcha_value" binding:"required"` // 用户输入的验证码值(必填)
	}

	// 绑定并验证请求参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 参数错误返回400
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,          // 错误代码
			"message": "请求参数错误",    // 错误信息
			"error":   err.Error(), // 详细错误
		})
		return
	}

	// 调用服务层验证验证码
	isValid := c.CaptchaService.VerifyCaptcha(req.CaptchaID, req.CaptchaValue)

	// 根据验证结果返回响应
	if isValid {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,       // 成功代码
			"message": "验证码正确", // 成功信息
		})
	} else {
		// 验证失败返回400
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,      // 错误代码
			"message": "验证码错误", // 错误信息
		})
	}
}
