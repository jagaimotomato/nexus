package data

import (
	"database/sql"
	"nexus/internal/conf"

	"context"
	"fmt"
	"nexus/internal/logger"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Data struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewData(cfg *conf.Config) (*Data, func(), error) {
	// 1. 初始化 MySQL
	db, err := newDB(cfg)
	if err != nil {
		return nil, nil, err
	}

	// 2. 初始化 Redis
	rdb, err := newRedis(cfg)
	if err != nil {
		return nil, nil, err
	}

	d := &Data{
		DB:  db,
		RDB: rdb,
	}

	// 返回清理函数，用于优雅停机
	cleanup := func() {
		logger.Log.Info("正在关闭数据层资源...")
		// 关闭 DB 连接
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			sqlDB.Close()
		}
		// 关闭 Redis 连接
		if rdb != nil {
			rdb.Close()
		}
	}

	return d, cleanup, nil
}

func newDB(cfg *conf.Config) (*gorm.DB, error) {
	ensureDatabase(cfg.Data.Database.Source)

	dsn := cfg.Data.Database.Source
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取 sql.DB 失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.Data.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Data.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Data.Database.ConnMaxLifetime) * time.Second)

	return db, nil
}

func newRedis(cfg *conf.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Data.Redis.Addr,
		Password: cfg.Data.Redis.Password,
		DB:       cfg.Data.Redis.DB,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("Redis连接失败: %w", err)
	}
	return rdb, nil
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
	if _, err := tempDB.Exec(createSQL); err != nil {
		logger.Log.Error("尝试创建数据库失败", zap.String("db_name", dbName), zap.Error(err))
	} else {
		logger.Log.Info("确保数据库存在", zap.String("db_name", dbName))
	}
}
