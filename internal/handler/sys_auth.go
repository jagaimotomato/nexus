package handler

import (
	"nexus/internal/logger"
	"nexus/internal/response"
	"nexus/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthApi struct{}

// 获取验证码接口
func (a *AuthApi) GetCaptcha(c *gin.Context) {
	id, b64s, err := service.GenerateCaptcha()
	if err != nil {
		logger.Log.Error("验证码生成失败", zap.Error(err))
		response.FailWithMessage(c, "验证码生成失败")
		return
	}
	response.OKWithData(c, gin.H{
		"captchaId":  id,
		"captchaImg": b64s,
	})
}

// Login
func (a *AuthApi) Login(c *gin.Context) {
	var req service.LoginParams
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.InvalidParams)
		return
	}

	// 调用service
	res, err := service.Login(req)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	response.OKWithData(c, res)
}

func (a *AuthApi) Logout(c *gin.Context) {
	// 1. 获取当前用户 ID (从中间件 JWTAuth 设置的上下文获取)
	// 虽然 JWT 是无状态的，不能像 Session 那样直接销毁，但这里可以做业务层面的清理
	userID, exists := c.Get("userID")
	if !exists {
		// 如果中间件通过了但这里拿不到 ID，属于异常情况
		response.FailWithCode(c, response.AuthFailed)
		return
	}

	username, _ := c.Get("username")
	
	// 2. 记录登出日志
	logger.Log.Info("用户发起退出登录", 
		zap.Any("user_id", userID), 
		zap.Any("username", username),
	)

	// 3. (可选) 高级安全处理：
	// 如果需要强制失效 Token，可以在这里解析 Header 里的 Token，
	// 并将其存入 Redis 的 "token_blacklist" 中，设置过期时间为 Token 的剩余有效期。
	// 然后在 JWT 中间件里增加一步检查 Redis 黑名单。
	// 简单场景下，只需前端丢弃 Token 即可。

	response.OKWithMessage(c, "退出成功")
}