package env

import "github.com/layasugar/laya/gcf"

const (
	defaultAppName    = "default-app"
	defaultAppMode    = "dev"
	defaultRunMode    = "debug"
	defaultAppVersion = "1.0.0"
	defaultHttpListen = "0.0.0.0:80"
	defaultGrpcListen = "0.0.0.0:10082"
	defaultAppUrl     = "http://127.0.0.1:80"

	_appName       = "app.name"
	_appMode       = "app.mode"
	_appRunMode    = "app.run_mode"
	_appVersion    = "app.version"
	_appHttpListen = "app.http_listen"
	_appGrpcListen = "app.grpc_listen"
	_appUrl        = "app.url"
)

// AppName 返回当前app名称
func AppName() string {
	if gcf.IsSet(_appName) {
		return gcf.GetString(_appName)
	}
	return defaultAppName
}

// AppMode 返回当前的环境
func AppMode() string {
	if gcf.IsSet(_appMode) {
		return gcf.GetString(_appMode)
	}
	return defaultAppMode
}

// RunMode 返回当前的运行模式
func RunMode() string {
	if gcf.IsSet(_appRunMode) {
		return gcf.GetString(_appRunMode)
	}
	return defaultRunMode
}

// AppVersion 返回app的版本号
func AppVersion() string {
	if gcf.IsSet(_appVersion) {
		return gcf.GetString(_appVersion)
	}
	return defaultAppVersion
}

// HttpListen 获取http监听地址
func HttpListen() string {
	if gcf.IsSet(_appHttpListen) {
		return gcf.GetString(_appHttpListen)
	}
	return defaultHttpListen
}

// GrpcListen 返回rpc监听地址
func GrpcListen() string {
	if gcf.IsSet(_appGrpcListen) {
		return gcf.GetString(_appGrpcListen)
	}
	return defaultGrpcListen
}

// AppUrl 返回当前app_url
func AppUrl() string {
	if gcf.IsSet(_appUrl) {
		return gcf.GetString(_appUrl)
	}
	return defaultAppUrl
}
