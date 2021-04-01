// Package log is a global internal glogs
// glogs: this is extend package, use https://github.com/uber-go/zap
package glogs

import (
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya/gconf"
	"github.com/layatips/laya/genv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
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
	InitSugar(c)
}

func InitSugar(c gconf.LogConf) {
	runMode := gconf.GetRunMode()
	loglevel := zapcore.InfoLevel
	setLevel(loglevel)

	var core zapcore.Core
	// 打印至文件中
	if runMode == "release" {
		configs := zap.NewProductionEncoderConfig()
		configs.EncodeTime = zapcore.ISO8601TimeEncoder
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   c.Path + genv.AppName() + "/app.log", // 日志文件的位置
			MaxSize:    c.MaxSize,                            // MB
			LocalTime:  true,                                 // 是否使用自己本地时间
			Compress:   true,                                 // 是否压缩/归档旧文件
			MaxAge:     c.MaxAge,                             // 保留旧文件的最大天数
			MaxBackups: c.MaxBackups,                         // 保留旧文件的最大个数
		})

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(configs),
			w,
			logLevel,
		)
		log.Printf("[glogs_sugar] logs open success at %s\n", c.Path+genv.AppName()+"/app.log")
	} else {
		// 打印在控制台
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), logLevel)
		log.Printf("[glogs_sugar] logs open success at console\n")
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
