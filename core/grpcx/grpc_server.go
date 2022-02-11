package grpcx

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

// GrpcServer struct
type GrpcServer struct {
	*grpc.Server
	contextPool sync.Pool
}

// NewGrpcServer create new GrpcServer with default configuration
func NewGrpcServer() *GrpcServer {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(serverInterceptor),
	}

	server := &GrpcServer{
		Server: grpc.NewServer(opts...),
	}

	server.contextPool = sync.Pool{
		New: func() interface{} {
			return &GrpcContext{
				server: server,
			}
		},
	}

	return server
}

func (gs *GrpcServer) Register(f func(s *GrpcServer)) {
	f(gs)
}

func (gs *GrpcServer) Run(addr string) (err error) {
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
