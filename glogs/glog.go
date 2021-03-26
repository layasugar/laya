// Package log is a global internal glogs
// glogs: this is extend package, use https://github.com/uber-go/zap
package glogs

import (
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya/gconf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

const RequestIDName = "request-id"

var (
	Sugar *zap.SugaredLogger
	//	ZapLog   *zap.Logger
	logLevel = zap.NewAtomicLevel()
)

// InitLog 初始化日志文件
func InitLog() {
	// 获取配置开启日志
	c := gconf.GetLogConf()
	if c.Open {
		InitSugar(c)
	}
}

func InitSugar(c gconf.LogConf) {
	loglevel := zapcore.InfoLevel
	switch c.LogLevel {
	case "INFO":
		loglevel = zapcore.InfoLevel
	case "ERROR":
		loglevel = zapcore.ErrorLevel
	}
	setLevel(loglevel)

	var core zapcore.Core
	// 打印至文件中
	if c.Driver == "file" {
		configs := zap.NewProductionEncoderConfig()
		configs.EncodeTime = zapcore.ISO8601TimeEncoder
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   c.Path, // 日志文件的位置
			MaxSize:    32,     // MB
			LocalTime:  true,   // 是否使用自己本地时间
			Compress:   true,   // 是否压缩/归档旧文件
			MaxAge:     90,     // 保留旧文件的最大天数
			MaxBackups: 300,    // 保留旧文件的最大个数
		})

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(configs),
			w,
			logLevel,
		)
	} else {
		// 打印在控制台
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), logLevel)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Sugar = logger.Sugar()
}

func setLevel(level zapcore.Level)                { logLevel.SetLevel(level) }
func Info(args ...interface{})                    { Sugar.Info(args...) }
func InfoF(template string, args ...interface{})  { Sugar.Infof(template, args...) }
func Warn(args ...interface{})                    { Sugar.Warn(args...) }
func WarnF(template string, args ...interface{})  { Sugar.Warnf(template, args...) }
func Error(args ...interface{})                   { Sugar.Error(args...) }
func ErrorF(template string, args ...interface{}) { Sugar.Errorf(template, args...) }

func InfoFR(c *gin.Context, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	template = "request_id=" + requestID + "," + template
	InfoF(template, args...)
}

func WarnFR(c *gin.Context, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	template = "request_id=" + requestID + "," + template
	WarnF(template, args...)
}

func ErrorFR(c *gin.Context, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	template = "request_id=" + requestID + "," + template
	ErrorF(template, args...)
}
