package laya

import (
	"github.com/layasugar/laya/grpc"
	"net"
	"time"
)

var (
	_ Context = &PbRPCContext{}
)

// PbRPCHandlerFunc handler for PbRPC
type PbRPCHandlerFunc func(*PbRPCContext)

// PbRPCContext pbrpc context
type PbRPCContext struct {
	server *PbRPCServer
	conn   net.Conn
	index  int8

	*LogContext
	*MemoryContext
}

// NewPbRPCContext new
func NewPbRPCContext() *PbRPCContext {
	return &PbRPCContext{
		LogContext:    nil,
		MemoryContext: NewMemoryContext(),
	}
}

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *PbRPCContext) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
func (c *PbRPCContext) Done() <-chan struct{} {
	return nil
}

// Err returns a non-nil error value after Done is closed,
// successive calls to Err return the same error.
// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
func (c *PbRPCContext) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *PbRPCContext) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}

// Next Goto next handler
func (c *PbRPCContext) Next() {
	if int(c.index+1) < len(c.server.handlers) {
		c.index++
		c.server.handlers[c.index](c)
	}
}

// ReadPackage return a package
func (c *PbRPCContext) ReadPackage() (*grpc.Package, error) {
	c.conn.SetReadDeadline(time.Now().Add(c.server.ReadTimeout))
	pkg := grpc.NewPackage()
	err := pkg.ReadIO(c.conn)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}

// WritePackage write a package
func (c *PbRPCContext) WritePackage(pkg *grpc.Package) (int, error) {
	c.conn.SetWriteDeadline(time.Now().Add(c.server.WriteTimeout))
	return pkg.WriteIO(c.conn)
}
