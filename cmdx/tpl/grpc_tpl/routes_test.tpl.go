package grpc_tpl

const RoutesTestTpl = `package routes

import (
	"github.com/layasugar/laya"

	"{{.goModName}}/controllers/test"
	"{{.goModName}}/pb"
)

// RegisterRpcRoutes 注册一组rpc路由
func RegisterRpcRoutes(s *laya.GrpcServer) {
	pb.RegisterGreeterServer(s.Server, test.Ctrl)
}`
