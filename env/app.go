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
)

// AppName 返回当前app名称
func AppName() string {
	if gcf.IsSet("app.name") {
		return gcf.GetString("app.name")
	}
	return defaultAppName
}

// AppMode 返回当前的环境
func AppMode() string {
	if gcf.IsSet("app.mode") {
		return gcf.GetString("app.mode")
	}
	return defaultAppMode
}

// RunMode 返回当前的运行模式
func RunMode() string {
	if gcf.IsSet("app.run_mode") {
		return gcf.GetString("app.run_mode")
	}
	return defaultRunMode
}

// AppVersion 返回app的版本号
func AppVersion() string {
	if gcf.IsSet("app.version") {
		return gcf.GetString("app.version")
	}
	return defaultAppVersion
}

// HttpListen 获取http监听地址
func HttpListen() string {
	if gcf.IsSet("app.http_listen") {
		return gcf.GetString("app.http_listen")
	}
	return defaultHttpListen
}

// GrpcListen 返回rpc监听地址
func GrpcListen() string {
	if gcf.IsSet("app.grpc_listen") {
		return gcf.GetString("app.grpc_listen")
	}
	return defaultGrpcListen
}

// AppUrl 返回当前app_url
func AppUrl() string {
	if gcf.IsSet("app.url") {
		return gcf.GetString("app.url")
	}
	return defaultAppUrl
}