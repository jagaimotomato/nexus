package router

import (
	"nexus/internal/conf"
	"nexus/internal/handler"
	"nexus/internal/logger"
	"nexus/internal/middleware"
	"nexus/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *conf.Config, authHandler *handler.AuthHandler, menuHandler *handler.MenuHandler, ipSecurity *service.IPSecurityService) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 注册中间件
	r.Use(middleware.GinAccessLogger())
	r.Use(middleware.CustomRecovery(logger.Log))
	r.Use(middleware.Cors())
	r.Use(middleware.Gzip())
	r.Use(middleware.RequestID())
	r.Use(middleware.IPBlacklist(ipSecurity))
	r.Use(middleware.RateLimitWithAutoBan(cfg.RateLimit.Qps, cfg.RateLimit.Burst, ipSecurity))

	public := r.Group("/api/v1")
	private := r.Group("/api/v1")
	private.Use(middleware.JWTAuth(cfg.Jwt))

	authHandler.RegisterPublic(public)
	authHandler.RegisterPrivate(private)
	menuHandler.RegisterPrivate(private)

	return r
}
