package server_tpl

const RoutesTestTpl = `package routes

import (
	"github.com/layasugar/laya"

	"{{.goModName}}/controllers/test"
)

// RegisterHttpTest 注册一组http路由
func RegisterHttpTest(r *laya.WebServer) {
	r.POST("/trace-http-test", test.Ctrl.HttpTraceTest)
}
`
