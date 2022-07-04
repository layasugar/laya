package cm

import (
	"context"
	"github.com/layasugar/laya/core/appx"
	"github.com/layasugar/laya/core/grpcx"
	"github.com/layasugar/laya/core/httpx"
	"github.com/opentracing/opentracing-go"
)

// ParseSpanByCtx 公共方法, 从ctx中获取
func ParseSpanByCtx(ctx interface{}, spanName string) opentracing.Span {
	var traceCtx *tracer.TraceContext
	switch ctx.(type) {
	case *httpx.WebContext:
		if webCtx, okInterface := ctx.(*httpx.WebContext); okInterface {
			traceCtx = webCtx.TraceContext
		}
	case *grpcx.GrpcContext:
		if grpcCtx, okInterface := ctx.(*grpcx.GrpcContext); okInterface {
			traceCtx = grpcCtx.TraceContext
		}
	case *appx.Context:
		if appCtx, okInterface := ctx.(*appx.Context); okInterface {
			traceCtx = appCtx.TraceContext
		}
	}
	if nil != traceCtx {
		return traceCtx.SpanStart(spanName)
	}
	return nil
}

// ParseLogIdByCtx 从context中解析出logId
func ParseLogIdByCtx(ctx context.Context) string {
	var requestId string
	switch ctx.(type) {
	case *httpx.WebContext:
		if webCtx, okInterface := ctx.(*httpx.WebContext); okInterface {
			requestId = webCtx.LogID()
		}
	case *grpcx.GrpcContext:
		if grpcCtx, okInterface := ctx.(*grpcx.GrpcContext); okInterface {
			requestId = grpcCtx.GetLogId()
		}
	case *appx.Context:
		if appCtx, okInterface := ctx.(*appx.Context); okInterface {
			requestId = appCtx.GetLogId()
		}
	}
	return requestId
}
