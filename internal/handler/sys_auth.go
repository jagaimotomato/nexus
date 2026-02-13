package handler

import (
	"nexus/internal/logger"
	"nexus/internal/response"
	"nexus/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

func (h *AuthHandler) RegisterPublic(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/captcha", h.GetCaptcha)
	}
}

func (h *AuthHandler) RegisterPrivate(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/logout", h.Logout)
	}
}

func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	id, b64s, err := h.svc.GetCaptcha()
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

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginParams
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithCode(c, response.InvalidParams)
		return
	}

	res, err := h.svc.Login(req)
	if err != nil {
		response.FailWithMessage(c, err.Error())
		return
	}

	response.OKWithData(c, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.FailWithCode(c, response.AuthFailed)
		return
	}

	username, _ := c.Get("username")

	logger.Log.Info("用户发起退出登录",
		zap.Any("user_id", userID),
		zap.Any("username", username),
	)

	h.svc.Logout()

	response.OKWithMessage(c, "退出成功")
}
