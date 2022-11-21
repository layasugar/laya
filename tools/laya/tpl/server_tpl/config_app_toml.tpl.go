package server_tpl

const ConfigAppTomlTpl = `## 应用名称, 运行模式, http端口, 应用外网地址, 是否开启入参出参打印, 应用版本号
[app]
url = ""
mod = "dev"
name = "{{.projectName}}"
run_mode = "debug"
params = true
version = "1.0.0"

## 应用日志的配置(路径和保留天数和打印在控制台还是文件中)(默认开启7天, 最大日志保留30个, 单个日志大小是128M, 7天内日志最大到3.75G)
[app.logger]
type = "file"
path = "/home/logs/app"
max_age = 7
max_count = 30

## 应用链路追踪上报(支持zipkin和jaeger), 采样率是0-1, 0是关闭链路追踪
## zipkin_addr参考设置http://127.0.0.1:9411/api/v2/spans
## jaeger_addr参考设置127.0.0.1:6831
[app.trace]
type = "jaeger"
addr = "127.0.0.1:6831"
mod = 1

## service name是唯一值, 对应服务名, 默认连接超时时间是(1500ms), 必须带上协议头(http,https,grpc)
[[services]]
name = "http_test"
addr = "http://127.0.0.1:10081"
[[services]]
name = "grpc_test"
addr = "grpc://127.0.0.1:10082"

## mysql配置, redis配置, mongo配置 自行合理配置, models/dao/base.go, 初始化配置
[mysql]
dsn = "root:123456@tcp(127.0.0.1:3306)/laya_template?charset=utf8&parseTime=True&loc=Local"
[redis]
addr = "127.0.0.1:6379"
db = 0
pwd = "123456"
[mongo]

## extra其他配置
[extra]
`
