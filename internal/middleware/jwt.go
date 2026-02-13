package middleware

import (
	"strings"

	"nexus/internal/conf"
	"nexus/internal/response"
	"nexus/internal/utils"

	"github.com/gin-gonic/gin"
)

// 上下文 Key 定义，防止硬编码字符串拼错
const (
	CtxUserID   = "userID"
	CtxUsername = "username"
	CtxRoles    = "roles"
)

func JWTAuth(jwtCfg conf.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.FailWithCode(c, response.AuthFailed)
			c.Abort()
			return
		}

		// 2. 校验 Bearer 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.FailWithDetailed(c, response.AuthFailed, "Token 格式错误")
			c.Abort()
			return
		}

		// 3. 解析 Token
		claims, err := utils.ParseToken(parts[1], utils.JWTConfig{
			Secret: jwtCfg.Secret,
			Expire: jwtCfg.Expire,
			Issuer: jwtCfg.Issuer,
		})
		if err != nil {
			response.FailWithDetailed(c, response.AuthFailed, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 4. 将用户信息存入上下文 (最关键的一步)
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Set(CtxRoles, claims.Roles) // 后续可以用来做角色鉴权

		c.Next()
	}
}
