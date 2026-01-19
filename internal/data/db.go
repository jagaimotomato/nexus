package data

import (
	"fmt"
	"nexus/internal/conf"
	"nexus/internal/logger"

	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	c := conf.GlobalConfig.Data.Database
	dsn := c.Source

	gormConfig := &gorm.Config{
		// 命名策略: 把结构体 Blacklist 映射为表名 blacklists (蛇形复数)
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		logger.Log.Fatal("MYSQL 数据库连接失败", zap.Error(err))
		panic(fmt.Errorf("MYSQL 数据库连接失败: %w", err))
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Log.Fatal("获取 sql.DB失败", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second)

	logger.Log.Info("MySQL 连接成功", zap.String("dsn", maskPassword(dsn)))
}

func maskPassword(dsn string) string {
	// 简单实现，实际可以正则替换
	if len(dsn) > 10 {
		return dsn[:5] + "***" + dsn[len(dsn)-5:]
	}
	return "***"
}