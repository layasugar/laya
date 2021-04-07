package gstore

import (
	"errors"
	"github.com/layatips/laya/gcache"
	"github.com/layatips/laya/gconf"
	"log"
	"time"
)

func InitMemory() *gcache.Cache {
	memc, err := gconf.GetMemConf()
	if err != nil && !errors.Is(err, gconf.Nil) {
		panic(err.Error())
	}
	if errors.Is(err, gconf.Nil) {
		memc = &gconf.MemConf{
			DefaultExp: 0,
			Cleanup:    600,
		}
	}
	Mem := gcache.New(time.Duration(memc.DefaultExp)*time.Second, time.Duration(memc.Cleanup)*time.Second)
	log.Printf("[store_memory] memory init success")
	return Mem
}
