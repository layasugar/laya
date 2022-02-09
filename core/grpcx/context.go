package grpcx

import (
	"context"
	"github.com/layasugar/laya/core/logx"
)

// HandlerFunc 请求的处理者
type HandlerFunc func(ctx Context)

// Context is the carrier of request and response
type Context interface {
	context.Context
	DataContext
	logx.LoggerContext
}
