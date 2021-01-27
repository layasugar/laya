package laya

import (
	"flag"
	"fmt"
	"github.com/LaYa-op/laya/config"
	"github.com/LaYa-op/laya/env"
	"github.com/LaYa-op/laya/glogs"
	"github.com/LaYa-op/laya/i18n"
	"github.com/LaYa-op/laya/store"
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
	err := config.InitConfig(configPath)
	if err != nil {
		panic(err)
	}

	cf := config.GetBaseConf()
	if cf.AppName != "" {
		env.SetAppName(cf.AppName)
	}
	if cf.RunMode != "" {
		env.SetRunMode(cf.RunMode)
	}
	if cf.HttpListen != "" {
		app.webServer = NewWebServer(cf.RunMode)
		if len(DefaultWebServerMiddlewares) > 0 {
			app.webServer.Use(DefaultWebServerMiddlewares...)
		}
	}

	glogs.Init()
	i18n.Init()
	store.InitDB()
	store.InitMdb()
	store.InitRdb()
	store.InitMemory()
	fmt.Printf("[app.Init] inited with: root_path=%s, config_dir=%s, app_name=%s, run_mode=%s\n",
		env.RootPath(), env.ConfRootPath(), env.AppName(), env.RunMode())

	return app
}

func (app *App) WebServer() *WebServer {
	return app.webServer
}

func (app *App) RunWebServer() {
	cf := config.GetBaseConf()
	err := app.webServer.Run(cf.HttpListen)
	if err != nil {
		fmt.Printf("Can't RunWebServer: %s\n", err.Error())
	}
}

var DefaultWebServerMiddlewares []gin.HandlerFunc
