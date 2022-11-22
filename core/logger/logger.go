// Package logger
// logger: this is extend package, use https://github.com/sirupsen/logrus
package logger

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/layasugar/laya/core/constants"
	"github.com/layasugar/laya/core/rotatelog"
	"github.com/sirupsen/logrus"
)

var sugar *logrus.Logger

var defaultConfig = &Config{
	appName:       "normal",
	appMode:       "debug",
	LogType:       "file",
	LogPath:       "/home/logs/app",
	childPath:     "%Y-%m-%d.log",
	RotationSize:  128 * 1024 * 1024,
	RotationCount: 30,
	RotationTime:  24 * time.Hour,
	MaxAge:        7 * 24 * time.Hour,
}

type Config struct {
	appName       string        // 应用名
	appMode       string        // 应用环境
	childPath     string        // 日志子路径+文件名
	LogType       string        // 日志类型
	LogPath       string        // 日志主路径
	LogLevel      string        // 日志等级
	RotationSize  int64         // 单个文件大小
	RotationCount uint          // 可以保留的文件个数
	RotationTime  time.Duration // 日志分割的时间
	MaxAge        time.Duration // 日志最大保留的天数
}

func GetSugar() *logrus.Logger {
	if sugar == nil {
		sugar = InitSugar(defaultConfig)
	}
	return sugar
}

func InitSugar(lc *Config) *logrus.Logger {
	level, err := logrus.ParseLevel(lc.LogLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
	logPath := fmt.Sprintf("%s/%s/%s", lc.LogPath, lc.appName, lc.childPath)
	if lc.LogType == "file" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(GetWriter(logPath, lc))
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}
	log.Printf("[app] logger success")
	return logrus.WithField("app_name", lc.appName).WithField("app_mode", lc.appMode).Logger
}

func Debug(logId, template string, args ...interface{}) {
	entry := GetSugar().WithField(constants.X_REQUESTID, logId)
	msg, entry := dealWithArgs(entry, template, args...)
	entry.Debug(msg)
}

func Info(logId, template string, args ...interface{}) {
	entry := GetSugar().WithField(constants.X_REQUESTID, logId)
	msg, entry := dealWithArgs(entry, template, args...)
	entry.Info(msg)
}

func Warn(logId, template string, args ...interface{}) {
	entry := GetSugar().WithField(constants.X_REQUESTID, logId)
	msg, entry := dealWithArgs(entry, template, args...)
	entry.Warn(msg)
}

func Error(logId, template string, args ...interface{}) {
	entry := GetSugar().WithField(constants.X_REQUESTID, logId)
	msg, entry := dealWithArgs(entry, template, args...)
	entry.Error(msg)
}

func dealWithArgs(entry *logrus.Entry, tmp string, args ...interface{}) (msg string, l *logrus.Entry) {
	l = entry
	if len(args) > 0 {
		var tmpArgs []interface{}
		for _, item := range args {
			if nil == item {
				continue
			}
			if fields, ok := item.(logrus.Fields); ok {
				l = l.WithFields(fields)
			} else {
				tmpArgs = append(tmpArgs, item)
			}
		}
		if len(tmpArgs) > 0 {
			msg = fmt.Sprintf(tmp, tmpArgs...)
		}
	}
	msg = tmp
	return
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
	var options []rotatelog.Option
	options = append(options,
		rotatelog.WithRotationSize(lc.RotationSize),
		rotatelog.WithRotationCount(lc.RotationCount),
		rotatelog.WithRotationTime(lc.RotationTime),
		rotatelog.WithMaxAge(lc.MaxAge))

	hook, err := rotatelog.New(
		filename,
		options...,
	)

	if err != nil {
		panic(err)
	}
	return hook
}
