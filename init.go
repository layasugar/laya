package laya

import (
	log "github.com/LaYa-op/laya/logger"
	"github.com/LaYa-op/laya/store/db"
	"github.com/LaYa-op/laya/store/redis"
)

func Init() {
	db.Init()
	redis.Init()
	log.Init()
}
