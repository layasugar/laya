// Package log is a global internal logger
// logger: this is extend package, use https://github.com/uber-go/zap
package llog

import (
	"github.com/LaYa-op/laya/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func init() {

}

var (
	zapLog *zap.SugaredLogger // 简易版日志文件
	//Logger *zap.Logger // 这个日志强大一些, 目前还用不到
	logLevel = zap.NewAtomicLevel()
)

type Config struct {
	Driver string `toml:"driver"`
	Path   string `toml:"path"`
}

// InitLog 初始化日志文件
func InitLog() {

	err := config.ReadFile(name, &config)
	logConf := config.GetLogConf()

	loglevel := zapcore.InfoLevel
	switch logConf.LogLevel {
	case "INFO":
		loglevel = zapcore.InfoLevel
	case "ERROR":
		loglevel = zapcore.ErrorLevel
	}
	setLevel(loglevel)

	var core zapcore.Core
	// 打印至文件中
	if logConf.LogType == "file" {
		config := zap.NewProductionEncoderConfig()
		config.EncodeTime = zapcore.ISO8601TimeEncoder
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logConf.LogPath,
			MaxSize:    128, // MB
			LocalTime:  true,
			Compress:   true,
			MaxBackups: 8, // 最多保留 n 个备份
		})

		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			w,
			logLevel,
		)
	} else {
		// 打印在控制台
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), logLevel)
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zapLog = logger.Sugar()
}

func setLevel(level zapcore.Level) {
	logLevel.SetLevel(level)
}

func Info(args ...interface{}) {
	zapLog.Info(args...)
}

func InfoF(template string, args ...interface{}) {
	zapLog.Infof(template, args...)
}

func Warn(args ...interface{}) {
	zapLog.Warn(args...)
}

func WarnF(template string, args ...interface{}) {
	zapLog.Warnf(template, args...)
}

func Error(args ...interface{}) {
	zapLog.Error(args...)
}

func ErrorF(template string, args ...interface{}) {
	zapLog.Errorf(template, args...)
}
