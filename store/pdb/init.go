package db

const mysqlConfKey = "mysql"

func init() {
	// 初始化数据库连接和redis连接
	var dbs = gconf.GetConfigMap(mysqlConfKey)

	// 解析dbs
	InitConn(dbs)
}
