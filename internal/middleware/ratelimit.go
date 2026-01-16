package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter 存储每个 IP 的限流器
type IPRateLimiter struct {
	ips sync.Map
	r   rate.Limit // 每秒产生多少令牌 (QPS)
	b   int        // 桶的大小 (突发并发量)
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		r: r,
		b: b,
	}
}

// GetLimiter 获取或创建一个针对特定 IP 的限流器
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	limiter, exists := i.ips.Load(ip)
	if !exists {
		// 创建一个新的限流器
		newLimiter := rate.NewLimiter(i.r, i.b)
		// 使用 LoadOrStore 防止并发创建覆盖
		actual, _ := i.ips.LoadOrStore(ip, newLimiter)
		return actual.(*rate.Limiter)
	}
	return limiter.(*rate.Limiter)
}

// RateLimit 中间件入口
// qps: 每秒允许的请求数
// burst: 允许瞬间突发的最大请求数
func RateLimit(qps int, burst int) gin.HandlerFunc {
	limiter := NewIPRateLimiter(rate.Limit(qps), burst)

	// 定期清理过期 IP (防止 map 无限增长内存泄漏) - 简易版略过，生产环境需加上

	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)

		// Allow() 非阻塞，如果桶里没令牌了，直接返回 false
		if !l.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "request too frequent",
			})
			return
		}

		c.Next()
	}
}