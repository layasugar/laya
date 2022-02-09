package tracex

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"log"
	"net/http"
)

// TracerContext 链路
type TracerContext interface {
	SpanFinish(span opentracing.Span)

	// SpanStart 开启子span
	SpanStart(name string) opentracing.Span

	// SpanInject 注入请求
	SpanInject(r *http.Request)
}

func (ctx *TraceContext) SpanFinish(span opentracing.Span) {
	if span != nil {
		span.Finish()
	}
}

func (ctx *TraceContext) SpanStart(name string) opentracing.Span {
	if t, err := getTracer(); err == nil {
		if t != nil {
			return t.StartSpan(name, opentracing.FollowsFrom(ctx.TopSpan.Context()))
		}
	}
	return nil
}

// SpanInject 将span注入到request
func (ctx *TraceContext) SpanInject(r *http.Request) {
	if t, err := getTracer(); err == nil {
		if t != nil {
			err = t.Inject(ctx.TopSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil {
				log.Printf("SpanInject, err: %s", err.Error())
			}
		}
	}
}

// TraceContext trace
type TraceContext struct {
	TopSpan opentracing.Span
}

var _ TracerContext = &TraceContext{}

// NewLogContext new obj
func NewLogContext(r *http.Request) *TraceContext {
	ctx := &TraceContext{}

	if t, err := getTracer(); err == nil {
		if t != nil {
			spanCtx, errno := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
			if errno != nil {
				ctx.TopSpan = t.StartSpan(r.RequestURI)
			} else {
				ctx.TopSpan = t.StartSpan(r.RequestURI, ext.RPCServerOption(spanCtx))
			}
		}
	}

	return ctx
}
