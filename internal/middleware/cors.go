package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true // 开发环境允许所有
    // 生产环境建议指定域名: config.AllowOrigins = []string{"https://my-admin.com"}
    config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Request-Id"}
    return cors.New(config)
}