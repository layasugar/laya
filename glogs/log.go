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

const (
	RequestIDName = "x-b3-traceid"
	HeaderAppName = "app-name"
)

var (
	Sugar    *zap.Logger
	logLevel = zap.NewAtomicLevel()
	//	ZapLog   *zap.Logger
	lc = &LogConfig{
		appName:      "default-app",        // 默认应用名称
		appMode:      "dev",                // 默认应用环境
		logType:      "file",               // 默认日志类型
		logPath:      "/home/logs/app",     // 默认文件目录
		childPath:    "glogs/%Y-%m-%d.log", // 默认子目录
		rotationSize: 32 * 1024 * 1024,     // 默认大小为32M
	}
)

type LogConfig struct {
	appName       string        // 应用名
	appMode       string        // 应用环境
	logType       string        // 日志类型
	logPath       string        // 日志主路径
	childPath     string        // 日志子路径+文件名
	rotationSize  int64         // 单个文件大小
	rotationCount uint          // 可以保留的文件个数
	rotationTime  time.Duration // 日志分割的时间
	maxAge        time.Duration // 日志最大保留的天数
}

type LogOptionFunc func(*LogConfig)

type CusLog struct {
	Logger *zap.Logger
	Config *LogConfig
}

// InitLog 初始化日志文件 logPath= /home/logs/app/appName/childPath
func InitLog(options ...LogOptionFunc) {
	for _, f := range options {
		f(lc)
	}

	Sugar = initSugar(lc)
}

// NewLogger 得到一个zap.Logger
func NewLogger(options ...LogOptionFunc) *CusLog {
	var cus = &CusLog{Config: lc}
	for _, f := range options {
		f(cus.Config)
	}

	cus.Logger = initSugar(cus.Config)
	return cus
}

func initSugar(lc *LogConfig) *zap.Logger {
	loglevel := zapcore.InfoLevel
	setLevel(loglevel)

	logPath := fmt.Sprintf("%s/%s/%s", lc.logPath, lc.appName, lc.childPath)

	var core zapcore.Core
	// 打印至文件中
	if lc.logType == "file" {
		configs := zap.NewProductionEncoderConfig()
		configs.FunctionKey = "func"
		configs.EncodeTime = timeEncoder
		w := zapcore.AddSync(GetWriter(
			logPath,
			rotatelogs.WithRotationSize(lc.rotationSize),
			rotatelogs.WithRotationCount(lc.rotationCount),
			rotatelogs.WithRotationTime(lc.rotationTime),
			rotatelogs.WithMaxAge(lc.maxAge),
		))

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

	filed := zap.Fields(zap.String("app_name", lc.appName), zap.String("app_mode", lc.appMode))
	return zap.New(core, filed, zap.AddCaller(), zap.AddCallerSkip(1))
	//Sugar = logger.Sugar()
}

func setLevel(level zapcore.Level) { logLevel.SetLevel(level) }

func Info(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	Sugar.Info(msg)
}
func InfoF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	originAppName := c.GetHeader(HeaderAppName)
	path := c.Request.RequestURI
	msg := fmt.Sprintf(template, args...)
	Sugar.Info(msg,
		zap.String("path", path),
		zap.String("request_id", requestID),
		zap.String("title", title),
		zap.String("origin_app_name", originAppName))
}

func Warn(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	Sugar.Warn(msg)
}
func WarnF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	originAppName := c.GetHeader(HeaderAppName)
	path := c.Request.RequestURI
	msg := fmt.Sprintf(template, args...)
	Sugar.Info(msg,
		zap.String("path", path),
		zap.String("request_id", requestID),
		zap.String("title", title),
		zap.String("origin_app_name", originAppName))
}

func Error(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	Sugar.Error(msg)
}
func ErrorF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	originAppName := c.GetHeader(HeaderAppName)
	path := c.Request.RequestURI
	msg := fmt.Sprintf(template, args...)
	Sugar.Info(msg,
		zap.String("path", path),
		zap.String("request_id", requestID),
		zap.String("title", title),
		zap.String("origin_app_name", originAppName))
}

func (l *CusLog) Info(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Info(msg)
}
func (l *CusLog) InfoF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	originAppName := c.GetHeader(HeaderAppName)
	path := c.Request.RequestURI
	msg := fmt.Sprintf(template, args...)
	l.Logger.Info(msg,
		zap.String("path", path),
		zap.String("request_id", requestID),
		zap.String("title", title),
		zap.String("origin_app_name", originAppName))
}

func (l *CusLog) Warn(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Warn(msg)
}
func (l *CusLog) WarnF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	originAppName := c.GetHeader(HeaderAppName)
	path := c.Request.RequestURI
	msg := fmt.Sprintf(template, args...)
	l.Logger.Info(msg,
		zap.String("path", path),
		zap.String("request_id", requestID),
		zap.String("title", title),
		zap.String("origin_app_name", originAppName))
}

func (l *CusLog) Error(template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	l.Logger.Error(msg)
}
func (l *CusLog) ErrorF(c *gin.Context, title string, template string, args ...interface{}) {
	requestID := c.GetHeader(RequestIDName)
	originAppName := c.GetHeader(HeaderAppName)
	path := c.Request.RequestURI
	msg := fmt.Sprintf(template, args...)
	l.Logger.Info(msg,
		zap.String("path", path),
		zap.String("request_id", requestID),
		zap.String("title", title),
		zap.String("origin_app_name", originAppName))
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

// 按天切割按大小切割
// filename 文件名
// rotationSize 每个文件的大小
// maxAge 文件最大保留天数
// rotationCount 最大保留文件个数
// rotationTime 设置文件分割时间
// rotationCount 设置保留的最大文件数量
func GetWriter(filename string, options ...rotatelogs.Option) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 stream-2021-5-20.log
	// demo.log是指向最新日志的连接
	// 保存7天内的日志，每1小时(整点)分割一第二天志
	hook, err := rotatelogs.New(
		filename,
		options...,
	)

	if err != nil {
		panic(err)
	}
	return hook
}

// 设置应用名称,默认值default-app
func SetLogAppName(appName string) LogOptionFunc {
	return func(c *LogConfig) {
		if appName != "" {
			c.appName = appName
		}
	}
}

// 设置环境变量,标识当前应用运行的环境,默认值dev
func SetLogAppMode(appMode string) LogOptionFunc {
	return func(c *LogConfig) {
		if appMode != "" {
			c.appMode = appMode
		}
	}
}

// 设置日志类型,日志类型目前分为2种,console和file,默认值file
func SetLogType(logType string) LogOptionFunc {
	return func(c *LogConfig) {
		if logType != "" {
			c.logType = logType
		}
	}
}

// 设置日志目录,这个是主目录,程序会给此目录拼接上项目名,子目录以及文件,默认值/home/logs/app
func SetLogPath(logPath string) LogOptionFunc {
	return func(c *LogConfig) {
		if logPath != "" {
			c.logPath = logPath
		}
	}
}

// 设置子目录—+文件名,保证一个类型的文件在同一个文件夹下面便于区分、默认值glogs/%Y-%m-%d.log
func SetLogChildPath(childPath string) LogOptionFunc {
	return func(c *LogConfig) {
		if childPath != "" {
			c.childPath = childPath
		}
	}
}

// 设置单个文件最大值byte,默认值32M
func SetLogMaxSize(size int64) LogOptionFunc {
	return func(c *LogConfig) {
		if size > 0 {
			c.rotationSize = size
		}
	}
}

// 设置文件最大保留时间、默认值7天
func SetLogMaxAge(maxAge time.Duration) LogOptionFunc {
	return func(c *LogConfig) {
		if maxAge != 0 {
			c.maxAge = maxAge
		}
	}
}

// 设置文件分割时间、默认值24*time.Hour(按天分割)
func SetRotationTime(rotationTime time.Duration) LogOptionFunc {
	return func(c *LogConfig) {
		if rotationTime != 0 {
			c.rotationTime = rotationTime
		}
	}
}

// 设置保留的最大文件数量、没有默认值(表示不限制)
func SetRotationCount(n uint) LogOptionFunc {
	return func(c *LogConfig) {
		if n != 0 {
			c.rotationCount = n
		}
	}
}
