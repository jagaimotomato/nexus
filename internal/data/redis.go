package data

import (
	"context"
	"fmt"
	"nexus/internal/conf"
	"nexus/internal/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RDB *redis.Client

func InitRedis() {
	cfg := conf.GlobalConfig.Data.Redis

	RDB = redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
		Password: cfg.Password,
		DB: cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		logger.Log.Fatal("Redis 连接失败", zap.Error(err), zap.String("addr", cfg.Addr))
		panic(fmt.Errorf("redis connect failed %w", err))
	}
	logger.Log.Info("Redis 连接成功", zap.String("addr", cfg.Addr))
}