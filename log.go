package log

import (
	"os"
	"sync"

	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logInstance *zap.Logger
	once        sync.Once
)

// Logger return zap logger instance
func Logger() *zap.Logger {
	once.Do(func() {
		logInstance = logger()
	})

	return logInstance
}

func logger() *zap.Logger {
	// 日志滚动和备份策略
	infoWriter := lumberjack.Logger{
		Filename:   defaultViperString("log.info_file", "log/info.log"), // 日志输出地址
		LocalTime:  defaultViperBool("log.filename_with_time", true),    // 日志文件名时间
		MaxSize:    defaultViperInt("log.file_max_size", 100),           // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: defaultViperInt("log.file_max_backups", 30),         // 日志文件最多保存多少个备份
		MaxAge:     defaultViperInt("log.file_max_age", 30),             // 文件最多保存多少天
		Compress:   defaultViperBool("log.file_compress", true),         // 是否压缩
	}

	errorWriter := lumberjack.Logger{
		Filename:   defaultViperString(viper.GetString("log.error_file"), "log/error.log"), // 日志输出地址
		LocalTime:  defaultViperBool("log.filename_with_time", true),                       // 日志文件名时间
		MaxSize:    defaultViperInt("log.file_max_size", 100),                              // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: defaultViperInt("log.file_max_backups", 30),                            // 日志文件最多保存多少个备份
		MaxAge:     defaultViperInt("log.file_max_age", 30),                                // 文件最多保存多少天
		Compress:   defaultViperBool("log.file_compress", true),                            // 是否压缩
	}

	// 日志输出格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "",
		MessageKey:     "msg",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // FullCallerEncoder
		EncodeName:     zapcore.FullNameEncoder,
	}

	consoleConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "",
		CallerKey:      "",
		MessageKey:     "msg",
		StacktraceKey:  "",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 根据配置调整日志级别, 支持http接口动态修改zap日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	// 设置输出源，输出格式，日志等级
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&infoWriter), zap.InfoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&errorWriter), zap.ErrorLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(consoleConfig), zapcore.AddSync(os.Stdout), atomicLevel),
	)

	// // 开启开发模式，堆栈跟踪
	// caller := zap.AddCaller()
	// // trace
	// trace := zap.AddStacktrace(zap.InfoLevel)

	// 开启文件及行号
	development := zap.Development()

	// 构造日志
	z := zap.New(core, development)

	// TODO: add custom field
	z = z.With(zap.Int64("uuid", time.Now().Unix()))

	return z
}
