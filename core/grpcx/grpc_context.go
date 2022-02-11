package grpcx

import (
	"github.com/layasugar/laya/core/alarmx"
	"github.com/layasugar/laya/core/logx"
	"github.com/layasugar/laya/core/metautils"
	"github.com/layasugar/laya/core/tracex"
	"github.com/layasugar/laya/gtools"
	uuid "github.com/satori/go.uuid"
	"time"
)

var (
	_ Context = &GrpcContext{}
)

// GrpcHandlerFunc handler for PbRPC
type GrpcHandlerFunc func(*GrpcContext)

// GrpcContext pbrpc contextx
type GrpcContext struct {
	server *GrpcServer

	*logx.LogContext
	*MemoryContext
	*tracex.TraceContext
	*alarmx.AlarmContext
}

// NewGrpcContext newCtx
func NewGrpcContext(name string, md metautils.NiceMD) *GrpcContext {
	logId := md.Get(gtools.RequestIdKey)
	if logId == "" {
		logId = gtools.Md5(uuid.NewV4().String())
	}

	ctx := &GrpcContext{
		LogContext:    logx.NewLogContext(logId),
		TraceContext:  tracex.NewTraceContext(name, md),
		MemoryContext: NewMemoryContext(),
	}
	ctx.Set(gtools.RequestIdKey, logId)
	return ctx
}

// Deadline returns the time when work done on behalf of this contextx
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *GrpcContext) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done returns a channel that's closed when work done on behalf of this
// contextx should be canceled. Done may return nil if this contextx can
// never be canceled. Successive calls to Done return the same value.
func (c *GrpcContext) Done() <-chan struct{} {
	return nil
}

// Err returns a non-nil error value after Done is closed,
// successive calls to Err return the same error.
// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the contextx was canceled
// or DeadlineExceeded if the contextx's deadline passed.
func (c *GrpcContext) Err() error {
	return nil
}

// Value returns the value associated with this contextx for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *GrpcContext) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
