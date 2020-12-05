package config

import (
	"github.com/LaYa-op/laya/i18n"
	"github.com/LaYa-op/laya/logger"
	"github.com/LaYa-op/laya/store/db"
	"github.com/LaYa-op/laya/store/redis"
)

var Path = "./config/app.toml"

type AppConfig struct {
	AppName    string         `toml:"app_name"`
	RunMode    string         `toml:"run_mode"`
	HTTPListen string         `toml:"http_listen"`
	AppVersion string         `toml:"version"`
	DBConfig   *db.Config     `toml:"mysql"`
	RDBConfig  *redis.Config  `toml:"redis"`
	I18nConfig *i18n.Config   `toml:"i18n"`
	LogConfig  *logger.Config `toml:"log"`
}
