package laya

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/layasugar/glogs"
	"github.com/layasugar/laya/gconf"
	"github.com/layasugar/laya/genv"
	"github.com/layasugar/laya/gpprof"
	"log"
)

type App struct {
	WebServer *gin.Engine
	cfg       *appCfg
}

type appCfg struct {
	cfgPath   string // 配置文件路径
	webServer bool   // 是否开启webserver
}

// DefaultApp 提供基础的日志服务, ginLog, webServer, 默认的配置路径"./conf/app.json"
func DefaultApp() *App {
	app := &App{
		cfg: &appCfg{
			webServer: true,
		},
	}

	app.initWithConfig()
	return app
}

func NewApp(options ...AppOption) *App {
	app := &App{cfg: &appCfg{}}
	for _, option := range options {
		option(app)
	}

	app.initWithConfig()
	return app
}

type AppOption func(app *App)

// 初始化app
func (app *App) initWithConfig() *App {
	// 初始化配置
	err := gconf.InitConfig(app.cfg.cfgPath)
	if err != nil {
		panic(err)
	}

	// 注册env
	app.registerEnv()

	// 开启日志系统
	glogs.InitLog(
		glogs.SetLogAppName(genv.AppName()),
		glogs.SetLogAppMode(genv.AppMode()),
		glogs.SetLogType(genv.LogType()),
		glogs.SetLogPath(genv.LogPath()),
	)

	// 是否需要初始化http服务
	if app.cfg.webServer {
		// 是否需要重定向gin日志输出
		if genv.RunMode() == gin.ReleaseMode {
			ginLogFile := genv.LogPath() + "/" + genv.AppName() + "/gin/%Y-%m-%d.log"
			gin.DefaultWriter = glogs.GetWriter(ginLogFile, glogs.DefaultConfig)
		}

		// 初始化http服务
		gin.SetMode(genv.RunMode())
		app.WebServer = gin.Default()

		// 开启必要中间件, requestID设置, 入参日志
		app.WebServer.Use(SetHeader, LogInParams)
	}

	// 是否开启pprof
	if genv.Pprof() {
		gpprof.StartPprof()
	}

	// 注册pprof监听函数和params监听函数和重载env函数
	gconf.RegisterConfigCharge(func() {
		if gconf.C.IsSet("app.pprof") {
			var newPprof = gconf.C.GetBool("app.pprof")
			var oldPprof = genv.Pprof()
			if oldPprof != newPprof {
				if newPprof {
					gpprof.StartPprof()
				} else {
					gpprof.StopPprof()
				}
			}
		} else {
			gpprof.StopPprof()
		}
	}, func() {
		app.registerEnv()
	})

	return app
}

// RunServer 运行服务
func (app *App) RunServer() {
	// 启动配置回调
	gconf.OnConfigCharge()

	// 启动web服务
	if app.cfg.webServer {
		log.Printf("%s %s %s starting at %q\n", genv.AppName(), genv.RunMode(), genv.AppUrl(), genv.HttpListen())
		err := app.WebServer.Run(genv.HttpListen())
		if err != nil {
			fmt.Printf("Can't RunWebServer: %s\n", err.Error())
		}
	}
}

func (app *App) Use(fc ...func()) {
	for _, f := range fc {
		f()
	}
}

// RegisterRouter 注册路由
func (app *App) RegisterRouter(rr func(*gin.Engine)) {
	rr(app.WebServer)
}

// set env
func (app *App) registerEnv() {
	genv.SetAppName(gconf.C.GetString("app.name"))
	genv.SetAppMode(gconf.C.GetString("app.mode"))
	genv.SetRunMode(gconf.C.GetString("app.run_mode"))
	genv.SetHttpListen(gconf.C.GetString("app.http_listen"))
	genv.SetAppUrl(gconf.C.GetString("app.url"))
	genv.SetPprof(gconf.C.GetBool("app.pprof"))
	genv.SetLogPath(gconf.C.GetString("app.logger"))
	genv.SetAppVersion(gconf.C.GetString("app.version"))

	if gconf.C.IsSet("app.params") {
		genv.SetParamLog(gconf.C.GetBool("app.params"))
	} else {
		genv.SetParamLog(true)
	}

	if genv.RunMode() == "release" {
		genv.SetLogType("file")
	}
}

// SetWebServer set web server
func SetWebServer() AppOption {
	return func(app *App) {
		app.cfg.webServer = true
	}
}

// SetConfigFile 设置配置文件
func SetConfigFile(filePath string) AppOption {
	return func(app *App) {
		app.cfg.cfgPath = filePath
	}
}
