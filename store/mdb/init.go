package mdb

import "github.com/layasugar/laya/gcnf"

const mongoConfKey = "mongo"

func init() {
	var mdbs = gcnf.GetConfigMap(mongoConfKey)
	InitConn(mdbs)
}
