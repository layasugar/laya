package laya

import (
	"github.com/LaYa-op/laya/store/db"
	"github.com/LaYa-op/laya/store/redis"
)

func Init() {
	db.Init()
	//mongodb.Init()
	redis.Init()
}
