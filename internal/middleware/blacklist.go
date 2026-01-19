package middleware

import (
	"net/http"
	"nexus/internal/service" // 引用 Service 层

	"github.com/gin-gonic/gin"
)

func IPBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 直接查 Service 的内存缓存
		if service.IsBlocked(clientIP) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access Denied: Your IP is blocked.",
			})
			return // 必须 return
		}

		c.Next()
	}
}