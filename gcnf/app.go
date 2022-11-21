package gcnf

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
	if gcnf.IsSet(_appName) {
		return gcnf.GetString(_appName)
	}
	return defaultAppName
}

// AppMode 返回当前的环境
func AppMode() string {
	if gcnf.IsSet(_appMode) {
		return gcnf.GetString(_appMode)
	}
	return defaultAppMode
}

// RunMode 返回当前的运行模式
func RunMode() string {
	if gcnf.IsSet(_appRunMode) {
		return gcnf.GetString(_appRunMode)
	}
	return defaultRunMode
}

// AppVersion 返回app的版本号
func AppVersion() string {
	if gcnf.IsSet(_appVersion) {
		return gcnf.GetString(_appVersion)
	}
	return defaultAppVersion
}

// HttpListen 获取http监听地址
func HttpListen() string {
	if gcnf.IsSet(_appHttpListen) {
		return gcnf.GetString(_appHttpListen)
	}
	return defaultHttpListen
}

// GrpcListen 返回rpc监听地址
func GrpcListen() string {
	if gcnf.IsSet(_appGrpcListen) {
		return gcnf.GetString(_appGrpcListen)
	}
	return defaultGrpcListen
}

// AppUrl 返回当前app_url
func AppUrl() string {
	if gcnf.IsSet(_appUrl) {
		return gcnf.GetString(_appUrl)
	}
	return defaultAppUrl
}
