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

	// GetTraceID 获取traceID
	GetTraceID() string
}

func (ctx *Context) SpanFinish(span opentracing.Span) {
	if nil != span {
		span.Finish()
	}
}

func (ctx *Context) SpanStart(name string) opentracing.Span {
	if t := getTracer(); t != nil {
		return t.StartSpan(name, opentracing.FollowsFrom(ctx.TopSpan.Context()))
	}
	return nil
}

// SpanInject 将span注入到request
func (ctx *Context) SpanInject(md metautils.NiceMD) {
	if t := getTracer(); t != nil {
		err := t.Inject(ctx.TopSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(md))
		if err != nil {
			log.Printf("SpanInject, err: %s", err.Error())
		}
	}
}

// GetTraceID 获取traceID
func (ctx *Context) GetTraceID() string {
	return ctx.TraceID
}

// Context trace
type Context struct {
	TopSpan opentracing.Span
	TraceID string
}

// NewTraceContext new traceCtx
func NewTraceContext(name string, headers map[string][]string) *Context {
	ctx := &Context{}

	if gcnf.ApiTrace() {
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
		ctx.TraceID = util.Md5(uuid.NewV4().String())
	}
	return ctx
}
