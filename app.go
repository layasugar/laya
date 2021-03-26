package laya

import (
	"flag"
	"fmt"
	"github.com/layatips/laya/gconf"
	"github.com/layatips/laya/genv"
	"github.com/layatips/laya/glogs"
	"github.com/layatips/laya/gi18n"
	"github.com/layatips/laya/gstore"
	"github.com/gin-gonic/gin"
)

type App struct {
	webServer *WebServer
}

func NewApp() *App {
	return new(App).InitWithConfig()
}

func (app *App) InitWithConfig() *App {
	var configPath string
	flag.StringVar(&configPath, "config_path", "", "配置文件地址：xx/xx/app.toml")
	flag.Parse()
	err := conf.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	cf := conf.GetBaseConf()
	if cf.AppName != "" {
		genv.SetAppName(cf.AppName)
	}
	if cf.RunMode != "" {
		genv.SetRunMode(cf.RunMode)
	}
	if cf.HttpListen != "" {
		app.webServer = NewWebServer(cf.RunMode)
		if len(DefaultWebServerMiddlewares) > 0 {
			app.webServer.Use(DefaultWebServerMiddlewares...)
		}
	}

	glogs.InitLog()
	gi18n.Init()
	gstore.InitDB()
	gstore.InitMdb()
	gstore.InitRdb()
	gstore.InitMemory()
	fmt.Printf("[app.InitLog] inited with: root_path=%s, config_dir=%s, app_name=%s, run_mode=%s\n",
		genv.RootPath(), genv.ConfRootPath(), genv.AppName(), genv.RunMode())

	return app
}

func (app *App) WebServer() *WebServer {
	return app.webServer
}

func (app *App) RunWebServer() {
	cf := conf.GetBaseConf()
	err := app.webServer.Run(cf.HttpListen)
	if err != nil {
		fmt.Printf("Can't RunWebServer: %s\n", err.Error())
	}
}

var DefaultWebServerMiddlewares []gin.HandlerFunc
