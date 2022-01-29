// 链路追踪

package glogs

import (
	"github.com/layasugar/laya/genv"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"log"
	"net/http"
)

const (
	TraceTypeJaeger = "jaeger"
	TraceTypeZipkin = "zipkin"
)

type (
	Tracer      = opentracing.Tracer
	Span        = opentracing.Span
	SpanContext = opentracing.SpanContext
)

// tracer 全局单例变量
var tracer Tracer

// InitTrace 初始化trace
func getTracer() (Tracer, error) {
	if tracer == nil {
		if genv.TraceMod() != 0 {
			var err error
			switch genv.TraceType() {
			case TraceTypeZipkin:
				tracer = newZkTracer(genv.AppName(), genv.LocalIP(), genv.TraceAddr(), genv.TraceMod())
				if err != nil {
					return nil, err
				}
			case TraceTypeJaeger:
				tracer = newJTracer(genv.AppName(), genv.TraceAddr(), genv.TraceMod())
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return tracer, nil
}

// SpanStart 开启第一个span
func SpanStart(name string) Span {
	t, err := getTracer()
	if err != nil {
		return nil
	}

	if t == nil {
		return nil
	}
	return t.StartSpan(name)
}

// SpanFinish span结束
func SpanFinish(span Span) {
	if span == nil {
		return
	}
	span.Finish()
}

// SpanStartByParent 通过上级span创建span
func SpanStartByParent(ctx SpanContext, name string) Span {
	t, err := getTracer()
	if err != nil {
		return nil
	}

	if t == nil {
		return nil
	}
	return t.StartSpan(name, opentracing.FollowsFrom(ctx))
}

// SpanStartByRequest 通过请求头创建span
func SpanStartByRequest(r *http.Request, name string) Span {
	t, err := getTracer()
	if err != nil {
		return nil
	}

	if t == nil {
		return nil
	}
	spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		return SpanStart(name)
	}
	return tracer.StartSpan(name, ext.RPCServerOption(spanCtx))
}

// SpanInject 将span信息注入请求头
func SpanInject(r *http.Request, span Span) {
	if tracer != nil {
		err := tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			log.Printf("inject header err: %s", err.Error())
		}
	}
}
