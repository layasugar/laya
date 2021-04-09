package laya

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya/gconf"
	"github.com/layatips/laya/genv"
	"github.com/layatips/laya/glogs"
	"github.com/layatips/laya/gmiddleware"
	"log"
	"path/filepath"
	"sync"
)

type App struct {
	WebServer *gin.Engine
	watcher   *gconf.Watcher
	watchLock sync.Mutex
}

func NewApp() *App {
	return new(App).InitWithConfig()
}

func (app *App) InitWithConfig() *App {
	err := gconf.InitConfig(genv.ConfigPath)
	if err != nil {
		panic(err)
	}

	cf, err := gconf.GetBaseConf()
	if err != nil && !errors.Is(err, gconf.Nil) {
		panic(err.Error())
	}
	if errors.Is(err, gconf.Nil) {
		cf = &gconf.BaseConf{
			AppName:    "default-app",
			HttpListen: "0.0.0.0:10080",
			RunMode:    "debug",
			AppVersion: "1.0.0",
			AppUrl:     "127.0.0.1:10080",
			ParamLog:   true,
		}
	}
	if cf.AppName != "" {
		genv.SetAppName(cf.AppName)
	}
	if cf.HttpListen != "" {
		genv.SetHttpListen(cf.HttpListen)
	}
	if cf.RunMode != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.AppVersion != "" {
		genv.SetAppVersion(cf.AppVersion)
	}
	if cf.AppUrl != "" {
		genv.SetAppUrl(cf.AppUrl)
	}
	genv.SetParamLog(cf.ParamLog)

	gin.SetMode(genv.RunMode())
	glogs.InitLog()
	app.WebServer = gin.Default()
	if len(DefaultWebServerMiddlewares) > 0 {
		app.WebServer.Use(DefaultWebServerMiddlewares...)
	}

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

var DefaultWebServerMiddlewares = []gin.HandlerFunc{
	gmiddleware.SetHeader,
	gmiddleware.LogInParams,
}
