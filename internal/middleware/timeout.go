package middleware

import (
	"net/http"
	"time"

	"nexus/internal/response" // 引入 response 包

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

// Timeout 封装官方库
func Timeout(d time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(d),
		timeout.WithHandler(func(c *gin.Context) {
			c.Next()
		}),
		timeout.WithResponse(func(c *gin.Context) {
			// 超时响应：HTTP 504 + 标准 JSON (业务码给 Error 500 即可，或者定义专门的 Timeout 码)
			response.Result(c, http.StatusGatewayTimeout, response.Error, "request timeout", nil)
		}),
	)
}