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
	return NewApp(SetLogger, SetGinLog, SetDing, SetTrace, SetWebServer, SetPprof)
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
		cf = &gconf.BaseConf{
			AppName:    "default-app",
			HttpListen: "0.0.0.0:10080",
			RunMode:    "debug",
			AppVersion: "1.0.0",
			AppUrl:     "127.0.0.1:10080",
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
func (app *App) registerEnv(cf *gconf.BaseConf) {
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
	if cf.LogPath != "" {
		genv.SetLogPath(cf.LogPath)
	}
	if cf.AppMode != "" {
		genv.SetAppMode(cf.AppMode)
	}
	if cf.Pprof {
		genv.SetPprof(cf.Pprof)
	}
	if genv.RunMode() == "release" {
		genv.SetLogType("file")
	}
	genv.SetParamLog(cf.ParamLog)
}

// Check trace it's on or not
func SetTrace(app *App) {
	tc, err := gconf.GetTraceConf()
	if errors.Is(err, gconf.Nil) {
		return
	}
	if err != nil {
		log.Printf("trace配置获取失败,err=%s", err.Error())
		return
	}
	if tc == nil {
		return
	} else {
		if tc.ZipkinAddr == "" {
			return
		}
		err = glogs.InitTrace(genv.AppName(), genv.HttpListen(), tc.ZipkinAddr, tc.Mod)
		if err != nil {
			log.Printf("trace初始化失败")
			return
		}
	}
}

// set ding ding pusher
func SetDing(app *App) {
	dc, err := gconf.GetDingConf()
	if errors.Is(err, gconf.Nil) {
		return
	}
	if err != nil {
		log.Printf("trace配置获取失败,err=%s", err.Error())
		return
	}
	if dc == nil {
		return
	} else {
		if dc.RobotKey == "" || dc.RobotHost == "" {
			return
		}
		glogs.InitDing(dc.RobotKey, dc.RobotHost)
	}
}

// set gin logger
func SetGinLog(app *App) {
	if genv.AppMode() == "release" {
		// 设置gin的请求日志
		ginLogFile := genv.LogPath() + "/" + genv.AppName() + "/gin-http" + "/%Y-%m-%d.log"
		var cfg = glogs.LogConfig{
			RotationSize: 64 * 1024 * 1024,
			NoBuffWrite:  true,
		}
		gin.DefaultWriter = glogs.GetWriter(
			ginLogFile,
			&cfg,
		)
	}
}

// set gin logger noBuffer
func SetGinLogNoBuffer(app *App) {
	if genv.AppMode() == gin.ReleaseMode {
		// 设置gin的请求日志
		ginLogFile := genv.LogPath() + "/" + genv.AppName() + "/gin-http" + "/%Y-%m-%d.log"
		var cfg = glogs.LogConfig{
			RotationSize: 64 * 1024 * 1024,
			NoBuffWrite:  true,
		}
		gin.DefaultWriter = glogs.GetWriter(
			ginLogFile,
			&cfg,
		)
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
