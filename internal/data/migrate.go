package data

import (
	"nexus/internal/logger"

	"go.uber.org/zap"
)

func Migrate() {
	logger.Log.Info("开始迁移数据库表结构...")

	err := DB.AutoMigrate(&User{}, &Role{}, &Blacklist{})
	if err != nil {
		logger.Log.Fatal("迁移数据库表结构失败", zap.Error(err))
	}

	logger.Log.Info("数据库表结构迁移完成")
}