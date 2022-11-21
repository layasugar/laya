package constants

import "time"

const (
	CTimer     = 5 * time.Second
	envLocalIP = "127.0.0.1"
)

// 配置文件的key
const (
	KEY_MYSQL                = "mysql"
	KEY_REDIS                = "redis"
	KEY_MONGO                = "mongo"
	KEY_ES                   = "es"
	KEY_SERVICES             = "services"
	KEY_APPCONSOLEAPILOG     = "app.console.api_log"
	KEY_APPCONSOLEAPITRACE   = "app.console.api_trace"
	KEY_APPCONSOLESDKLOG     = "app.console.sdk_log"
	KEY_APPCONSOLEMYSQLLOG   = "app.console.mysql_log"
	KEY_APPCONSOLEMYSQLTRACE = "app.console.mysql_trace"
	KEY_APPCONSOLEREDISTRACE = "app.console.redis_trace"
	KEY_APPCONSOLEMONGOTRACE = "app.console.mongo_trace"
	KEY_APPCONSOLEESTRACE    = "app.console.es_trace"
	KEY_APPLOGGERPATH        = "app.logger.path"
	KEY_APPLOGGERTYPE        = "app.logger.type"
	KEY_APPLOGGERMAXAGE      = "app.logger.max_age"
	KEY_APPLOGGERMAXCOUNT    = "app.logger.max_count"
	KEY_APPTRACETYPE         = "app.trace.type"
	KEY_APPTRACEADDR         = "app.trace.addr"
	KEY_APPTRACEMOD          = "app.trace.mod"
)

// 默认参数
const (
	DEFAULT_CONFIGFILE          = "conf/app.toml"
	DEFAULT_LOGPATH             = "/home/logs/app"
	DEFAULT_LOGTYPE             = "console"
	DEFAULT_LOGMAXAGE           = 7 * 24 * time.Hour
	DEFAULT_LOGMAXCOUNT uint    = 30
	DEFAULT_TRACETYPE           = ""
	DEFAULT_TRACEADDR           = ""
	DEFAULT_TRACEMOD    float64 = 0
)

const (
	X_FORWARDEDFOR = "X-Forwarded-For" // 获取真实ip
	X_REALIP       = "X-Real-IP"       // 获取真实ip
	X_REQUESTID    = "x_request_id"    // 日志key
)
