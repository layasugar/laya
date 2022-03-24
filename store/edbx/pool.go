// edbx 不考虑并发读写问题, 项目初始化就会初始化连接, 后面只会取
// map并发读没有问题

package edbx

import (
	"github.com/elastic/go-elasticsearch/v7"
)

const (
	defaultEdbName = "default"
)

var p map[string]*elasticsearch.Client

func getEdb(databaseName string) *elasticsearch.Client {
	pool := getPool()
	return pool[databaseName]
}

func setEdb(databaseName string, db *elasticsearch.Client) {
	pool := getPool()
	pool[databaseName] = db
}

func getPool() map[string]*elasticsearch.Client {
	if p == nil {
		p = make(map[string]*elasticsearch.Client)
	}
	return p
}
