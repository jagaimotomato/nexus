package middleware

import (
	"net/http"
	"sync"
	"time"

	"nexus/internal/service" // 引用 Service 层

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// 简单的违规计数器 (生产环境建议用 Redis)
var (
	violationStore sync.Map
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
			count := incrementViolation(ip)

			// 3. 判断是否触发自动封禁
			// 规则: 如果短时间内违规超过 20 次，封禁 24 小时
			if count > 20 {
				service.BanIP(ip, "Auto Ban: Rate Limit Abuse", 24*time.Hour)
				
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "IP Blocked due to abusive behavior",
				})
				return
			}

			// 4. 普通限流提示
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, slow down",
			})
			return
		}

		c.Next()
	}
}

// incrementViolation 增加违规次数 (简单的内存实现，带简单的重置机制)
// 实际生产中这里应该用 Redis.Incr 并设置 1 分钟过期
func incrementViolation(ip string) int {
	v, _ := violationStore.LoadOrStore(ip, 0)
	count := v.(int) + 1
	violationStore.Store(ip, count)

	// 简单清理逻辑: 启动一个协程，1分钟后把这个 IP 的违规数清零
	// (注意: 这只是极简演示，防止内存泄漏。严谨做法是用 Redis)
	if count == 1 {
		time.AfterFunc(1*time.Minute, func() {
			violationStore.Delete(ip)
		})
	}
	return count
}