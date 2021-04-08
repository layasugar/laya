package genv

var (
	envAppName    = ""
	envRunMode    = ""
	envHttpListen = ""
	envAppVersion = ""
	envAppUrl     = ""
	envParamLog   = true
)

const (
	_DefaultAppName    = "default-app"
	_DefaultHttpListen = "0.0.0.0:10080"
	_DefaultRunMode    = "debug"
	_DefaultAppVersion = "1.0.0"
	_DefaultAppUrl     = "127.0.0.1:10080"
)

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

// 设置运行模式
func SetRunMode(runMode string) {
	envRunMode = runMode
}

// 返回当前的运行模式
func RunMode() string {
	if envRunMode == "" {
		SetRunMode(_DefaultRunMode)
	}
	return envRunMode
}

// 设置app的版本号
func SetAppVersion(appVersion string) {
	envAppVersion = appVersion
}

// 返回app的版本号
func AppVersion() string {
	if envAppVersion == "" {
		SetAppVersion(_DefaultAppVersion)
	}
	return envAppVersion
}

// 设置监听端口
func SetHttpListen(httpListen string) {
	envHttpListen = httpListen
}

// 返回当前监听端口
func HttpListen() string {
	if envHttpListen == "" {
		SetHttpListen(_DefaultHttpListen)
	}
	return envHttpListen
}

// 设置app_url
func SetAppUrl(appUrl string) {
	envAppUrl = appUrl
}

// 返回当前app_url
func AppUrl() string {
	if envAppUrl == "" {
		SetAppUrl(_DefaultAppUrl)
	}
	return envAppUrl
}

// 设置是否打印入参和出参
func SetParamLog(ParamLog bool) {
	envParamLog = ParamLog
}

// 返回是否打印入参和出参
func ParamLog() bool {
	return envParamLog
}
