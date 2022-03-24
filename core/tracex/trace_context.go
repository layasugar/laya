package tracex

import (
	"github.com/layasugar/laya/core/metautils"
	"github.com/layasugar/laya/env"
	"github.com/layasugar/laya/tools"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkinOt "github.com/openzipkin-contrib/zipkin-go-opentracing"
	uuid "github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	"log"
)

// TracerContext 链路
type TracerContext interface {
	SpanFinish(span opentracing.Span)

	// SpanStart 开启子span
	SpanStart(name string) opentracing.Span

	// SpanInject 注入请求
	SpanInject(md metautils.NiceMD)

	// GetTraceID 获取traceID
	GetTraceID() string
}

func (ctx *TraceContext) SpanFinish(span opentracing.Span) {
	if nil != span {
		span.Finish()
	}
}

func (ctx *TraceContext) SpanStart(name string) opentracing.Span {
	if t := getTracer(); t != nil {
		return t.StartSpan(name, opentracing.FollowsFrom(ctx.TopSpan.Context()))
	}
	return nil
}

// SpanInject 将span注入到request
func (ctx *TraceContext) SpanInject(md metautils.NiceMD) {
	if t := getTracer(); t != nil {
		err := t.Inject(ctx.TopSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(md))
		if err != nil {
			log.Printf("SpanInject, err: %s", err.Error())
		}
	}
}

// GetTraceID 获取traceID
func (ctx *TraceContext) GetTraceID() string {
	return ctx.TraceID
}

// TraceContext trace
type TraceContext struct {
	TopSpan opentracing.Span
	TraceID string
}

var _ TracerContext = &TraceContext{}

// NewTraceContext new traceCtx
func NewTraceContext(name string, headers map[string][]string) *TraceContext {
	ctx := &TraceContext{}

	if env.ApiTrace() {
		if t := getTracer(); t != nil {
			if len(headers) == 0 {
				ctx.TopSpan = t.StartSpan(name)
			} else {
				spanCtx, errno := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
				if errno != nil {
					ctx.TopSpan = t.StartSpan(name)
				} else {
					ctx.TopSpan = t.StartSpan(name, ext.RPCServerOption(spanCtx))
				}
			}
		}
	}

	if ctx.TopSpan != nil {
		spanCtx := ctx.TopSpan.Context()
		switch spanCtx.(type) {
		case jaeger.SpanContext:
			js := spanCtx.(jaeger.SpanContext)
			ctx.TraceID = js.TraceID().String()
		case zipkinOt.SpanContext:
			zs := spanCtx.(zipkinOt.SpanContext)
			ctx.TraceID = zs.TraceID.String()
		}
	}

	if ctx.TraceID == "" {
		ctx.TraceID = tools.Md5(uuid.NewV4().String())
	}
	return ctx
}
