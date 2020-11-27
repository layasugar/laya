package redis

import (
	"context"
	"github.com/LaYa-op/laya/config"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var Rdb *redis.Client

// 初始化redis
func Init() {
	path := ""
	Configs := config.ListFiles(path)

	log.Printf("[store_db] DB_INIT with %d cluster\n", len(Configs))
	var config Config
	for _, name := range Configs {
		err := config.ReadFile(name, &config)
		if err != nil {
			log.Printf("[store_db] parse db config %s failed,err= %s\n", name, err)
			continue
		}
	}
	if config.Open {
		conn(&config)
	}
}

func conn(conf *Config) {
	options := redis.Options{
		Addr:        conf.Addr,                                     // Redis地址
		DB:          conf.DB,                                       // Redis库
		PoolSize:    conf.PoolSize,                                 // Redis连接池大小
		MaxRetries:  conf.MaxRetries,                               // 最大重试次数
		IdleTimeout: time.Second * time.Duration(conf.IdleTimeout), // 空闲链接超时时间
	}
	if conf.Pwd != "" {
		options.Password = conf.Pwd
	}
	Rdb = redis.NewClient(&options)
	pong, err := Rdb.Ping(context.Background()).Result()
	if err == redis.Nil {
		log.Printf("[store redis] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		log.Printf("[store_db] redis conn err,err= %s\n", err)
		panic(err)
	} else {
		log.Printf("[store_db] redis conn success,suc= %s\n", pong)
	}
}
