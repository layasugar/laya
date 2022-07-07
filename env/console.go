package env

const (
	_appConsoleApiLog     = "app.console.api_log"
	_appConsoleApiTrace   = "app.console.api_trace"
	_appConsoleSdkLog     = "app.console.sdk_log"
	_appConsoleMysqlLog   = "app.console.mysql_log"
	_appConsoleMysqlTrace = "app.console.mysql_trace"
	_appConsoleRedisTrace = "app.console.redis_trace"
	_appConsoleMongoTrace = "app.console.mongo_trace"
	_appConsoleEsTrace    = "app.console.es_trace"
)

// ApiLog 返回api是否打印入参和出参
func ApiLog() bool {
	if gcnf.IsSet(_appConsoleApiLog) {
		return gcnf.GetBool(_appConsoleApiLog)
	}
	return true
}

// ApiTrace 返回api是否提交链路追踪
func ApiTrace() bool {
	if gcnf.IsSet(_appConsoleApiTrace) {
		return gcnf.GetBool(_appConsoleApiTrace)
	}
	return true
}

// SdkLog 返回是否打印内部服务调用日志
func SdkLog() bool {
	if gcnf.IsSet(_appConsoleSdkLog) {
		return gcnf.GetBool(_appConsoleSdkLog)
	}
	return true
}

// MysqlLog 返回是否打印mysql查询日志
func MysqlLog() bool {
	if gcnf.IsSet(_appConsoleMysqlLog) {
		return gcnf.GetBool(_appConsoleMysqlLog)
	}
	return true
}

// MysqlTrace 返回是否提交mysql链路追踪
func MysqlTrace() bool {
	if gcnf.IsSet(_appConsoleMysqlTrace) {
		return gcnf.GetBool(_appConsoleMysqlTrace)
	}
	return true
}

// RedisTrace 返回是否提交redis链路追踪
func RedisTrace() bool {
	if gcnf.IsSet(_appConsoleRedisTrace) {
		return gcnf.GetBool(_appConsoleRedisTrace)
	}
	return true
}

// MongoTrace 返回是否提交mongo链路追踪
func MongoTrace() bool {
	if gcnf.IsSet(_appConsoleMongoTrace) {
		return gcnf.GetBool(_appConsoleMongoTrace)
	}
	return true
}

// EsTrace 返回是否提交es链路追踪
func EsTrace() bool {
	if gcnf.IsSet("app.console.es_trace") {
		return gcnf.GetBool("app.console.es_trace")
	}
	return true
}
