// Package log is a global internal logger
// logger: this is extend package, use https://github.com/uber-go/zap
package log

import (
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

var (
	ZapLog *zap.SugaredLogger // 简易版日志文件
	//Logger *zap.Logger // 这个日志强大一些, 目前还用不到
	logLevel = zap.NewAtomicLevel()
)

var path = "./config/db/db.toml"

type Config struct {
	Driver   string `toml:"driver"`
	Path     string `toml:"path"`
	LogLevel string `toml:"log_level"`
}

// InitLog 初始化日志文件
func Init() {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Printf("[store_db] parse db config %s failed,err= %s\n", path, err)
		return
	}

	loglevel := zapcore.InfoLevel
	switch config.LogLevel {
	case "INFO":
		loglevel = zapcore.InfoLevel
	case "ERROR":
		loglevel = zapcore.ErrorLevel
	}
	setLevel(loglevel)

	var core zapcore.Core
	// 打印至文件中
	if config.Driver == "file" {
		configs := zap.NewProductionEncoderConfig()
		configs.EncodeTime = zapcore.ISO8601TimeEncoder
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Path,
			MaxSize:    128, // MB
			LocalTime:  true,
			Compress:   true,
			MaxBackups: 8, // 最多保留 n 个备份
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
	ZapLog = logger.Sugar()
}

func setLevel(level zapcore.Level) {
	logLevel.SetLevel(level)
}

func Info(args ...interface{}) {
	ZapLog.Info(args...)
}

func InfoF(template string, args ...interface{}) {
	ZapLog.Infof(template, args...)
}

func Warn(args ...interface{}) {
	ZapLog.Warn(args...)
}

func WarnF(template string, args ...interface{}) {
	ZapLog.Warnf(template, args...)
}

func Error(args ...interface{}) {
	ZapLog.Error(args...)
}

func ErrorF(template string, args ...interface{}) {
	ZapLog.Errorf(template, args...)
}
