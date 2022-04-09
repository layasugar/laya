# laya

基本框架, 只支持(http, Grpc, 基础服务)

## 快速开始

[使用模板快速构建项目](https://github.com/layasugar/laya-template)

## 功能

- [x] 应用初始化, 包含http和grpc应用
- [x] 配置文件初始化, 配置文件热重载signal 等基础功能
- [x] 提供全局传递的WebContext, GrpcContext, Context
- [x] 提供完全兼容的ginRoute和GrpcRoute
- [x] 提供完善的日志功能(包含grpc和http的日志跟踪)
- [x] http中间件与grpc拦截器完成日志和链路追踪
- [x] 链路支持zipkin与jaeger(包含http与grpc和基础服务)
- [x] 链路追踪包含mysql
- [x] 链路追踪包含redis
- [x] 日志打印error时触发alarm
- [x] 基础app配置增加环境配置(有用环境做应用隔离的需求)
- [x] 精准控制mysql,redis,mongo,es是否接入链路和日志
- [x] gcal支持请求除内部服务外的其他http服务
- [ ] 使用cmux实现多路复用[cmux-github](https://github.com/soheilhy/cmux), 使http与grpc监听同一个端口

## 工具

- [x] 一键初始化目录结构到当前目录
- [x] 一键生成db.model


