package middleware

import (
	"context"
	"net/http"
	"nexus/internal/data" // 引入 data 包使用 RDB
	"nexus/internal/response"
	"nexus/internal/service"
	"sync"
	"time"

	"nexus/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// IPRateLimiter 包装标准库的限流器
type IPRateLimiter struct {
	ips sync.Map
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{r: r, b: b}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	limiter, exists := i.ips.Load(ip)
	if !exists {
		newLimiter := rate.NewLimiter(i.r, i.b)
		actual, _ := i.ips.LoadOrStore(ip, newLimiter)
		return actual.(*rate.Limiter)
	}
	return limiter.(*rate.Limiter)
}

// RateLimitWithAutoBan 限流并自动封禁
// qps: 每秒允许请求数
// burst: 允许瞬间爆发数
func RateLimitWithAutoBan(qps int, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(qps), burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		// 1. 如果限流器不允许通过 (令牌用光了)
		if !l.Allow() {
			// 2. 增加违规计数
			count := incrementViolationRedis(ip)

			if count > 20 {
				logger.Log.Warn("触发自动封禁", zap.String("ip", ip), zap.Int("violation_count", count))
				
				service.BanIP(ip, "Auto Ban: Rate Limit Abuse", 24*time.Hour)
				
				// 封禁响应：HTTP 403 + 标准 JSON
				response.Result(c, http.StatusForbidden, response.Forbidden, "IP Blocked due to abusive behavior", nil)
				c.Abort()
				return
			}

			// 限流响应：HTTP 429 + 标准 JSON
			response.Result(c, http.StatusTooManyRequests, response.TooManyRequests, "Too many requests, slow down", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func incrementViolationRedis(ip string) int {
	ctx := context.Background()
	key := "violation:" + ip

	// 原子递增
	count, err := data.RDB.Incr(ctx, key).Result()
	if err != nil {
		logger.Log.Error("Redis Incr error", zap.Error(err))
		return 0
	}

	// 设置过期时间
	if count == 1 {
		data.RDB.Expire(ctx, key, 1*time.Minute)
	}

	return int(count)
}