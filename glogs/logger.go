// Package glogs is a global internal glogs
// glogs: this is extend package, use https://github.com/uber-go/zap
package glogs

import (
	"fmt"
	rl "github.com/layasugar/laya/glogs/logger"
	"github.com/layasugar/laya/genv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	DefaultChildPath    = "glogs/%Y-%m-%d.logger" // 默认子目录
	DefaultRotationSize = 128 * 1024 * 1024       // 默认大小为128M
	DefaultRotationTime = 24 * time.Hour          // 默认每天轮转一次
	DefaultNoBuffWrite  = false                   // 默认不开启无缓冲写入

	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var (
	Sugar *zap.Logger

	RequestIDName    = "x-b3-traceid"
	HeaderAppName    = "app-name"
	KeyPath          = "path"
	KeyTitle         = "title"
	KeyOriginAppName = "origin_app_name"
)

type Config struct {
	appName       string        // 应用名
	appMode       string        // 应用环境
	logType       string        // 日志类型
	logPath       string        // 日志主路径
	childPath     string        // 日志子路径+文件名
	RotationSize  int64         // 单个文件大小
	RotationCount uint          // 可以保留的文件个数
	NoBuffWrite   bool          // 设置无缓冲日志写入
	RotationTime  time.Duration // 日志分割的时间
	MaxAge        time.Duration // 日志最大保留的天数
}

func getSugar() *zap.Logger {
	if Sugar == nil {
		cfg := Config{
			appName:       genv.AppName(),
			appMode:       genv.RunMode(),
			logType:       genv.LogType(),
			logPath:       genv.LogPath(),
			childPath:     DefaultChildPath,
			RotationSize:  DefaultRotationSize,
			RotationCount: genv.LogMaxCount(),
			NoBuffWrite:   DefaultNoBuffWrite,
			RotationTime:  DefaultRotationTime,
			MaxAge:        genv.LogMaxAge(),
		}

		Sugar = InitSugar(&cfg)
	}

	return Sugar
}

func InitSugar(lc *Config) *zap.Logger {
	loglevel := zapcore.InfoLevel
	defaultLogLevel := zap.NewAtomicLevel()
	defaultLogLevel.SetLevel(loglevel)

	logPath := fmt.Sprintf("%s/%s/%s", lc.logPath, lc.appName, lc.childPath)

	var core zapcore.Core
	// 打印至文件中
	if lc.logType == "file" {
		configs := zap.NewProductionEncoderConfig()
		configs.FunctionKey = "func"
		configs.EncodeTime = timeEncoder

		w := zapcore.AddSync(GetWriter(logPath, lc))

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(configs),
			w,
			defaultLogLevel,
		)
		log.Printf("[glogs_sugar] logger success")
	} else {
		// 打印在控制台
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), defaultLogLevel)
		log.Printf("[glogs_sugar] logger success")
	}

	filed := zap.Fields(zap.String("app_name", lc.appName), zap.String("app_mode", lc.appMode))
	return zap.New(core, filed, zap.AddCaller(), zap.AddCallerSkip(3))
}

func Info(r *http.Request, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(r, LevelInfo, msg, LevelInfo, fields...)
}

func Warn(r *http.Request, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(r, LevelWarn, msg, LevelWarn, fields...)
}

func Error(r *http.Request, template string, args ...interface{}) {
	msg, fields := dealWithArgs(template, args...)
	writer(r, LevelError, msg, LevelError, fields...)
}

func dealWithArgs(tmp string, args ...interface{}) (msg string, f []zap.Field) {
	var tmpArgs []interface{}
	for _, item := range args {
		if zapField, ok := item.(zap.Field); ok {
			f = append(f, zapField)
		} else {
			tmpArgs = append(tmpArgs, item)
		}
	}
	msg = fmt.Sprintf(tmp, tmpArgs...)
	return
}

func writer(r *http.Request, level, msg string, title string, fields ...zap.Field) {
	if r == nil {
		fields = append(fields, zap.String(KeyTitle, title))
		switch level {
		case LevelInfo:
			getSugar().Info(msg, fields...)
		case LevelWarn:
			getSugar().Warn(msg, fields...)
		case LevelError:
			getSugar().Error(msg, fields...)
		}
		return
	}

	requestID := r.Header.Get(RequestIDName)
	originAppName := r.Header.Get(HeaderAppName)
	path := r.RequestURI
	fields = append(fields, zap.String(KeyPath, path),
		zap.String(RequestIDName, requestID),
		zap.String(KeyTitle, title),
		zap.String(KeyOriginAppName, originAppName))

	switch level {
	case LevelInfo:
		getSugar().Info(msg, fields...)
	case LevelWarn:
		getSugar().Warn(msg, fields...)
	case LevelError:
		getSugar().Error(msg, fields...)
	}
	return
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

// GetWriter 按天切割按大小切割
// filename 文件名
// RotationSize 每个文件的大小
// MaxAge 文件最大保留天数
// RotationCount 最大保留文件个数
// RotationTime 设置文件分割时间
// RotationCount 设置保留的最大文件数量
func GetWriter(filename string, lc *Config) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 stream-2021-5-20.logger
	// demo.log是指向最新日志的连接
	// 保存7天内的日志，每1小时(整点)分割一第二天志
	var options []rl.Option
	if lc.NoBuffWrite {
		options = append(options, rl.WithNoBuffer())
	}
	options = append(options,
		rl.WithRotationSize(lc.RotationSize),
		rl.WithRotationCount(lc.RotationCount),
		rl.WithRotationTime(lc.RotationTime),
		rl.WithMaxAge(lc.MaxAge),
		rl.ForceNewFile())

	hook, err := rl.New(
		filename,
		options...,
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func String(key string, value interface{}) zap.Field {
	v := fmt.Sprintf("%v", value)
	return zap.String(key, v)
}
