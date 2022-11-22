package trace

import (
	"github.com/layasugar/laya/core/metautils"
	"github.com/layasugar/laya/core/util"
	"github.com/layasugar/laya/gcnf"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkinOt "github.com/openzipkin-contrib/zipkin-go-opentracing"
	uuid "github.com/satori/go.uuid"
	"github.com/uber/jaeger-client-go"
	"log"
)

// Trace 链路
type Trace interface {
	SpanFinish(span opentracing.Span)

	// SpanStart 开启子span
	SpanStart(name string) opentracing.Span

	// SpanInject 注入请求
	SpanInject(md metautils.NiceMD)

	// TraceID 获取traceID
	TraceID() string
}

// Context trace
type Context struct {
	topSpan opentracing.Span
	traceID string
}

// NewTraceContext new traceCtx
func NewTraceContext(name string, headers map[string][]string) *Context {
	ctx := &Context{}

	if gcnf.ApiTrace() {
		if t := getTracer(); t != nil {
			if len(headers) == 0 {
				ctx.topSpan = t.StartSpan(name)
			} else {
				spanCtx, errno := t.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
				if errno != nil {
					ctx.topSpan = t.StartSpan(name)
				} else {
					ctx.topSpan = t.StartSpan(name, ext.RPCServerOption(spanCtx))
				}
			}
		}
	}

	if ctx.topSpan != nil {
		spanCtx := ctx.topSpan.Context()
		switch spanCtx.(type) {
		case jaeger.SpanContext:
			js := spanCtx.(jaeger.SpanContext)
			ctx.traceID = js.TraceID().String()
		case zipkinOt.SpanContext:
			zs := spanCtx.(zipkinOt.SpanContext)
			ctx.traceID = zs.TraceID.String()
		}
	}

	if ctx.traceID == "" {
		ctx.traceID = util.Md5(uuid.NewV4().String())
	}
	return ctx
}

func (ctx *Context) SpanFinish(span opentracing.Span) {
	if nil != span {
		span.Finish()
	}
}

func (ctx *Context) SpanStart(name string) opentracing.Span {
	if t := getTracer(); t != nil {
		return t.StartSpan(name, opentracing.FollowsFrom(ctx.topSpan.Context()))
	}
	return nil
}

// SpanInject 将span注入到request
func (ctx *Context) SpanInject(md metautils.NiceMD) {
	if t := getTracer(); t != nil {
		err := t.Inject(ctx.topSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(md))
		if err != nil {
			log.Printf("SpanInject, err: %s", err.Error())
		}
	}
}

// TraceID 获取traceID
func (ctx *Context) TraceID() string {
	return ctx.traceID
}
