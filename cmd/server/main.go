package main

import (
	"fmt"
	"os"

	"nexus/internal/conf"
	"nexus/internal/logger"

	"github.com/spf13/cobra"
)

// 定义全局变量，用于接收 --config 参数
var cfgFile string

// 1. 定义 Root Command (主命令)
var rootCmd = &cobra.Command{
	Use: "nexus",
	Short: "Nexus Server",
	// 核心运行逻辑
    Run: func (cmd *cobra.Command, args []string){
		cfg := conf.InitConfig()
		fmt.Printf("服务启动成功，当前配置: %+v\n", cfg)
		logger.InitLogger(&cfg.Log)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
