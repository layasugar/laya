package gcnf

import (
	"time"
)

// LogPath 返回日志基本路径
func LogPath() string {
	if gcf.IsSet(_appLoggerPath) {
		return gcf.GetString(_appLoggerPath)
	}
	return defaultLogPath
}

// LogType 返回日志类型
func LogType() string {
	if gcf.IsSet(_appLoggerType) {
		return gcf.GetString(_appLoggerType)
	}
	return defaultLogType
}

// LogMaxAge 返回日志默认保留7天
func LogMaxAge() time.Duration {
	if gcf.IsSet(_appLoggerMaxAge) {
		return time.Duration(gcf.GetInt(_appLoggerMaxAge)) * 24 * time.Hour
	}
	return defaultLogMaxAge
}

// LogMaxCount 返回日志默认限制为30个
func LogMaxCount() uint {
	if gcf.IsSet(_appLoggerMaxCount) {
		return uint(gcf.GetInt(_appLoggerMaxCount))
	}
	return defaultLogMaxCount
}
