package genv

import (
	"os"
)

var (
	envAppName = ""
	envRunMode = ""
)

const (
	_DefaultAppName = "unknown"
	_DefaultRunMode = "debug"
)

// 设置运行模式
func SetRunMode(runMode string) {
	envRunMode = runMode
}

// 返回当前的运行模式
func RunMode() string {
	if envRunMode == "" {
		if os.Getenv("GIN_RUN_MODE") != "" {
			SetRunMode(os.Getenv("GIN_RUN_MODE"))
		} else {
			SetRunMode(_DefaultRunMode)
		}
	}
	return envRunMode
}

// 设置app名称
func SetAppName(appName string) {
	envAppName = appName
}

// 返回当前app名称
func AppName() string {
	if envAppName == "" {
		SetAppName(_DefaultAppName)
	}
	return envAppName
}
