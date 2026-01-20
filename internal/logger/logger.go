package logger

import (
	"nexus/internal/conf"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger // 给业务逻辑用
	AccessLog *zap.Logger // 给 Gin 中间件用
)

// InitLogger 初始化入口
func InitLogger(cfg *conf.Log) {
	// 1. 确保日志根目录存在 (例如 ./logs)
    // Lumberjack 虽然会尝试创建，但显式创建更稳健
    if ok, _ := pathExists(cfg.Director); !ok {
        _ = os.MkdirAll(cfg.Director, os.ModePerm)
    }
	// 1. 公用的编码器配置 (JSON)
	encoder := getEncoder(cfg.Format)

	// --------------------------------------------------------
	// 2. 初始化业务 Logger (app.log)
	// --------------------------------------------------------
	// 负责写入文件 (app.log)
	appWriter := getLogWriter(cfg, "app.log")
	
	// 同时输出到控制台和文件
	appSyncer := zapcore.NewMultiWriteSyncer(appWriter, zapcore.AddSync(os.Stdout))

	// 解析日志级别
	level := getLogLevel(cfg.Level)

	// 创建业务 Core
	appCore := zapcore.NewCore(encoder, appSyncer, level)
	
	// 生成业务 Logger (开启 Caller 显示行号)
	Log = zap.New(appCore, zap.AddCaller())
	zap.ReplaceGlobals(Log) // 替换全局方便使用

	// --------------------------------------------------------
	// 3. 初始化访问 Logger (access.log)
	// --------------------------------------------------------
	// 负责写入文件 (access.log)
	accessWriter := getLogWriter(cfg, "access.log")
	
	// AccessLog 通常不需要打印到控制台，也不需要行号(Caller)，只写文件即可
	// 级别固定为 Info
	accessCore := zapcore.NewCore(encoder, accessWriter, zap.InfoLevel)
	
	AccessLog = zap.New(accessCore) // 注意：这里没有 AddCaller，因为是在中间件里调用的，行号没意义
}

// ---- 下面是辅助函数，避免代码重复 ----

// 获取编码器 (JSON/Console)
func getEncoder(format string) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	
	if format == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 获取日志写入器 (Lumberjack)
func getLogWriter(cfg *conf.Log, fileName string) zapcore.WriteSyncer {
	lumberLogger := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Director, fileName), // 拼接路径
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.AddSync(lumberLogger)
}

// 获取日志级别
func getLogLevel(levelStr string) zapcore.Level {
	switch strings.ToLower(levelStr) {
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

// 辅助函数: 判断路径是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}