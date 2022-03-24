package env

import (
	"github.com/layasugar/laya/gcf"
	"time"
)

const (
	defaultLogPath          = "/home/logs/app"
	defaultLogType          = "console"
	defaultLogMaxAge        = 7 * 24 * time.Hour
	defaultLogMaxCount uint = 30
)

// LogPath 返回日志基本路径
func LogPath() string {
	if gcf.IsSet("app.logger.path") {
		return gcf.GetString("app.logger.path")
	}
	return defaultLogPath
}

// LogType 返回日志类型
func LogType() string {
	if gcf.IsSet("app.logger.type") {
		return gcf.GetString("app.logger.type")
	}
	return defaultLogType
}

// LogMaxAge 返回日志默认保留7天
func LogMaxAge() time.Duration {
	if gcf.IsSet("app.logger.max_age") {
		return time.Duration(gcf.GetInt("app.logger.max_age")) * 24 * time.Hour
	}
	return defaultLogMaxAge
}

// LogMaxCount 返回日志默认限制为30个
func LogMaxCount() uint {
	if gcf.IsSet("app.logger.max_count") {
		return uint(gcf.GetInt("app.logger.max_count"))
	}
	return defaultLogMaxCount
}
