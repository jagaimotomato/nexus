package main

import (
	"context"
	"fmt"
	"net/http"
	"nexus/internal/conf"
	"nexus/internal/data"
	"nexus/internal/handler"
	"nexus/internal/logger"
	"nexus/internal/router"
	"nexus/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "nexus",
		Short: "Nexus Server",
		Run: func(cmd *cobra.Command, args []string) {
			runServer()
		},
	}
}

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runServer() {
	cfg := conf.InitConfig()
	logger.InitLogger(&cfg.Log)
	defer logger.Log.Sync()

	logger.Log.Info("服务启动中...", zap.Any("config", cfg.Server))

	d, cleanup, err := data.NewData(cfg)
	if err != nil {
		logger.Log.Fatal("数据层初始化失败", zap.Error(err))
	}
	defer cleanup()

	data.Migrate(d.DB)

	authService := service.NewAuthService(d, cfg)
	menuService := service.NewMenuService(d)
	ipSecurityService := service.NewIPSecurityService(d)

	authHandler := handler.NewAuthHandler(authService)
	menuHandler := handler.NewMenuHandler(menuService)

	service.InitData(data.NewUserRepo(d.DB), data.NewRoleRepo(d.DB))
	ipSecurityService.LoadBlacklistAtStartup()

	r := router.NewRouter(cfg, authHandler, menuHandler, ipSecurityService)

	addr := fmt.Sprintf(":%d", cfg.Server.HttpPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		logger.Log.Info("HTTP服务启动成功", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("服务异常退出", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("正在关闭服务...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("服务强制关闭", zap.Error(err))
	}

	logger.Log.Info("服务已安全停止")
}
