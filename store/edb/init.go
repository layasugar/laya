package edb

import "github.com/layasugar/laya/gcnf"

const esConfKey = "es"

func init() {
	var edbs = gcnf.GetConfigMap(esConfKey)

	InitConn(edbs)
}
