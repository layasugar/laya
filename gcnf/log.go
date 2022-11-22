package gcnf

import (
	"time"
)

// LogPath 返回日志基本路径
func LogPath() string {
	if IsSet(_appLoggerPath) {
		return GetString(_appLoggerPath)
	}
	return defaultLogPath
}

// LogType 返回日志类型
func LogType() string {
	if IsSet(_appLoggerType) {
		return GetString(_appLoggerType)
	}
	return defaultLogType
}

// LogMaxAge 返回日志默认保留7天
func LogMaxAge() time.Duration {
	if IsSet(_appLoggerMaxAge) {
		return time.Duration(GetInt(_appLoggerMaxAge)) * 24 * time.Hour
	}
	return defaultLogMaxAge
}

// LogMaxCount 返回日志默认限制为30个
func LogMaxCount() uint {
	if IsSet(_appLoggerMaxCount) {
		return uint(GetInt(_appLoggerMaxCount))
	}
	return defaultLogMaxCount
}
