package gstore

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/layatips/laya/gconf"
	"log"
	"time"
)

// 初始化redis
func InitRdb(cf *gconf.RdbConf) *redis.Client {
	return connRdb(cf.DB, cf.PoolSize, cf.MaxRetries, cf.IdleTimeout, cf.Addr, cf.Pwd)
}

func connRdb(db, poolSize, maxRetries, idleTimeout int, addr, pwd string) *redis.Client {
	options := redis.Options{
		Addr:        addr,                                     // Redis地址
		DB:          db,                                       // Redis库
		PoolSize:    poolSize,                                 // Redis连接池大小
		MaxRetries:  maxRetries,                               // 最大重试次数
		IdleTimeout: time.Second * time.Duration(idleTimeout), // 空闲链接超时时间
	}
	if pwd != "" {
		options.Password = pwd
	}
	Rdb := redis.NewClient(&options)
	pong, err := Rdb.Ping(context.Background()).Result()
	if err == redis.Nil {
		log.Printf("[store_redis] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		log.Printf("[store_redis] redis connRdb err,err=%s\n", err)
		panic(err)
	} else {
		log.Printf("[store_redis] redis connRdb success,suc=%s\n", pong)
	}
	return Rdb
}

func RdbSurvive(rdb *redis.Client) error {
	err := rdb.Ping(context.Background()).Err()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
