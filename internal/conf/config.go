package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server Server `mapstructure:"server"`
	Data   Data   `mapstructure:"data"`
	Jwt    Jwt    `mapstructure:"jwt"`
	Log    Log    `mapstructure:"log"`
	RateLimit RateLimit `mapstructure:"ratelimit"`
}

type Server struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	HttpPort     int    `mapstructure:"http_port"`
	GrpcPort     int    `mapstructure:"grpc_port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

type Data struct {
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
}

type Database struct {
	Driver          string `mapstructure:"driver"`
	Source          string `mapstructure:"source"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type Redis struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type Jwt struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
	Issuer string `mapstructure:"issuer"`
}

type Log struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Director   string `mapstructure:"director"`
	LinkName   string `mapstructure:"link_name"`
	ShowLine   bool   `mapstructure:"show_line"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

type RateLimit struct {
	Qps int `mapstructure:"qps"`
	Burst int `mapstructure:"burst"`
}

var GlobalConfig Config

func InitConfig() *Config {
	// 告诉 Viper 配置文件在哪里、叫什么
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") // 文件类型
	viper.AddConfigPath("./configs") // 配置文件路径
	viper.AddConfigPath("../../configs") // 配置文件路径

	if err := viper.ReadInConfig(); err != nil {
		// 如果是“文件未找到”错误，可以根据情况处理，否则直接 panic
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            panic("未找到配置文件 config.yaml，请检查路径")
        } else {
            panic(fmt.Errorf("读取配置文件出错: %w", err))
        }
	}

	// 将读取到的内容解析到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		panic(fmt.Errorf("配置解析失败: %w", err))
	}

	fmt.Println("✅ 配置加载成功:", viper.ConfigFileUsed())

	return &GlobalConfig
}