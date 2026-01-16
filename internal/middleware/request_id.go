package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-Id"
const ContextKeyRequestID = "request_id"

// RequestID 为每个请求生成唯一的 UUID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 尝试从请求头获取（用于串联上下游服务，比如 Nginx -> Master）
		reqID := c.GetHeader(HeaderXRequestID)

		// 2. 如果没有，则生成一个新的
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// 3. 写入 Context，供后续业务逻辑/日志使用
		c.Set(ContextKeyRequestID, reqID)

		// 4. 写入响应头，方便前端/客户排查问题
		c.Writer.Header().Set(HeaderXRequestID, reqID)

		c.Next()
	}
}