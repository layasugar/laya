package laya

import (
	"encoding/json"
	"github.com/LaYa-op/laya/config"
	"github.com/LaYa-op/laya/store/db"
	"github.com/LaYa-op/laya/store/mongodb"
	"github.com/LaYa-op/laya/store/redis"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2/util/log"
)

func init() {
	InitEnv()
	db.Init("./config/db")
	mongodb.Init()
	redis.Init()
}

func InitEnv() {
	_, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	err = config.LoadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(config.Get(ENV, "i18n").Bytes(), &I18n.Conf)
	log.Info(I18n.Conf)
	if err != nil {
		panic(err)
	}

	// get mysql config
	err = json.Unmarshal(config.Get(ENV, "database").Bytes(), &MysqlConf)
	if err != nil {
		panic(err)
	}

	// get cache config
	err = json.Unmarshal(config.Get(ENV, "cache").Bytes(), &RedisConf)
	if err != nil {
		panic(err)
	}

	// get delayerServer config
	DelayServer = config.Get(ENV, "delayServer").String("http://127.0.0.1:9278")
}
