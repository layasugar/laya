package env

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/LaYa-op/laya/utils"
)

var (
	envAppName     = ""
	envRootPath    = ""
	envIDC         = ""
	envRunMode     = ""
	envConfDirName = ""
)

const (
	_DefaultAppName     = "unknown"
	_DefaultIDC         = "test"
	_DefaultRunMode     = "debug"
	_defaultConfDirName = "config"
)

// SetRootPath 设置应用的根目录
func SetRootPath(rootPath string) {
	if envRootPath != "" && rootPath != envRootPath {
		panic("app root path cannot set twice")
	}
	envRootPath = rootPath
	log.Printf("[env] SetRootPath=%s\n", rootPath)
}

// RootPath 返回应用的根目录
func RootPath() string {
	if envRootPath == "" {
		if os.Getenv("GXE_ROOT_PATH") != "" {
			SetRootPath(os.Getenv("GXE_ROOT_PATH"))
		} else {
			SetRootPath(detectRootPath())
		}
	}
	return envRootPath
}

// SetIDC 设置应用的idc
func SetIDC(idc string) {
	envIDC = idc
	log.Printf("[env] SetIDC= %s\n", idc)
}

// IDC 返回应用的机房
func IDC() string {
	if envIDC == "" {
		SetIDC(_DefaultIDC)
	}
	return envIDC
}

// SetRunMode 设置运行模式
func SetRunMode(runMode string) {
	envRunMode = runMode
	log.Printf("[env] SetRunMode= %s\n", runMode)
}

// RunMode 返回当前的运行模式
func RunMode() string {
	if envRunMode == "" {
		if os.Getenv("GXE_RUN_MODE") != "" {
			SetRunMode(os.Getenv("GXE_RUN_MODE"))
		} else {
			SetRunMode(_DefaultRunMode)
		}
	}
	return envRunMode
}

// SetConfDirName 设置配置文件根目录名
func SetConfDirName(confDirName string) {
	envConfDirName = confDirName
	log.Printf("[env] SetConfDirName= %s\n", confDirName)
}

// ConfRootPath 返回配置文件根目录绝对地址
func ConfRootPath() string {
	if envConfDirName == "" {
		SetConfDirName(_defaultConfDirName)
	}
	return filepath.Join(RootPath(), envConfDirName)
}

// SetAppName 设置app名称
func SetAppName(appName string) {
	envAppName = appName
	log.Printf("[env] SetAppName=%s\n", appName)
}

// AppName 返回当前app名称
func AppName() string {
	if envAppName == "" {
		SetAppName(_DefaultAppName)
	}
	return envAppName
}

// DataRootPath 返回data 目录的绝对地址
func DataRootPath() string {
	return filepath.Join(RootPath(), "data")
}

// LogRootPath 返回log根目录的绝对地址
func LogRootPath() string {
	return filepath.Join(RootPath(), "log")
}

// detectRootPath 自动寻找rootpath
func detectRootPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic("DefaultApp can't get current directory: " + err.Error())
	}
	var dir string
	defer func() {
		log.Printf("[env] auto detect rootPath= %s\n", dir)
	}()

	bindir := filepath.Dir(os.Args[0])
	if !filepath.IsAbs(bindir) {
		bindir = filepath.Join(pwd, bindir)
	}
	// 如果有和可执行文件平级的conf目录，则当前目录就是根目录
	// 这通常是直接在代码目录里go build然后直接执行生成的结果
	dir = filepath.Join(bindir, "config")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		dir = bindir
		return dir
	}
	// 如果有和可执行文件上级平级的conf目录，则上层目录就是根目录
	// 这一般是按标准进行部署
	dir = filepath.Join(filepath.Dir(bindir), "config")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		dir = filepath.Dir(bindir)
		return dir
	}
	// 如果都没有，但可执行文件的父目录名称为bin，则bin的上一层就是根目录
	// 这种情况适用于配置目录名为：etc, config, configs等情况
	if filepath.Base(bindir) == "bin" {
		dir = filepath.Dir(bindir)
		return dir
	}
	// 如果都不符合，当前路径就是根目录，这一般是使用go run main.go的情况
	dir = pwd
	return dir
}

var ip string

// LocalIP 本机IP
func LocalIP() string {
	if ip != "" {
		return ip
	}

	ip, _ = utils.LocalIP()
	if ip == "" {
		ip = "unknown"
	}
	return ip
}

var pid int

// PID 得到 PID
func PID() int {
	if pid != 0 {
		return pid
	}

	pid = os.Getpid()
	pidstring = strconv.Itoa(pid)
	return pid
}

var pidstring string

// PIDString 得到PID 字符串形式
func PIDString() string {
	if pidstring == "" {
		PID()
	}

	return pidstring
}
