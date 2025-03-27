package initialize

import (
	"chat/internal/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

// 日志组件,可以复用
// SetupLogger 初始化日志组件，设置日志级别、编码方式和输出目标
func SetupLogger() {
	// 创建一个可原子性修改的日志级别
	level := zap.NewAtomicLevel()
	level.SetLevel(zapcore.DebugLevel)

	// 创建控制台日志编码器，配置日志消息的格式
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:       "message",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "logger",
		CallerKey:        "caller",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       CustomTimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
		ConsoleSeparator: "",
	})

	// 定义两个日志输出核心，一个输出到标准输出，一个输出到文件
	cores := [...]zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, level),
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(getwritesync()),
			level,
		),
	}

	// 创建并初始化全局日志对象
	global.Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()

	}(global.Logger)

	// 记录日志组件初始化成功的信息
	global.Logger.Info("日志组件初始化成功")
}

// getwritesync 返回一个实现了zapcore.WriteSyncer接口的lumberjack.Logger对象，
// 用于日志文件的滚动和备份
func getwritesync() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   global.Config.ZapConfig.Filename,
		MaxSize:    global.Config.ZapConfig.MaxSize,
		MaxBackups: global.Config.ZapConfig.MaxBackups,
		MaxAge:     global.Config.ZapConfig.MaxAge,
		LocalTime:  true,
	}

	return zapcore.AddSync(lumberJackLogger)
}

// CustomTimeEncoder 自定义时间编码器，将时间格式化为"2006-01-02 15:04:05"的形式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
