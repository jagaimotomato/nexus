package utils

import (
	"context"
	"nexus/internal/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisStore 实现base64Captcha.Store接口
type RedisStore struct {
	Expiration time.Duration
	Client     *redis.Client
}

// Set 将验证码存入 redis
func (s *RedisStore) Set(id string, value string) error {
	ctx := context.Background()
	key := "captcha:" + id
	if s.Client == nil {
		return nil
	}
	// 存入redis 有效期5分钟
	err := s.Client.Set(ctx, key, value, s.Expiration).Err()
	if err != nil {
		logger.Log.Error("设置验证码失败", zap.Error(err))
		return err
	}
	return nil
}

// get 从redis 获取验证码
func (s *RedisStore) Get(id string, clear bool) string {
	ctx := context.Background()
	key := "captcha:" + id
	if s.Client == nil {
		return ""
	}
	// 从redis 获取验证码
	value, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		logger.Log.Error("获取验证码失败", zap.Error(err))
		return ""
	}
	return value
}

// verify 验证验证码
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	v := s.Get(id, clear)
	return v == answer
}
