package laya

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya/gconf"
	"github.com/layatips/laya/genv"
	"github.com/layatips/laya/glogs"
	"github.com/layatips/laya/gmiddleware"
	"io"
	"log"
	"os"
	"path"
)

type App struct {
	WebServer *gin.Engine
}

func NewApp() *App {
	return new(App).InitWithConfig()
}

func (app *App) InitWithConfig() *App {
	var configPath string
	flag.StringVar(&configPath, "config_path", "", "配置文件地址：xx/xx/app.json")
	flag.Parse()
	if configPath == "" {
		configPath = "./conf/app.json"
	}
	err := gconf.InitConfig(configPath)
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
			GinLog:     "/home/logs/app/default-app/gin_http.log",
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
	if cf.GinLog != "" {
		genv.SetGinLog(cf.GinLog)
	}
	if cf.ParamLog {
		genv.SetParamLog(cf.ParamLog)
	}
	gin.SetMode(genv.RunMode())
	app.WebServer = gin.Default()
	if len(DefaultWebServerMiddlewares) > 0 {
		app.WebServer.Use(DefaultWebServerMiddlewares...)
	}
	if genv.RunMode() == "debug" {
		err := os.MkdirAll(path.Dir(genv.GinLog()), os.ModeDir)
		if err != nil {
			log.Printf("[store_gin_log] Could not create log path")
		}
		logfile, err := os.Create(genv.GinLog())
		if err != nil {
			log.Printf("[store_gin_log] Could not create log file")
		}
		gin.DefaultWriter = io.MultiWriter(logfile)
	}

	glogs.InitLog()
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

var DefaultWebServerMiddlewares = []gin.HandlerFunc{
	gmiddleware.SetHeader,
}
