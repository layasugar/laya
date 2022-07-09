package rdb

import "github.com/layasugar/laya/gcnf"

const redisConfKey = "redis"

func init() {
	var rdbs = gcnf.GetConfigMap(redisConfKey)

	InitConn(rdbs)
}
