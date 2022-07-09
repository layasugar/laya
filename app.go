package laya

import (
	"flag"
	"fmt"
	"github.com/layasugar/laya/gcal"
	"github.com/layasugar/laya/gcnf"
	"github.com/layasugar/laya/gcnf/env"
	"github.com/layasugar/laya/store/dbx"
	"github.com/layasugar/laya/store/edbx"
	"github.com/layasugar/laya/store/mdbx"
	"github.com/layasugar/laya/store/rdbx"
	"log"
)

type (
	WebContext     = httpx.WebContext
	WebServer      = httpx.WebServer
	WebHandlerFunc = httpx.WebHandlerFunc

	GrpcContext = grpcx.GrpcContext
	GrpcServer  = grpcx.GrpcServer

	App struct {
		// webServer 目前web引擎使用gin
		webServer *httpx.WebServer

		// grpcServer
		grpcServer *grpcx.GrpcServer

		// scene 是web还是grpc
		scene int
	}

	AppConfig struct {
		// HttpListen Web web 服务监听的地址
		HTTPListen string
		// PbPRCListen PbRPC服务监听的地址
		PbRPCListen string
	}
)

const (
	webApp = iota
	grpcApp
	defaultApp
)

const (
	mysqlConfKey    = "mysql"
	redisConfKey    = "redis"
	mongoConfKey    = "mongo"
	esConfKey       = "es"
	servicesConfKey = "services"
)

// DefaultApp 默认应用不带有web或者grpc, 可作为服务使用
func DefaultApp() *App {
	app := new(App)

	app.initWithConfig(-1)
	return app
}

// WebApp web app
func WebApp() *App {
	app := new(App)

	app.initWithConfig(webApp)
	return app
}

// GrpcApp grpc app
func GrpcApp() *App {
	app := new(App)

	app.initWithConfig(grpcApp)
	return app
}

// 初始化app
func (app *App) initWithConfig(scene int) *App {
	app.scene = scene

	// 接收命令行参数
	var f string
	flag.StringVar(&f, "config", "", "set a config file")
	flag.Parse()

	// 初始化配置
	err := gcnf.InitConfig(f)
	if err != nil {
		panic(err)
	}

	// 注册env
	app.register()

	switch scene {
	case webApp:
		if env.HttpListen() == "" {
			panic("app.http_listen is null")
		}
		app.webServer = httpx.NewWebServer(env.RunMode())
		if len(httpx.DefaultWebServerMiddlewares) > 0 {
			app.webServer.Use(httpx.DefaultWebServerMiddlewares...)
		}
	case grpcApp:
		if env.GrpcListen() == "" {
			panic("app.http_listen is null")
		}
		app.grpcServer = grpcx.NewGrpcServer()
	}

	return app
}

// RunServer 运行Web服务
func (app *App) RunServer() {
	switch app.scene {
	case webApp:
		// 启动web服务
		log.Printf("[app] Listening and serving %s on %s\n", "HTTP", env.HttpListen())
		err := app.webServer.Run(env.HttpListen())
		if err != nil {
			fmt.Printf("Can't RunWebServer: %s\n", err.Error())
		}
	case grpcApp:
		// 启动grpc服务
		log.Printf("[app] Listening and serving %s on %s\n", "GRPC", env.GrpcListen())
		err := app.grpcServer.Run(env.GrpcListen())
		if err != nil {
			log.Fatalf("Can't RunGrpcServer, GrpcListen: %s, err: %s", env.GrpcListen(), err.Error())
		}
	case defaultApp:
	}
}

// Use 提供一个加载函数
func (app *App) Use(fc ...func()) {
	for _, f := range fc {
		f()
	}
}

// register cal db services
func (app *App) register() {
	// 初始化调用gcal
	var services = gcnf.GetConfigMap(servicesConfKey)
	if len(services) > 0 {
		err := gcal.LoadService(services)
		if err != nil {
			log.Printf("[app] init load services error: %s", err.Error())
		}
	}

	// 初始化数据库连接和redis连接
	var dbs = gcnf.GetConfigMap(mysqlConfKey)
	var rdbs = gcnf.GetConfigMap(redisConfKey)
	var mdbs = gcnf.GetConfigMap(mongoConfKey)
	var edbs = gcnf.GetConfigMap(esConfKey)

	// 解析dbs
	dbx.InitConn(dbs)
	// 解析rdbs
	rdbx.InitConn(rdbs)
	// 解析mongo
	mdbx.InitConn(mdbs)
	// 解析es
	edbx.InitConn(edbs)
}

// SetNoLogParams 设置不需要打印的路由
func (app *App) SetNoLogParams(path ...string) {
	for _, v := range path {
		httpx.NoLogParamsRules.NoLogParams[v] = v
	}
}

// SetNoLogParamsPrefix 设置不需要打印入参和出参的路由前缀
func (app *App) SetNoLogParamsPrefix(path ...string) {
	for _, v := range path {
		httpx.NoLogParamsRules.NoLogParamsPrefix = append(httpx.NoLogParamsRules.NoLogParamsPrefix, v)
	}
}

// SetNoLogParamsSuffix 设置不需要打印的入参和出参的路由后缀
func (app *App) SetNoLogParamsSuffix(path ...string) {
	for _, v := range path {
		httpx.NoLogParamsRules.NoLogParamsSuffix = append(httpx.NoLogParamsRules.NoLogParamsSuffix, v)
	}
}

// WebServer 获取WebServer的指针
func (app *App) WebServer() *httpx.WebServer {
	return app.webServer
}

// GrpcServer 获取PbRPCServer的指针
func (app *App) GrpcServer() *grpcx.GrpcServer {
	return app.grpcServer
}

// NewContext 基础服务提供一个NewContext
func (app *App) NewContext(spanName string) *Context {
	return NewDefaultContext(spanName)
}
