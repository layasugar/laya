package rdb

const redisConfKey = "redis"

func init() {
	var rdbs = gconf.GetConfigMap(redisConfKey)

	InitConn(rdbs)
}
