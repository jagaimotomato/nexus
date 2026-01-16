package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CustomRecovery(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // 记录 Panic 堆栈到 Zap (error.log)
                logger.Error("Server Panic Recovered",
                    zap.Any("error", err),
                    zap.Stack("stack"),
                )
                c.AbortWithStatusJSON(500, gin.H{"error": "Internal Server Error"})
            }
        }()
        c.Next()
    }
}