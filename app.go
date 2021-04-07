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
			ParamsLog:  true,
		}
	}
	if cf.AppName != "" {
		genv.SetAppName(cf.AppName)
	}
	if cf.HttpListen != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.RunMode != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.AppVersion != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.AppUrl != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.GinLog != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.ParamsLog {
		genv.SetRunMode(cf.RunMode)
	}

	log.Printf("%s %s %s starting at %q\n", cf.AppName, cf.RunMode, cf.AppUrl, cf.HttpListen)
	glogs.InitLog()
	return app
}

//func (app *App) WebServer() *gin.Engine {
//	return app.webServer
//}

func (app *App) RunWebServer() {
	app.WebServer = gin.Default()
	if len(DefaultWebServerMiddlewares) > 0 {
		app.WebServer.Use(DefaultWebServerMiddlewares...)
	}
	cf := gconf.GetBaseConf()
	err := app.WebServer.Run(cf.HttpListen)
	if genv.RunMode() == "debug" {
		err := os.MkdirAll(path.Dir(cf.GinLog), os.ModeDir)
		if err != nil {
			log.Printf("[store_gin_log] Could not create log path")
		}
		logfile, err := os.Create(cf.GinLog)
		if err != nil {
			log.Printf("[store_gin_log] Could not create log file")
		}
		gin.DefaultWriter = io.MultiWriter(logfile)
	}
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
