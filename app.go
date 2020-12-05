package laya

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/LaYa-op/laya/config"
	"github.com/LaYa-op/laya/env"
	"github.com/LaYa-op/laya/logger"
	"github.com/LaYa-op/laya/store/db"
	"github.com/LaYa-op/laya/store/redis"
	"github.com/gin-gonic/gin"
)

type App struct {
	config    *config.AppConfig
	webServer *WebServer
}

func NewApp() *App {
	return new(App).Init()
}

func (app *App) Init() *App {
	return app.InitWithConfigName(config.Path)
}

func (app *App) InitWithConfigName(fn string) *App {
	cf := config.AppConfig{}

	if _, err := toml.DecodeFile(fn, &cf); err != nil {
		panic(fmt.Sprintf("Can't load config file %s: %s", fn, err.Error()))
	}

	return app.InitWithConfig(&cf)
}

func (app *App) InitWithConfig(config *config.AppConfig) *App {
	if config == nil {
		panic("Can't initial App with nil config")
	}
	app.config = config

	if app.config.AppName != "" {
		env.SetAppName(app.config.AppName)
	}

	if app.config.RunMode != "" {
		env.SetRunMode(app.config.RunMode)
	}

	if app.config.HTTPListen != "" {
		app.webServer = NewWebServer(app.config.RunMode)
		if len(DefaultWebServerMiddlewares) > 0 {
			app.webServer.Use(DefaultWebServerMiddlewares...)
		}
	}

	app.InitLog()
	app.InitDb()

	logger.ZapLog.Infof("[app.Init] inited with: root_path=%s, config_dir=%s, app_name=%s, run_mode=%s",
		env.RootPath(), env.ConfRootPath(), env.AppName(), env.RunMode())

	return app
}

func (app *App) WebServer() *WebServer {
	return app.webServer
}

func (app *App) RunWebServer() {
	err := app.webServer.Run(app.config.HTTPListen)
	if err != nil {
		logger.ZapLog.Warnf("Can't RunWebServer: %s", err.Error())
	}
}

func (app *App) InitLog() {
	logger.Init(app.config.LogConfig)
}

func (app *App) InitDb() {
	db.Init(app.config.DBConfig)
	redis.Init(app.config.RDBConfig)
}

var DefaultWebServerMiddlewares []gin.HandlerFunc
