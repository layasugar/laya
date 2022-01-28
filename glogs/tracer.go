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

func StopSpan(span Span) {
	if span == nil {
		return
	}
	span.Finish()
}
func StartSpan(name string) Span {
	t, err := getTracer()
	if err != nil {
		return nil
	}

	if t == nil {
		return nil
	}
	return t.StartSpan(name)
}
func StartSpanP(ctx SpanContext, name string) Span {
	t, err := getTracer()
	if err != nil {
		return nil
	}

	if t == nil {
		return nil
	}
	return t.StartSpan(name, opentracing.FollowsFrom(ctx))
}
func StartSpanR(r *http.Request, name string) Span {
	t, err := getTracer()
	if err != nil {
		return nil
	}

	if t == nil {
		return nil
	}
	spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
	if err != nil {
		log.Println(err.Error())
	}
	return tracer.StartSpan(name, ext.RPCServerOption(spanCtx))
}
