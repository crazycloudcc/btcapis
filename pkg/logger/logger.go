package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// 日志级别常量
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

var (
	// Log 全局日志实例
	Log *logrus.Logger
)

// InitLogger 初始化日志系统
func InitLogger(level string) {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 设置日志级别
	switch strings.ToLower(level) {
	case DebugLevel:
		Log.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		Log.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		Log.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}
}

// Debug 调试日志
func Debug(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

// Info 信息日志
func Info(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Warn 警告日志
func Warn(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error 错误日志
func Error(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}
