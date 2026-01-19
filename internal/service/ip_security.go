package service

import (
	"nexus/internal/logger"
	"sync"
	"time"

	"nexus/internal/data"

	"go.uber.org/zap"
)

// 内存缓存系统
type ipCacheSystem struct {
   blockedIPs sync.Map // 结构 IP-> 过期时间（time.Time)
}

// 单例实例
var IPCache = &ipCacheSystem{}

// 启动加载
func LoadBlacklistAtStartup() {
	var list []data.Blacklist
	// 查询未过期的记录
	err := data.DB.Where("is_active = ? AND (expires_at > ? OR expires_at IS NULL)", true, time.Now()).Find(&list).Error
	if err != nil {
    logger.Log.Error("加载黑名单失败", zap.Error(err))
	return
	}
	count := 0
	for _, item := range list {
		expire := time.Now().Add(100 *365 *24 * time.Hour)
		if item.ExpiresAt != nil {
			expire = *item.ExpiresAt
		}
		IPCache.blockedIPs.Store(item.IP, expire)
		count++
	}
	logger.Log.Info("黑名单加载完成", zap.Int("count", count))
}

// IsBlocked 检查 IP 是否被封禁 (中间件调用，极速)
func IsBlocked(ip string) bool {
	val, ok := IPCache.blockedIPs.Load(ip)
	if !ok {
		return false
	}

	expireTime := val.(time.Time)
	// 检查是否过期
	if time.Now().After(expireTime) {
		IPCache.blockedIPs.Delete(ip) // 懒惰删除
		return false
	}
	return true
}

// BanIP 执行封禁 (自动封禁/手动封禁通用)
func BanIP(ip string, reason string, duration time.Duration) {
	// A. 先更新内存 (立即生效)
	expire := time.Now().Add(duration)
	IPCache.blockedIPs.Store(ip, expire)

	// B. 异步写入数据库 (持久化)
	go func() {
		newBan := data.Blacklist{
			IP:        ip,
			Reason:    reason,
			ExpiresAt: &expire,
			IsActive:  true,
		}

		// Upsert: 如果存在则更新过期时间
		result := data.DB.Where(data.Blacklist{IP: ip}).
			Assign(data.Blacklist{ExpiresAt: &expire, Reason: reason}).
			FirstOrCreate(&newBan)

		if result.Error != nil {
			logger.Log.Error("封禁IP落库失败", zap.String("ip", ip), zap.Error(result.Error))
		} else {
			logger.Log.Warn("已封禁IP", zap.String("ip", ip), zap.String("reason", reason))
		}
	}()
}