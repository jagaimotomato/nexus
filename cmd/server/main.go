package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Server.HttpPort),
			Handler: r,
		}

		// 3. 在 goroutine 中启动服务
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Log.Fatal("listen: ", zap.Error(err))
			}
		}()

		// 4. 监听中断信号
		quit := make(chan os.Signal, 1)
		// kill (无参) 默认发送 syscall.SIGTERM
		// kill -2 发送 syscall.SIGINT (Ctrl+C)
		// kill -9 发送 syscall.SIGKILL (无法捕获)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit // 阻塞直到收到信号
		logger.Log.Info("关闭服务...")

		// 5. 设置超时上下文，执行 Shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Log.Fatal("强制关闭服务: ", zap.Error(err))
		}

		logger.Log.Info("服务退出")
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
