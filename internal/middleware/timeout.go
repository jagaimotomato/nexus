package middleware

import (
	"net/http"
	"time"

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
			// 超时后的响应
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timeout"})
		}),
	)
}