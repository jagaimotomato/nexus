package utils

import (
	"errors"
	"time"

	"nexus/internal/conf"
	"nexus/internal/data"

	"github.com/golang-jwt/jwt/v5"
)

// UserClaims 自定义 Token 载荷
type UserClaims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"` // 关键：放入角色 Key (如 ["admin", "editor"])
	jwt.RegisteredClaims
}

// GenerateToken 生成 Token
// 接收完整的 User 对象，以便提取 ID 和 Roles
func GenerateToken(user *data.User) (string, error) {
	cfg := conf.GlobalConfig.Jwt

	// 1. 提取用户角色 Key
	var roleKeys []string
	for _, role := range user.Roles {
		roleKeys = append(roleKeys, role.Key)
	}

	// 2. 构造 Claims
	claims := UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		Roles:    roleKeys,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.Expire) * time.Second)),
			Issuer:    cfg.Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 3. 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*UserClaims, error) {
	cfg := conf.GlobalConfig.Jwt

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}