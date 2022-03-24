package env

import "github.com/layasugar/laya/gcf"

// ApiLog 返回api是否打印入参和出参
func ApiLog() bool {
	if gcf.IsSet("app.console.api_log") {
		return gcf.GetBool("app.console.api_log")
	}
	return true
}

// ApiTrace 返回api是否提交链路追踪
func ApiTrace() bool {
	if gcf.IsSet("app.console.api_trace") {
		return gcf.GetBool("app.console.api_trace")
	}
	return true
}

// SdkLog 返回是否打印内部服务调用日志
func SdkLog() bool {
	if gcf.IsSet("app.console.sdk_log") {
		return gcf.GetBool("app.console.sdk_log")
	}
	return true
}

// MysqlLog 返回是否打印mysql查询日志
func MysqlLog() bool {
	if gcf.IsSet("app.console.mysql_log") {
		return gcf.GetBool("app.console.mysql_log")
	}
	return true
}

// MysqlTrace 返回是否提交mysql链路追踪
func MysqlTrace() bool {
	if gcf.IsSet("app.console.mysql_trace") {
		return gcf.GetBool("app.console.mysql_trace")
	}
	return true
}

// RedisTrace 返回是否提交redis链路追踪
func RedisTrace() bool {
	if gcf.IsSet("app.console.redis_trace") {
		return gcf.GetBool("app.console.redis_trace")
	}
	return true
}

// MongoTrace 返回是否提交mongo链路追踪
func MongoTrace() bool {
	if gcf.IsSet("app.console.mongo_trace") {
		return gcf.GetBool("app.console.mongo_trace")
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
