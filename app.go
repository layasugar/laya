// surprise

package laya

import (
	"flag"
	"fmt"
	"github.com/layasugar/laya/core/appx"
	"github.com/layasugar/laya/core/grpcx"
	"github.com/layasugar/laya/core/httpx"
	"github.com/layasugar/laya/gcal"
	"github.com/layasugar/laya/gconf"
	"github.com/layasugar/laya/genv"
	"github.com/layasugar/laya/gstore/dbx"
	"github.com/layasugar/laya/gstore/edbx"
	"github.com/layasugar/laya/gstore/mdbx"
	"github.com/layasugar/laya/gstore/rdbx"
	"log"
)

type (
	Context = appx.Context

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
	err := gconf.InitConfig(f)
	if err != nil {
		panic(err)
	}

	// 注册env
	app.registerEnv()

	// db init and rdb init
	app.initDbConn()

	switch scene {
	case webApp:
		if genv.HttpListen() == "" {
			panic("app.http_listen is null")
		}
		app.webServer = httpx.NewWebServer(genv.RunMode())
		if len(httpx.DefaultWebServerMiddlewares) > 0 {
			app.webServer.Use(httpx.DefaultWebServerMiddlewares...)
		}
	case grpcApp:
		if genv.GrpcListen() == "" {
			panic("app.http_listen is null")
		}
		app.grpcServer = grpcx.NewGrpcServer()
	}

	// 注册pprof监听函数和params监听函数和重载env函数
	gconf.RegisterConfigCharge(func() {
		app.registerEnv()
	})

	// 启动配置回调
	gconf.OnConfigCharge()

	return app
}

// RunServer 运行Web服务
func (app *App) RunServer() {
	switch app.scene {
	case webApp:
		// 启动web服务
		log.Printf("[app] Listening and serving %s on %s\n", "HTTP", genv.HttpListen())
		err := app.webServer.Run(genv.HttpListen())
		if err != nil {
			fmt.Printf("Can't RunWebServer: %s\n", err.Error())
		}
	case grpcApp:
		// 启动grpc服务
		log.Printf("[app] Listening and serving %s on %s\n", "GRPC", genv.GrpcListen())
		err := app.grpcServer.Run(genv.GrpcListen())
		if err != nil {
			log.Fatalf("Can't RunGrpcServer, GrpcListen: %s, err: %s", genv.GrpcListen(), err.Error())
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

// set env
func (app *App) registerEnv() {
	genv.SetAppUrl(gconf.V.GetString("app.url"))
	genv.SetAppName(gconf.V.GetString("app.name"))
	log.Printf("[app] app.name %s\n", genv.AppName())
	genv.SetAppMode(gconf.V.GetString("app.mode"))
	genv.SetRunMode(gconf.V.GetString("app.run_mode"))
	log.Printf("[app] app.run_mode %s\n", genv.RunMode())
	genv.SetHttpListen(gconf.V.GetString("app.http_listen"))
	genv.SetGrpcListen(gconf.V.GetString("app.grpc_listen"))

	if gconf.V.IsSet("app.params") {
		genv.SetParamLog(gconf.V.GetBool("app.params"))
	} else {
		genv.SetParamLog(true)
	}
	genv.SetAppVersion(gconf.V.GetString("app.gversion"))

	// 日志
	genv.SetLogPath(gconf.V.GetString("app.logger.path"))
	genv.SetLogType(gconf.V.GetString("app.logger.type"))
	genv.SetLogMaxAge(gconf.V.GetInt("app.logger.max_age"))
	genv.SetLogMaxCount(gconf.V.GetInt("app.logger.max_count"))

	// tracex
	genv.SetTraceType(gconf.V.GetString("app.trace.type"))
	genv.SetTraceAddr(gconf.V.GetString("app.trace.addr"))
	genv.SetTraceMod(gconf.V.GetFloat64("app.trace.mod"))

	// alarmx
	genv.SetAlarmType(gconf.V.GetString("app.alarm.type"))
	genv.SetAlarmKey(gconf.V.GetString("app.alarm.key"))
	genv.SetAlarmHost(gconf.V.GetString("app.alarm.addr"))

	// 初始化调用gcal
	var services = gconf.GetConfigMap(servicesConfKey)
	if len(services) > 0 {
		err := gcal.LoadService(services)
		if err != nil {
			log.Printf("[app] init load services error: %s", err.Error())
		}
	}
}

// 初始化数据库连接和redis连接
func (app *App) initDbConn() {
	var dbs = gconf.GetConfigMap(mysqlConfKey)
	var rdbs = gconf.GetConfigMap(redisConfKey)
	var mdbs = gconf.GetConfigMap(mongoConfKey)
	var edbs = gconf.GetConfigMap(esConfKey)

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
func (app *App) NewContext(logId string, spanName string) *appx.Context {
	return appx.NewDefaultContext(logId, spanName)
}
