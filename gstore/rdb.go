package gstore

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// 初始化redis
func InitRdb(db, poolSize, maxRetries, idleTimeout int, addr, pwd string) *redis.Client {
	return connRdb(db, poolSize, maxRetries, idleTimeout, addr, pwd)
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
	_, err := Rdb.Ping(context.Background()).Result()
	if err == redis.Nil {
		log.Printf("[gstore_redis] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		log.Printf("[gstore_redis] redis fail, err=%s", err)
		panic(err)
	} else {
		log.Printf("[gstore_redis] redis success")
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
