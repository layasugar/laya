package laya

import (
	"net"
	"runtime"
	"sync"
)

// PbRPCServer struct
type PbRPCServer struct {
	Server
	contextPool sync.Pool
	handlers    []PbRPCHandlerFunc
}

// NewPbRPCServer create new PbRPCServer with default configuration
func NewPbRPCServer() *PbRPCServer {
	server := &PbRPCServer{}

	server.contextPool = sync.Pool{
		New: func() interface{} {
			return &PbRPCContext{
				server: server,
			}
		},
	}
	server.RequestHandler = server.connHandle
	server.MaxWorker = uint32(runtime.NumCPU() * 64)
	server.SetReadTimeout(1500)
	server.SetWriteTimeout(1500)

	return server
}

// AddHandler 添加处理函数
func (pbs *PbRPCServer) AddHandler(handlers ...PbRPCHandlerFunc) {
	pbs.handlers = append(pbs.handlers, handlers...)
}

func (pbs *PbRPCServer) connHandle(conn net.Conn) {
	context := pbs.contextPool.Get().(*PbRPCContext)
	context.conn = conn
	context.LogContext = NewLogContext(context.req)
	for int(context.index) < len(pbs.handlers) {
		pbs.handlers[context.index](context)
		context.index++
	}
	context.conn = nil
	context.index = 0
	context.LogContext = nil
	pbs.contextPool.Put(context)
}
