package store

import (
	"github.com/patrickmn/go-cache"
	"log"
	"time"
)

var ChannelCache *cache.Cache
var ChannelBaseCache *cache.Cache
var ChannelPassageCache *cache.Cache
var ChannelTypeCache *cache.Cache

func InitMemory() {
	ChannelCache = cache.New(0, 1000*time.Minute)
	ChannelBaseCache = cache.New(0, 1000*time.Minute)
	ChannelPassageCache = cache.New(0, 1000*time.Minute)
	ChannelTypeCache = cache.New(0, 1000*time.Minute)
	log.Printf("[store_memory] memory init success")
}
