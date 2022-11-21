package http_tpl

const MainTpl = `package main

import (
	"github.com/layasugar/laya"

	"{{.goModName}}/middlewares"
	"{{.goModName}}/routes"
)

// webAppSetup 初始化服务设置
func webAppSetup() *laya.App {
	app := laya.WebApp()

	// register global middlewares
	app.WebServer().Use(middlewares.TestMiddleware())

	// register routes
	app.WebServer().Register(routes.RegisterHttpTest)

	// 屏蔽不需要打印出入参路由分组
	//app.SetNoLogParamsPrefix("/admin")

	return app
}

func main() {
	app := webAppSetup()

	app.RunServer()
}
`
