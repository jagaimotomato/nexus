package main

import (
	"fmt"
	"os"

	"nexus/internal/conf"
	"nexus/internal/data"
	"nexus/internal/logger"
	"nexus/internal/router"
	"nexus/internal/service"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// 定义全局变量，用于接收 --config 参数
var cfgFile string

// 1. 定义 Root Command (主命令)
var rootCmd = &cobra.Command{
	Use:   "nexus",
	Short: "Nexus Server",
	// 核心运行逻辑
	Run: func(cmd *cobra.Command, args []string) {
		cfg := conf.InitConfig()
		fmt.Printf("服务启动成功，当前配置: %+v\n", cfg)
		logger.InitLogger(&cfg.Log)
		// 确保程序退出前日志刷盘
		defer logger.Log.Sync()

		logger.Log.Info("服务启动中...", zap.Any("config", cfg.Server))

		// 初始化数据层
		data.InitDB()
		data.InitRedis()

		data.Migrate()
		service.LoadBlacklistAtStartup()
		service.InitData()

		r := router.InitRouter(cfg)

		addr := fmt.Sprintf(":%d", cfg.Server.HttpPort)
		logger.Log.Info("服务启动成功", zap.String("addr", addr))

		if err := r.Run(addr); err != nil {
			logger.Log.Fatal("服务启动失败", zap.Error(err))
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initTables() {
	logger.Log.Info("初始化数据库表结构...")

	err := data.DB.AutoMigrate(&data.User{}, &data.Role{}, &data.Blacklist{})
	if err != nil {
		logger.Log.Fatal("初始化数据库表结构失败", zap.Error(err))
	}

	logger.Log.Info("数据库表结构初始化完成")
}
