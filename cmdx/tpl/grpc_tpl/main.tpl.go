package grpc_tpl

const MainTpl = `package main

import (
	"github.com/layasugar/laya"
	"{{.goModName}}/models/dao"
	"{{.goModName}}/routes"
)

// grpcAppSetup 初始化服务设置
func grpcAppSetup() *laya.App {
	app := laya.GrpcApp()

	// open db connection and global do before something
	app.Use(dao.Init)

	// rpc 路由
	app.GrpcServer().Register(routes.RegisterRpcRoutes)

	return app
}

func main() {
	app := grpcAppSetup()

	app.RunServer()
}
`
