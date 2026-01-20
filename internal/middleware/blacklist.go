package middleware

import (
	"net/http"
	"nexus/internal/response"
	"nexus/internal/service" // 引用 Service 层

	"github.com/gin-gonic/gin"
)

func IPBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 直接查 Service 的内存缓存
		if service.IsBlocked(clientIP) {
			// 使用统一响应结构，但保持 HTTP 403 状态码
			response.Result(c, http.StatusForbidden, response.Forbidden, "Access Denied: Your IP is blocked.", nil)
			
			// 重要：response.Result 只是写 JSON，不会中断后续 Handler，必须调用 Abort
			c.Abort() 
			return
		}

		c.Next()
	}
}