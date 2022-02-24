package mdbx

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/layasugar/laya/core/appx"
	"github.com/layasugar/laya/core/grpcx"
	"github.com/layasugar/laya/core/httpx"
	"github.com/layasugar/laya/core/tracex"
	"go.mongodb.org/mongo-driver/event"
	"sync"
)

const (
	tSpanName = "mongo"
)

type tracer struct {
	spans sync.Map
}

func NewTracer() *tracer {
	return &tracer{}
}

const prefix = "mongodb."

func (t *tracer) HandleStartedEvent(ctx context.Context, evt *event.CommandStartedEvent) {
	if evt == nil {
		return
	}

	var traceCtx *tracex.TraceContext
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
		if defaultCtx, okInterface := ctx.(*appx.Context); okInterface {
			traceCtx = defaultCtx.TraceContext
		}
	}

	if traceCtx != nil {
		span := traceCtx.SpanStart(tSpanName)
		if nil != span {
			ext.DBType.Set(span, "mongo")
			ext.DBInstance.Set(span, evt.DatabaseName)
			ext.DBStatement.Set(span, evt.Command.String())
			span.SetTag("db.host", evt.ConnectionID)
			ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
			ext.Component.Set(span, "golang-mongo")
		}
		t.spans.Store(evt.RequestID, span)
	}
}

func (t *tracer) HandleSucceededEvent(ctx context.Context, evt *event.CommandSucceededEvent) {
	if evt == nil {
		return
	}
	if rawSpan, ok := t.spans.Load(evt.RequestID); ok {
		defer t.spans.Delete(evt.RequestID)
		if span, ok := rawSpan.(opentracing.Span); ok {
			defer span.Finish()
			//span.SetTag(prefix+"reply", string(evt.Reply))
			span.SetTag(prefix+"duration", evt.DurationNanos)
		}
	}
}

func (t *tracer) HandleFailedEvent(ctx context.Context, evt *event.CommandFailedEvent) {
	if evt == nil {
		return
	}
	if rawSpan, ok := t.spans.Load(evt.RequestID); ok {
		defer t.spans.Delete(evt.RequestID)
		if span, ok := rawSpan.(opentracing.Span); ok {
			defer span.Finish()
			ext.Error.Set(span, true)
			span.SetTag(prefix+"duration", evt.DurationNanos)
			span.LogFields(log.String("failure", evt.Failure))
		}
	}
}
