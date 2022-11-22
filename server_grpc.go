package laya

import (
	"context"
	"log"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/layasugar/laya/core/constants"
	"github.com/layasugar/laya/core/metautils"
)

// GrpcServer struct
type GrpcServer struct {
	*grpc.Server
	opts   []grpc.UnaryServerInterceptor
	routes []func(server *GrpcServer)
}

// NewGrpcServer create new GrpcServer with default configuration
func NewGrpcServer() *GrpcServer {
	server := &GrpcServer{
		opts: []grpc.UnaryServerInterceptor{
			serverInterceptor,
		},
	}

	return server
}

func (gs *GrpcServer) Use(f ...func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)) {
	for _, vf := range f {
		gs.opts = append(gs.opts, vf)
	}
}

func (gs *GrpcServer) Register(f ...func(s *GrpcServer)) {
	gs.routes = append(gs.routes, f...)
}

func (gs *GrpcServer) Run(addr string) (err error) {
	// 初始化server, 将多个拦截器构建成一个拦截器
	gs.Server = grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(gs.opts...)),
	)

	// 注册路由
	for _, vf := range gs.routes {
		vf(gs)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	// 在给定的gRPC服务器上注册服务器反射服务
	reflection.Register(gs.Server)

	// Serve方法在lis上接受传入连接，为每个连接创建一个ServerTransport和server的goroutine。
	// 该goroutine读取gRPC请求，然后调用已注册的处理程序来响应它们
	err = gs.Server.Serve(lis)
	return
}

// serverInterceptor 提供服务的拦截器, 重写context, 记录出入参, 记录链路追踪
func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 初始化context
	md := metautils.ExtractIncoming(ctx)
	newCtx := NewContext(constants.SERVERGRPC, info.FullMethod, md, nil)

	// 入参 header->meta
	if env.ApiLog() {
		reqByte, _ := tools.CJson.Marshal(req)
		mdByte, _ := tools.CJson.Marshal(md)
		newCtx.InfoF("%s", string(reqByte),
			newCtx.Field("header", string(mdByte)),
			newCtx.Field("path", info.FullMethod),
			newCtx.Field("protocol", protocol),
			newCtx.Field("title", "入参"))
	}

	resp, err := handler(newCtx, req)

	if env.ApiLog() {
		respByte, _ := tools.CJson.Marshal(resp)
		newCtx.InfoF("%s", string(respByte), newCtx.Field("title", "出参"))
	}
	newCtx.SpanFinish(newCtx.TopSpan)
	return resp, err
}
