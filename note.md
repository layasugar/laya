## 日志
## 功能代办
 - 链路支持zipkin与jaeger(包含http与grpc) 1
 - 链路追踪必须包含数据库(慢查询和error都必须进链路追踪)基于尾部连贯采样
 - 请求支持http, https, grpc
 - http => rpc
 - http => http 1
 - rpc => http
 - rpc => rpc 
 - 请求支持代理请求
 - 请求需要打印日志, 包含入参出参

## 工具代办
 - 一键初始化目录结构到当前目录
 - 一键生成db.model
 - mysql redis mongo es 封装一层(支持链路追踪)
 - 请求支持代理请求+链路跟踪
 