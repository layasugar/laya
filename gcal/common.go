package gcal

import (
	"context"
	"github.com/layasugar/laya/core/appx"
	"github.com/layasugar/laya/core/grpcx"
	"github.com/layasugar/laya/core/httpx"
	"github.com/layasugar/laya/core/metautils"
	"github.com/layasugar/laya/gcal/converter"
	"github.com/layasugar/laya/gcal/pool"
	"github.com/layasugar/laya/gcal/protocol"
	"github.com/layasugar/laya/gcal/service"
	"github.com/layasugar/laya/tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
)

var pbTc = &pool.Pool{}

// HTTPRequest 别名
type HTTPRequest = protocol.HTTPRequest

// HTTPHead 别名
type HTTPHead = protocol.HTTPHead

// ConverterType 别名
type ConverterType = converter.ConverterType

// JSONConverter 别名
var JSONConverter = converter.JSON

// FORMConverter 别名
var FORMConverter = converter.FORM

// RAWConverter 别名
var RAWConverter = converter.RAW

// LoadService load one service from struct
func LoadService(configs []map[string]interface{}) error {
	return service.LoadService(configs)
}

func GetRpcConn(serverName string) *grpc.ClientConn {
	srv, ok := service.GetService(serverName)
	if !ok {
		return nil
	}

	curConnKey := pool.Key{
		Schema: "tcp",
		Addr:   srv.GetAddr(),
	}
	tcConn, _ := pbTc.Get(curConnKey)
	if tcConn == nil {
		conn, errDial := grpc.Dial(srv.GetAddr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(clientInterceptor),
		)
		if errDial != nil {
			return nil
		}
		//c := pool.Func{
		//	Factory: func() (interface{}, error) {
		//		return grpc.Dial(srv.GetAddr(),
		//			grpc.WithTransportCredentials(insecure.NewCredentials()),
		//			grpc.WithUnaryInterceptor(clientInterceptor),
		//		)
		//	},
		//}
		//pbTc.SetFunc(curConnKey, c)
		defer pbTc.Put(curConnKey, conn)
		return conn
	}
	conn := tcConn.(*grpc.ClientConn)
	return conn
}

// clientInterceptor 提供客户端的拦截器, 注入trace, 注入logId
func clientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var x = make(map[string][]string)
	// 反射ctx, 判断是webContext, 还是grpcContext
	if oldCtx, ok := ctx.(*httpx.WebContext); ok {
		x[tools.RequestIdKey] = []string{oldCtx.LogID()}
		oldCtx.SpanInject(x)
	}

	if oldCtx, ok := ctx.(*grpcx.GrpcContext); ok {
		x[tools.RequestIdKey] = []string{oldCtx.GetLogId()}
		oldCtx.SpanInject(x)
	}

	if oldCtx, ok := ctx.(*appx.Context); ok {
		x[tools.RequestIdKey] = []string{oldCtx.GetLogId()}
		oldCtx.SpanInject(x)
	}

	// 转换key为小写不然rst
	var md = make(metautils.NiceMD)
	for k, v := range x {
		key := strings.ToLower(k)
		if len(v) > 0 {
			md.Set(key, v[0])
		}
	}

	newCtx := md.ToOutgoing(context.Background())
	return invoker(newCtx, method, req, reply, cc, opts...)
}
