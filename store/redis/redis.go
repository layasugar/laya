package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var Rdb *redis.Client

type Config struct {
	Open        bool   `toml:"open"`
	DB          int    `toml:"db"`
	Addr        string `toml:"addr"`
	Pwd         string `toml:"pwd"`
	PoolSize    int    `toml:"poolSize"`
	MaxRetries  int    `toml:"maxRetries"`
	IdleTimeout int    `toml:"idleTimeout"`
}

// 初始化redis
func Init(config *Config) {
	if config.Open {
		conn(config)
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
		log.Printf("[store_redis] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		log.Printf("[store_redis] redis conn err,err=%s\n", err)
		panic(err)
	} else {
		log.Printf("[store_redis] redis conn success,suc=%s\n", pong)
	}
}
