package grpc_tpl

const ControllersTestTestTpl = `package test

import (
	"context"
	"errors"
	"github.com/layasugar/laya"

	"{{.goModName}}/models/page/test"
	"{{.goModName}}/pb"
)

// SayHello 基础测试使用
func (ctrl *controller) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: in.Name}, nil
}

// GrpcTraceTest 测试http请求和链路追踪(grpc_to_http grpc_to_grpc)
func (ctrl *controller) GrpcTraceTest(ctx context.Context, in *pb.GrpcTraceTestReq) (*pb.HelloReply, error) {
	// 转换ctx
	newCtx := ctx.(*laya.GrpcContext)

	// 参数验证
	if in.Kind == 0 {
		return nil, errors.New("请传入kind")
	}

	// 业务处理
	resp, err := test.RpcTraceTest(newCtx, in)
	if err != nil {
		return nil, err
	}

	// 响应
	return &pb.HelloReply{Message: resp.Code}, nil
}
`
