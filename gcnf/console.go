package gcnf

// ApiLog 返回api是否打印入参和出参
func ApiLog() bool {
	if IsSet(_appConsoleApiLog) {
		return GetBool(_appConsoleApiLog)
	}
	return true
}

// ApiTrace 返回api是否提交链路追踪
func ApiTrace() bool {
	if gcf.IsSet(_appConsoleApiTrace) {
		return gcf.GetBool(_appConsoleApiTrace)
	}
	return true
}

// SdkLog 返回是否打印内部服务调用日志
func SdkLog() bool {
	if gcf.IsSet(_appConsoleSdkLog) {
		return gcf.GetBool(_appConsoleSdkLog)
	}
	return true
}

// MysqlLog 返回是否打印mysql查询日志
func MysqlLog() bool {
	if gcf.IsSet(_appConsoleMysqlLog) {
		return gcf.GetBool(_appConsoleMysqlLog)
	}
	return true
}

// MysqlTrace 返回是否提交mysql链路追踪
func MysqlTrace() bool {
	if gcf.IsSet(_appConsoleMysqlTrace) {
		return gcf.GetBool(_appConsoleMysqlTrace)
	}
	return true
}

// RedisTrace 返回是否提交redis链路追踪
func RedisTrace() bool {
	if gcf.IsSet(_appConsoleRedisTrace) {
		return gcf.GetBool(_appConsoleRedisTrace)
	}
	return true
}

// MongoTrace 返回是否提交mongo链路追踪
func MongoTrace() bool {
	if gcf.IsSet(_appConsoleMongoTrace) {
		return gcf.GetBool(_appConsoleMongoTrace)
	}
	return true
}

// EsTrace 返回是否提交es链路追踪
func EsTrace() bool {
	if gcf.IsSet("app.console.es_trace") {
		return gcf.GetBool("app.console.es_trace")
	}
	return true
}
