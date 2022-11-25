package db

import "github.com/layasugar/laya/gcnf"

const mysqlConfKey = "mysql"

func init() {
	// 初始化数据库连接和redis连接
	var dbs = gcnf.GetConfigMap(mysqlConfKey)

	// 解析dbs
	InitConn(dbs)
}
