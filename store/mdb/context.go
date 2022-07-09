package mdb

import (
	"context"
	"github.com/layasugar/laya/env"
	"github.com/layasugar/laya/store/cm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
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
	if env.MongoTrace() {
		span := cm.ParseSpanByCtx(ctx, tSpanName)

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
	if env.MongoTrace() {
		if rawSpan, ok := t.spans.Load(evt.RequestID); ok {
			defer t.spans.Delete(evt.RequestID)
			if span, ok := rawSpan.(opentracing.Span); ok {
				defer span.Finish()
				//span.SetTag(prefix+"reply", string(evt.Reply))
				span.SetTag(prefix+"duration", evt.DurationNanos)
			}
		}
	}
}

func (t *tracer) HandleFailedEvent(ctx context.Context, evt *event.CommandFailedEvent) {
	if evt == nil {
		return
	}
	if env.MongoTrace() {
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
}
