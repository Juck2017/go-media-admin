package utils

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

/**
	param:
		levelString, 日志级别文字
		logFilePath, 日志文件路径
		isDevMode, 是否是开发者模式
 */
func InitLogger(levelString string, logFilePath string, isDevMode bool) *zap.SugaredLogger {
	// 设置日志级别
	var level zapcore.Level
	switch levelString {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}
	encoder := getEncoder() // 获取编码后日志
	writeSyncer := getLogWriter(logFilePath, isDevMode) // 判断是否需要进行切割日志操作,并进行日志切割
	core := zapcore.NewCore(encoder, writeSyncer, level)
	// zap.AddCaller()  添加将调用函数信息记录到日志中的功能。
	return zap.New(core, zap.AddCaller()).Sugar()
}

// 获取编码后日志
func getEncoder() zapcore.Encoder {
	// encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间编码器

	// 在日志文件中使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// NewConsoleEncoder 打印更符合人们观察的方式
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 判断是否需要进行切割日志操作,并进行日志切割
func getLogWriter(logFilePath string, isDevMode bool) zapcore.WriteSyncer {
	var logger io.Writer
	if isDevMode {
		// 开发模式下，日志直接输出到控制台，方便开发者查看和调试
		logger = os.Stdout
	} else {
		// 使用Lumberjack进行日志切割归档
		logger = &lumberjack.Logger{
			Filename:   logFilePath, // 日志文件的位置
			MaxSize:    10, // 在进行切割之前，日志文件的最大大小（以MB为单位）
			MaxBackups: 5, // 保留旧文件的最大个数
			MaxAge:     30, // 保留旧文件的最大天数
			Compress:   false, // 是否压缩/归档旧文件
		}
	}
	return zapcore.AddSync(logger)
}
