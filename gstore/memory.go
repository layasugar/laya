package gstore

import (
	"github.com/layatips/laya/gcache"
	"log"
	"time"
)

var GCache *gcache.Cache

func InitMemory() {
	GCache = gcache.New(0, 1000*time.Minute)
	log.Printf("[store_memory] memory init success")
}
