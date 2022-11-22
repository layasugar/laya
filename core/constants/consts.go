package constants

import "time"

// 配置文件的key
const (
	KEY_APPNAME                  = "app.name"
	KEY_APPMODE                  = "app.mode"
	KEY_APPVERSION               = "app.version"
	KEY_APPLISTEN                = "app.listen"
	KEY_APPLOGGERPATH            = "app.logger.path"
	KEY_APPLOGGERLEVEL           = "app.logger.level"
	KEY_APPLOGGERTYPE            = "app.logger.type"
	KEY_APPLOGGERMAXAGE          = "app.logger.max_age"
	KEY_APPLOGGERMAXCOUNT        = "app.logger.max_count"
	KEY_APPLOGGERMAXSIZE         = "app.logger.max_size"
	KEY_APPLOGGERMAXTIME         = "app.logger.max_time"
	KEY_APPLOGPARAMSLOGURI       = "app.logger.params.log_uri"
	KEY_APPLOGPARAMSLOGPREFIXURI = "app.logger.params.log_prefix_uri"
	KEY_APPLOGPARAMSLOGSUFFIXURI = "app.logger.params.log_suffix_uri"
	KEY_APPTRACETYPE             = "app.trace.type"
	KEY_APPTRACEADDR             = "app.trace.addr"
	KEY_APPTRACEMOD              = "app.trace.mod"
	KEY_APPALARMTYPE             = "app.alarm.type"
	KEY_APPALARMKEY              = "app.alarm.key"
	KEY_APPALARMADDR             = "app.alarm.addr"
	KEY_MYSQL                    = "mysql"
	KEY_PG                       = "pgsql"
	KEY_MONGO                    = "mongo"
	KEY_ES                       = "es"
	KEY_REDIS                    = "redis"
	KEY_SERVICES                 = "services"
)

// 默认参数
const (
	DEFAULT_CONFIGFILE          = "conf/app.toml"
	DEFAULT_LOGPATH             = "/home/logs/app"
	DEFAULT_LOGLEVEL            = "info"
	DEFAULT_LOGTYPE             = "file"
	DEFAULT_LOGMAXAGE           = 7 * 24 * time.Hour
	DEFAULT_LOGMAXTIME          = 24 * time.Hour
	DEFAULT_LOGMAXCOUNT uint    = 30
	DEFAULT_LOGMAXSIZE  int64   = 134217728
	DEFAULT_TRACETYPE           = ""
	DEFAULT_TRACEADDR           = ""
	DEFAULT_TRACEMOD    float64 = 0
	DEFAULT_APPNAME             = "normal"
	DEFAULT_APPMODE             = "dev"
	DEFAULT_APPVERSION          = "1.0.0"
	DEFAULT_LISTEN              = "0.0.0.0:10080"
	DEFAULT_NULLSTRING          = ""
)

const (
	X_FORWARDEDFOR  = "X-Forwarded-For" // 获取真实ip
	X_REALIP        = "X-Real-IP"       // 获取真实ip
	X_REQUESTID     = "x_request_id"    // 日志key
	PROTOCOLHTTP    = "HTTP"
	PROTOCOLGRPC    = "GRPC"
	TRACETYPEJAEGER = "jaeger"
	TRACETYPEZIPKIN = "zipkin"
)

const (
	SERVERGIN SERVERTYPE = iota + 1
	SERVERGRPC
	SERVERNORMAL
)

const (
	RUNMODEDEV  RUNMODE = "dev"
	RUNMODETEST RUNMODE = "test"
	RUNMODEPRE  RUNMODE = "pre"
	RUNMODEPROD RUNMODE = "prod"
)
