// rdbx 不考虑并发读写问题, 项目初始化就会初始化连接, 后面只会取
// map并发读没有问题

package rdbx

import (
	"github.com/go-redis/redis/v8"
)

const (
	defaultDbName = "default"
)

var dbPool map[string]*redis.Client

func getRDB(name string) *redis.Client {
	pool := getPool()
	return pool[name]
}

func setRDB(databaseName string, db *redis.Client) {
	pool := getPool()
	pool[databaseName] = db
}

func getPool() map[string]*redis.Client {
	if dbPool == nil {
		dbPool = make(map[string]*redis.Client)
	}
	return dbPool
}
