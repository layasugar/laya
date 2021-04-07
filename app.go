package laya

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya/gconf"
	"github.com/layatips/laya/genv"
	"github.com/layatips/laya/glogs"
	"github.com/layatips/laya/gmiddleware"
	"log"
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

	cf := gconf.GetBaseConf()
	if cf.AppName != "" {
		genv.SetAppName(cf.AppName)
	}
	if cf.RunMode != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.HttpListen != "" {
		app.WebServer = gin.Default()
		if len(DefaultWebServerMiddlewares) > 0 {
			app.WebServer.Use(DefaultWebServerMiddlewares...)
		}
	}
	log.Printf("%s %s %s starting at %q\n", cf.AppName, cf.RunMode, cf.AppUrl, cf.HttpListen)
	glogs.InitLog()
	return app
}

//func (app *App) WebServer() *gin.Engine {
//	return app.webServer
//}

func (app *App) RunWebServer() {
	cf := gconf.GetBaseConf()
	err := app.WebServer.Run(cf.HttpListen)
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
