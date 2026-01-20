package middleware

import (
	"net/http"
	"nexus/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CustomRecovery(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Server Panic Recovered",
					zap.Any("error", err),
					zap.Stack("stack"),
				)
				// 崩溃恢复：HTTP 500 + 标准 JSON
				response.Result(c, http.StatusInternalServerError, response.Error, "Internal Server Error", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}