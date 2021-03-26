package genv

import (
	"github.com/layatips/laya/gutils"
	"os"
	"path/filepath"
)

var (
	envAppName     = ""
	envRunMode     = ""
	envConfDirName = ""
)

const (
	_DefaultAppName     = "unknown"
	_DefaultRunMode     = "debug"
	_defaultConfDirName = "gconf"
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

// 设置配置文件根目录名
func SetConfDirName(confDirName string) {
	envConfDirName = confDirName
}

// 返回配置文件根目录绝对地址
func ConfRootPath() string {
	if envConfDirName == "" {
		SetConfDirName(_defaultConfDirName)
	}
	return filepath.Join(RootPath(), envConfDirName)
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

// 返回data 目录的绝对地址
func DataRootPath() string {
	return filepath.Join(RootPath(), "data")
}

// 返回log根目录的绝对地址
func LogRootPath() string {
	return filepath.Join(RootPath(), "log")
}

// 自动寻找rootPath
func detectRootPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic("DefaultApp can't get current directory: " + err.Error())
	}
	var dir string

	binDir := filepath.Dir(os.Args[0])
	if !filepath.IsAbs(binDir) {
		binDir = filepath.Join(pwd, binDir)
	}
	// 如果有和可执行文件平级的conf目录，则当前目录就是根目录
	// 这通常是直接在代码目录里go build然后直接执行生成的结果
	dir = filepath.Join(binDir, "gconf")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		dir = binDir
		return dir
	}
	// 如果有和可执行文件上级平级的conf目录，则上层目录就是根目录
	// 这一般是按标准进行部署
	dir = filepath.Join(filepath.Dir(binDir), "gconf")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		dir = filepath.Dir(binDir)
		return dir
	}
	// 如果都没有，但可执行文件的父目录名称为bin，则bin的上一层就是根目录
	// 这种情况适用于配置目录名为：etc, gconf, configs等情况
	if filepath.Base(binDir) == "bin" {
		dir = filepath.Dir(binDir)
		return dir
	}
	// 如果都不符合，当前路径就是根目录，这一般是使用go run main.go的情况
	dir = pwd
	return dir
}

var ip string

// 本机IP
func LocalIP() string {
	if ip != "" {
		return ip
	}
	ip, _ = gutils.LocalIP()
	if ip == "" {
		ip = "unknown"
	}
	return ip
}
