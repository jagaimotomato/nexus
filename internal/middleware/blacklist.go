package middleware

import (
	"net/http"
	"nexus/internal/response"
	"nexus/internal/service"

	"github.com/gin-gonic/gin"
)

func IPBlacklist(ipSecurity *service.IPSecurityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if ipSecurity.IsBlocked(clientIP) {
			response.Result(c, http.StatusForbidden, response.Forbidden, "Access Denied: Your IP is blocked.", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
