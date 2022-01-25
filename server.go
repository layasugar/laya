package laya

import (
	"fmt"
	"github.com/layasugar/laya/gserver"
	"go.uber.org/zap"
	"log"
	"net"
	"time"
)

// Server 提供较为通用的服务器程序
type Server struct {
	server gserver.Server

	// Exported fields
	Logger        *zap.Logger
	Mode          string
	MaxWorker     uint32
	PrepareWorker uint32
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration

	// Exported Method needs to apply
	RequestHandler func(con net.Conn)

	rejectHandler      func(net.Conn, error)
	acceptErrorHandler func(error)
}

func (gs *Server) Run(listen ...string) (err error) {
	// if length of length is zero, should search listen address in environment variables
	if len(listen) != 1 {
		return fmt.Errorf("%s", "listen address must be only one")
	}
	gs.server = new(gserver.TCPServer)
	return gs.run("TCP", listen[0])
}

func (gs *Server) RunUnix(path ...string) (err error) {
	if len(path) != 1 {
		return fmt.Errorf("listen address must be only one")
	}
	gs.server = new(gserver.UnixServer)

	return gs.run("Unix", path[0])
}

func (gs *Server) run(serverName string, listen string) error {

	gs.server.SetHandler(gs.RequestHandler)
	gs.server.SetWorkersPoolSize(gs.PrepareWorker, gs.MaxWorker)

	if gs.acceptErrorHandler != nil {
		gs.server.SetAcceptErrorHandler(gs.acceptErrorHandler)
	}
	if gs.rejectHandler != nil {
		gs.server.SetRejectHandler(gs.rejectHandler)
	}

	defer gs.server.Close()
	log.Printf("[app] Listening and serving %s on %s\n", serverName, listen)
	return gs.server.Run(listen)
}

func (gs *Server) RunGrace(addr string, timeouts ...time.Duration) error {
	return nil
}

func (gs *Server) SetLogger(logger *zap.Logger) {
	gs.Logger = logger
}

func (gs *Server) SetReadTimeout(ms int64) {
	gs.ReadTimeout = time.Duration(ms) * time.Millisecond
}

func (gs *Server) SetWriteTimeout(ms int64) {
	gs.WriteTimeout = time.Duration(ms) * time.Millisecond
}

func (gs *Server) SetRejectHandler(h func(net.Conn, error)) {
	gs.rejectHandler = h
}

func (gs *Server) SetAcceptErrorHandler(h func(error)) {
	gs.acceptErrorHandler = h
}

func (gs *Server) CountBusyWorkers() uint32 {
	return gs.server.CountBusyWorkers()
}

func (gs *Server) CountAvailableWorkers() uint32 {
	return gs.server.CountAvailableWorkers()
}
