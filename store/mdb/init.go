package mdb

const mongoConfKey = "mongo"

func init() {
	var mdbs = gconf.GetConfigMap(mongoConfKey)
	InitConn(mdbs)
}
