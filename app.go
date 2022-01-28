// surprise

package laya

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/layasugar/laya/gcal"
	"github.com/layasugar/laya/gconf"
	"github.com/layasugar/laya/genv"
	"log"
)

type App struct {
	// webEngine 目前web引擎使用gin
	webServer *WebServer

	// PbRPCServer
	pbRpcServer *PbRPCServer
}

type AppConfig struct {
	// HttpListen  Web web 服务监听的地址
	HTTPListen string
	// PbPRCListen PbRPC服务监听的地址
	PbRPCListen string
}

// DefaultApp 提供基础的日志服务, ginLog, webServer
func DefaultApp() *App {
	app := new(App)

	app.initWithConfig()
	return app
}

func NewApp() *App {
	app := &App{}
	app.initWithConfig()
	return app
}

type AppOption func(app *App)

// 初始化app
func (app *App) initWithConfig() *App {
	// 初始化配置
	err := gconf.InitConfig()
	if err != nil {
		panic(err)
	}

	// 注册env
	app.registerEnv()

	// 是否初始化web或者rpc
	if genv.HttpListen() != "" {
		app.webServer = NewWebServer(genv.RunMode())
		if len(DefaultWebServerMiddlewares) > 0 {
			app.webServer.Use(DefaultWebServerMiddlewares...)
		}
	}

	if genv.PbRpcListen() != "" {
		app.pbRpcServer = NewPbRPCServer()
	}

	// 注册pprof监听函数和params监听函数和重载env函数
	gconf.RegisterConfigCharge(func() {
		app.registerEnv()
	})

	// 启动配置回调
	gconf.OnConfigCharge()

	return app
}

// RunWebServer 运行Web服务
func (app *App) RunWebServer() {
	// 启动web服务
	log.Printf("[app] Listening and serving %s on %s\n", "HTTP", genv.HttpListen())
	err := app.webServer.Run(genv.HttpListen())
	if err != nil {
		fmt.Printf("Can't RunWebServer: %s\n", err.Error())
	}
}

// RunPbRPCServer 运行PbRPC服务
func (app *App) RunPbRPCServer() {
	err := app.pbRpcServer.Run(genv.PbRpcListen())
	if err != nil {
		log.Fatalf("Can't RunPbRPCServer, PbRPCListen=%s, err=%s", genv.PbRpcListen(), err.Error())
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
	genv.SetRunMode(gconf.V.GetString("app.run_mode"))
	log.Printf("[app] app.run_mode %s\n", genv.RunMode())
	genv.SetHttpListen(gconf.V.GetString("app.http_listen"))
	genv.SetPbRpcListen(gconf.V.GetString("app.pbrpc_liten"))

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

	// trace
	genv.SetTraceType(gconf.V.GetString("app.trace.type"))
	genv.SetTraceAddr(gconf.V.GetString("app.trace.addr"))
	genv.SetTraceMod(gconf.V.GetFloat64("app.trace.mod"))

	// 初始化调用gcal
	var services []map[string]interface{}
	s := gconf.V.Get("services")
	switch s.(type) {
	case []interface{}:
		si := s.([]interface{})
		for _, item := range si {
			if sim, ok := item.(map[string]interface{}); ok {
				services = append(services, sim)
			}
		}
	default:
		log.Printf("[app] init config error: services config")
	}
	if len(services) > 0 {
		err := gcal.LoadService(services)
		if err != nil {
			log.Printf("[app] init load services error: %s", err.Error())
		}
	}
}

// SetNoLogParams 设置不需要打印的路由
func (app *App) SetNoLogParams(path ...string) {
	for _, v := range path {
		noLogParamsRules.NoLogParams[v] = v
	}
}

// SetNoLogParamsPrefix 设置不需要打印入参和出参的路由前缀
func (app *App) SetNoLogParamsPrefix(path ...string) {
	for _, v := range path {
		noLogParamsRules.NoLogParamsPrefix = append(noLogParamsRules.NoLogParamsPrefix, v)
	}
}

// SetNoLogParamsSuffix 设置不需要打印的入参和出参的路由后缀
func (app *App) SetNoLogParamsSuffix(path ...string) {
	for _, v := range path {
		noLogParamsRules.NoLogParamsSuffix = append(noLogParamsRules.NoLogParamsSuffix, v)
	}
}

// WebServer 获取WebServer的指针
func (app *App) WebServer() *WebServer {
	return app.webServer
}

// PbRPCServer 获取PbRPCServer的指针
func (app *App) PbRPCServer() *PbRPCServer {
	return app.pbRpcServer
}

// DefaultWebServerMiddlewares 默认的Http Server中间件
// 其实应该保证TowerLogware 不panic，但是无法保证，多一个recovery来保证业务日志崩溃后依旧有访问日志
var DefaultWebServerMiddlewares = []WebHandlerFunc{
	LogParams,
	SetTrace,
	ginHandler2WebHandler(gin.Recovery()),
	recovery,
}

func recovery(ctx *WebContext) {
	defer func() {
		if err := recover(); err != nil {
			ctx.RspHTTPStatusCode = 500
			ctx.ErrMsg = fmt.Sprintf("%s", err)
			panic(err)
		}
	}()
	ctx.Next()
}
