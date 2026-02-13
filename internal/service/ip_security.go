package service

import (
	"context"
	"errors"
	"nexus/internal/data"
	"nexus/internal/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type IPSecurityService struct {
	db         *gorm.DB
	rdb        *redis.Client
	blockedIPs sync.Map
}

func NewIPSecurityService(d *data.Data) *IPSecurityService {
	return &IPSecurityService{
		db:  d.DB,
		rdb: d.RDB,
	}
}

func (s *IPSecurityService) LoadBlacklistAtStartup() {
	var list []data.Blacklist
	err := s.db.Where("is_active = ? AND (expires_at > ? OR expires_at IS NULL)", true, time.Now()).Find(&list).Error
	if err != nil {
		logger.Log.Error("加载黑名单失败", zap.Error(err))
		return
	}

	count := 0
	for _, item := range list {
		expire := time.Now().Add(100 * 365 * 24 * time.Hour)
		if item.ExpiresAt != nil {
			expire = *item.ExpiresAt
		}
		s.blockedIPs.Store(item.IP, expire)
		count++
	}

	logger.Log.Info("黑名单加载完成", zap.Int("count", count))
}

func (s *IPSecurityService) IsBlocked(ip string) bool {
	val, ok := s.blockedIPs.Load(ip)
	if !ok {
		return false
	}

	expireTime := val.(time.Time)
	if time.Now().After(expireTime) {
		s.blockedIPs.Delete(ip)
		return false
	}

	return true
}

func (s *IPSecurityService) BanIP(ip string, reason string, duration time.Duration) {
	expire := time.Now().Add(duration)
	s.blockedIPs.Store(ip, expire)

	go func() {
		newBan := data.Blacklist{
			IP:        ip,
			Reason:    reason,
			ExpiresAt: &expire,
			IsActive:  true,
		}

		result := s.db.Where(data.Blacklist{IP: ip}).
			Assign(data.Blacklist{ExpiresAt: &expire, Reason: reason}).
			FirstOrCreate(&newBan)

		if result.Error != nil {
			logger.Log.Error("封禁IP落库失败", zap.String("ip", ip), zap.Error(result.Error))
		} else {
			logger.Log.Warn("已封禁IP", zap.String("ip", ip), zap.String("reason", reason))
		}
	}()
}

func (s *IPSecurityService) IncrementViolation(ip string, ttl time.Duration) (int, error) {
	if s.rdb == nil {
		return 0, errors.New("redis client is nil")
	}

	ctx := context.Background()
	key := "violation:" + ip

	count, err := s.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		if err := s.rdb.Expire(ctx, key, ttl).Err(); err != nil {
			return int(count), err
		}
	}

	return int(count), nil
}
