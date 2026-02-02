package router

import (
	"nexus/internal/conf"
	"nexus/internal/logger"
	"nexus/internal/middleware"

	"nexus/internal/api"
	"nexus/internal/handler"

	"github.com/gin-gonic/gin"
)

func InitRouter(cfg *conf.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 注册中间件
	r.Use(middleware.GinAccessLogger())
	r.Use(middleware.CustomRecovery(logger.Log))
	r.Use(middleware.Cors())
	r.Use(middleware.Gzip())
	r.Use(middleware.RequestID())
	r.Use(middleware.IPBlacklist())
	r.Use(middleware.RateLimitWithAutoBan(cfg.RateLimit.Qps, cfg.RateLimit.Burst))

	registerAPIRoutes(r)

	return r
}

func registerAPIRoutes(r *gin.Engine) {
	authAPI := &api.AuthApi{}
	
	public := r.Group("/api/v1")
	{
		public.POST("/auth/captcha", authAPI.GetCaptcha)
		public.POST("/auth/login", authAPI.Login)
	}
    menuHandler := &handler.MenuHandler{}
	private := r.Group("/api/v1")
	private.Use(middleware.JWTAuth())
	{
		private.POST("/auth/logout", authAPI.Logout)

		menu := private.Group("/menus")
		{
			menu.GET("", menuHandler.GetList)
			menu.POST("", menuHandler.Create)
			menu.PUT("/:id", menuHandler.Update)
			menu.DELETE("/:id", menuHandler.Delete)
		}
	}

}