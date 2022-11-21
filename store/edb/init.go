package edb

const esConfKey = "es"

func init() {
	var edbs = gconf.GetConfigMap(esConfKey)

	InitConn(edbs)
}
