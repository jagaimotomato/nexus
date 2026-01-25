package data

import (
	"database/sql"
	"fmt"
	"nexus/internal/conf"
	"nexus/internal/logger"
	"strings"

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

	// 1. 在连接 GORM 之前，确保数据库存在
	ensureDatabase(dsn)

	gormConfig := &gorm.Config{
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

func ensureDatabase(dsn string) {
	startIndex := strings.LastIndex(dsn, ")/")
	if startIndex == -1 {
		logger.Log.Warn("DSN 格式无法自动解析，跳过自动创建数据库步骤")
		return
	}
	dbNameStart := startIndex + 2

	endIndex := strings.Index(dsn[dbNameStart:], "?")
	var dbName string
	if endIndex == -1 {
		dbName = dsn[dbNameStart:]
	} else {
		dbName = dsn[dbNameStart : dbNameStart+endIndex]
	}

	if dbName == "" {
		return
	}

	baseDSN := dsn[:startIndex+2]
	tempDB, err := sql.Open("mysql", baseDSN)
	if err != nil {
		logger.Log.Warn("无法建立临时连接，跳过创建数据库", zap.Error(err))
		return
	}
	defer tempDB.Close()

	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci", dbName)
	_, err = tempDB.Exec(createSQL)
	if err != nil {
		logger.Log.Error("尝试创建数据库失败", zap.String("db_name", dbName), zap.Error(err))
		// 这里不 Panic，让后续 GORM 连接去报错，可能用户权限不足但库已存在
	} else {
		logger.Log.Info("确保数据库存在", zap.String("db_name", dbName))
	}
}

func maskPassword(dsn string) string {
	// 简单实现，实际可以正则替换
	if len(dsn) > 10 {
		return dsn[:5] + "***" + dsn[len(dsn)-5:]
	}
	return "***"
}
