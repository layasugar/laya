package laya

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/layasugar/glogs"
	"github.com/layasugar/laya/gconf"
	"github.com/layasugar/laya/genv"
	"github.com/layasugar/laya/gpprof"
	"log"
	"path/filepath"
	"sync"
)

type App struct {
	WebServer *gin.Engine
	watcher   *gconf.Watcher
	watchLock sync.Mutex
}

func DefaultApp() *App {
	return NewApp(SetLogger, SetGinLog, SetWebServer, SetPprof)
}

func NewApp(options ...AppOption) *App {
	app := new(App).initWithConfig()
	for _, option := range options {
		option(app)
	}
	return app
}

type AppOption func(app *App)

func (app *App) initWithConfig() *App {
	err := gconf.InitConfig(genv.ConfigPath)
	if err != nil {
		panic(err)
	}

	cf, err := gconf.GetBaseConf()
	if err != nil && !errors.Is(err, gconf.Nil) {
		panic(err.Error())
	}
	if errors.Is(err, gconf.Nil) {
		cf = &gconf.App{
			Name:       "default-app",
			HttpListen: "0.0.0.0:10080",
			RunMode:    "debug",
			Version:    "1.0.0",
			Url:        "127.0.0.1:10080",
			ParamLog:   true,
			LogPath:    "/home/logs/app",
		}
	}
	app.registerEnv(cf)
	return app
}

func (app *App) RunWebServer() {
	log.Printf("%s %s %s starting at %q\n", genv.AppName(), genv.RunMode(), genv.AppUrl(), genv.HttpListen())
	err := app.WebServer.Run(genv.HttpListen())
	if err != nil {
		fmt.Printf("Can't RunWebServer: %s\n", err.Error())
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

// RegisterWatcher
func (app *App) RegisterFileWatcher(path string, fh gconf.WatcherEventHandler) {
	app.watchLock.Lock()
	defer app.watchLock.Unlock()

	if app.watcher == nil {
		app.watcher = gconf.NewWatcher(filepath.Dir(path), 65535)
		//默认需要监听所有的event
		go app.watcher.RegisterFileWatcher(filepath.Base(path), fh)
	}
}

// set env
func (app *App) registerEnv(cf *gconf.App) {
	if cf.Name != "" {
		genv.SetAppName(cf.Name)
	}
	if cf.HttpListen != "" {
		genv.SetHttpListen(cf.HttpListen)
	}
	if cf.RunMode != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.Version != "" {
		genv.SetAppVersion(cf.Version)
	}
	if cf.Url != "" {
		genv.SetAppUrl(cf.Url)
	}
	if cf.LogPath != "" {
		genv.SetLogPath(cf.LogPath)
	}
	if cf.Mode != "" {
		genv.SetAppMode(cf.Mode)
	}
	if cf.Pprof {
		genv.SetPprof(cf.Pprof)
	}
	if genv.RunMode() == "release" {
		genv.SetLogType("file")
	}
	genv.SetParamLog(cf.ParamLog)
}

// set gin logger
func SetGinLog(app *App) {
	if genv.AppMode() == "release" {
		// 设置gin的请求日志
		ginLogFile := genv.LogPath() + "/" + genv.AppName() + "/gin-http" + "/%Y-%m-%d.log"
		gin.DefaultWriter = glogs.GetWriter(ginLogFile, glogs.DefaultConfig)
	}
}

// set gin logger noBuffer
func SetGinLogNoBuffer(app *App) {
	if genv.AppMode() == gin.ReleaseMode {
		// 设置gin的请求日志
		ginLogFile := genv.LogPath() + "/" + genv.AppName() + "/gin-http" + "/%Y-%m-%d.log"
		var cfg = *glogs.DefaultConfig
		cfg.NoBuffWrite = true
		gin.DefaultWriter = glogs.GetWriter(ginLogFile, &cfg)
	}
}

// set app logger
func SetLogger(app *App) {
	glogs.InitLog(
		glogs.SetLogAppName(genv.AppName()),
		glogs.SetLogAppMode(genv.AppMode()),
		glogs.SetLogType(genv.LogType()),
	)
}

// set app logger noBuffer
func SetLoggerNoBuffer(app *App) {
	glogs.InitLog(
		glogs.SetLogAppName(genv.AppName()),
		glogs.SetLogAppMode(genv.AppMode()),
		glogs.SetLogType(genv.LogType()),
		glogs.SetNoBuffWriter(),
	)
}

// set web server
func SetWebServer(app *App) {
	gin.SetMode(genv.RunMode())
	app.WebServer = gin.Default()
}

// open pprof
func SetPprof(app *App) {
	if genv.Pprof() {
		gpprof.InitPprof()
	}
}
