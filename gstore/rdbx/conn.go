package rdbx

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

const (
	defaultRedisPoolMinIdle = 2 // 连接池空闲连接数量
)

// InitRdb 初始化redis
func InitRdb(cfg redis.Options) *redis.Client {
	return connRdb(cfg)
}

func connRdb(options redis.Options) *redis.Client {
	if options.MinIdleConns == 0 {
		options.MinIdleConns = defaultRedisPoolMinIdle
	}
	Rdb := redis.NewClient(&options)
	_, err := Rdb.Ping(context.Background()).Result()
	if err == redis.Nil {
		log.Printf("[app.rdbx] Nil reply returned by Rdb when key does not exist.")
	} else if err != nil {
		log.Printf("[app.rdbx] redis fail, err: %s", err)
		panic(err)
	} else {
		log.Printf("[app.rdbx] redis success")
	}
	Rdb.AddHook(NewHook())
	return Rdb
}

// RdbSurvive redis存活检测
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

func InitConn(m []map[string]interface{}) {
	for _, item := range m {
		var cfg redis.Options
		var name string

		if nameIf, ok := item["name"]; ok {
			if nameStr, okInterface := nameIf.(string); okInterface {
				name = nameStr
			}
		}

		if addr, ok := item["addr"]; ok {
			if addrStr, okInterface := addr.(string); okInterface {
				cfg.Addr = addrStr
			}
		}

		if db, ok := item["db"]; ok {
			if dbInt, okInterface := db.(int64); okInterface {
				cfg.DB = int(dbInt)
			}
		}

		if pwd, ok := item["pwd"]; ok {
			if pwdStr, okInterface := pwd.(string); okInterface {
				cfg.Password = pwdStr
			}
		}

		if name == "" {
			continue
		}

		setRDB(name, connRdb(cfg))
	}
}

func GetClient(name ...string) *redis.Client {
	if len(name) > 0 {
		return getRDB(name[0])
	} else {
		return getRDB(defaultDbName)
	}
}
