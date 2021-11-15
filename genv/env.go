package genv

var (
	envAppName    = ""
	envAppMode    = ""
	envRunMode    = ""
	envHttpListen = ""
	envAppVersion = ""
	envAppUrl     = ""
	envParamLog   = true
	envLogPath    = ""
	envLogType    = "console"
	envPprof      = false
	ConfigPath    = "./conf/app.json"
)

const (
	_DefaultAppName    = "default-app"
	_DefaultAppMode    = "dev"
	_DefaultRunMode    = "debug"
	_DefaultHttpListen = "0.0.0.0:10080"
	_DefaultAppVersion = "1.0.0"
	_DefaultAppUrl     = "127.0.0.1:10080"
	_DefaultLogPath    = "/home/logs/app"
	DefaultParamLog    = true
	DefaultPprof       = false
)

// SetAppName 设置app名称
func SetAppName(appName string) {
	envAppName = appName
}

// AppName 返回当前app名称
func AppName() string {
	if envAppName == "" {
		SetAppName(_DefaultAppName)
	}
	return envAppName
}

// SetAppMode 设置app运行环境
func SetAppMode(appMode string) {
	envAppMode = appMode
}

// AppMode 返回当前app运行环境
func AppMode() string {
	if envAppMode == "" {
		SetAppMode(_DefaultAppMode)
	}
	return envAppMode
}

// SetRunMode 设置运行模式
func SetRunMode(runMode string) {
	envRunMode = runMode
}

// RunMode 返回当前的运行模式
func RunMode() string {
	if envRunMode == "" {
		SetRunMode(_DefaultRunMode)
	}
	return envRunMode
}

// SetAppVersion 设置app的版本号
func SetAppVersion(appVersion string) {
	envAppVersion = appVersion
}

// AppVersion 返回app的版本号
func AppVersion() string {
	if envAppVersion == "" {
		SetAppVersion(_DefaultAppVersion)
	}
	return envAppVersion
}

// SetHttpListen 设置监听端口
func SetHttpListen(httpListen string) {
	envHttpListen = httpListen
}

// HttpListen 返回当前监听端口
func HttpListen() string {
	if envHttpListen == "" {
		SetHttpListen(_DefaultHttpListen)
	}
	return envHttpListen
}

// SetAppUrl 设置app_url
func SetAppUrl(appUrl string) {
	envAppUrl = appUrl
}

// AppUrl 返回当前app_url
func AppUrl() string {
	if envAppUrl == "" {
		SetAppUrl(_DefaultAppUrl)
	}
	return envAppUrl
}

// SetParamLog 设置是否打印入参和出参
func SetParamLog(ParamLog bool) {
	envParamLog = ParamLog
}

// ParamLog 返回是否打印入参和出参
func ParamLog() bool {
	return envParamLog
}

// SetLogPath 设置日志路径
func SetLogPath(path string) {
	envLogPath = path
}

// LogPath 返回日志基本路径
func LogPath() string {
	if envLogPath == "" {
		SetLogPath(_DefaultLogPath)
	}
	return envLogPath
}

// SetLogType 设置日志类型
func SetLogType(path string) {
	envLogType = path
}

// LogType 返回日志类型
func LogType() string {
	return envLogType
}

func SetPprof(pprof bool) {
	envPprof = pprof
}

func Pprof() bool {
	return envPprof
}
