// Package log is a global internal glogs
// glogs: this is extend package, use https://github.com/uber-go/zap
package glogs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"time"
)

const RequestIDName = "x-b3-traceid"

var (
	Sugar *zap.Logger
	//	ZapLog   *zap.Logger
	logLevel = zap.NewAtomicLevel()
)

// InitLog 初始化日志文件 logPath= /home/logs/app/appName
func InitLog(appName, logType, logPath string) {
	initSugar(appName, logPath, logType)
}

func initSugar(appName, logPath, logType string) {
	loglevel := zapcore.InfoLevel
	setLevel(loglevel)

	var core zapcore.Core
	// 打印至文件中
	if logType == "file" {
		configs := zap.NewProductionEncoderConfig()
		configs.FunctionKey = "func"
		configs.EncodeTime = timeEncoder
		//w := zapcore.AddSync(&lumberjack.Logger{
		//	Filename:   logPath, // 日志文件的位置
		//	MaxSize:    32,      // MB
		//	LocalTime:  true,    // 是否使用自己本地时间
		//	Compress:   false,   // 是否压缩/归档旧文件
		//	MaxAge:     90,      // 保留旧文件的最大天数
		//	MaxBackups: 300,     // 保留旧文件的最大个数
		//})

		w := zapcore.AddSync(getWriter(logPath))

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(configs),
			w,
			logLevel,
		)
		log.Printf("[glogs_sugar] log success")
	} else {
		// 打印在控制台
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), logLevel)
		log.Printf("[glogs_sugar] log success")
	}

	filed := zap.Fields(zap.String("app_name", appName))
	Sugar = zap.New(core, filed, zap.AddCaller(), zap.AddCallerSkip(1))
	//Sugar = logger.Sugar()
}

func setLevel(level zapcore.Level) { logLevel.SetLevel(level) }

func Info(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args)
	Sugar.Info(msg)
}
func InfoF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	msg := fmt.Sprintf(template, args)
	Sugar.Info(msg, zap.String("request_id", requestID), zap.String("title", title))
}

func Warn(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args)
	Sugar.Warn(msg)
}
func WarnF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	msg := fmt.Sprintf(template, args)
	Sugar.Info(msg, zap.String("request_id", requestID), zap.String("title", title))
}

func Error(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args)
	Sugar.Error(msg)
}
func ErrorF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	msg := fmt.Sprintf(template, args)
	Sugar.Info(msg, zap.String("request_id", requestID), zap.String("title", title))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	var layout = "2006-01-02 15:04:05"
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, layout)
		return
	}

	enc.AppendString(t.Format(layout))
}

// 按天切割
func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 stream-2021-5-20.log
	// demo.log是指向最新日志的连接
	// 保存7天内的日志，每1小时(整点)分割一第二天志
	hook, err := rotatelogs.New(
		filename+"/%Y-%m-%d.log",
		//rotatelogs.WithLinkName(filename+"/app.log"),
		//rotatelogs.WithMaxAge(time.Hour*24*365),
		rotatelogs.WithRotationTime(time.Hour*24),
	)

	if err != nil {
		panic(err)
	}
	return hook
}
