// 链路追踪

package trace

import (
	"github.com/layasugar/laya/gcnf/env"
	"github.com/opentracing/opentracing-go"
	"log"
)

const (
	TRACETYPEJAEGER = "jaeger"
	TRACETYPEZIPKIN = "zipkin"
)

// tracer 全局单例变量
var tracer opentracing.Tracer

// InitTrace 初始化trace
func getTracer() opentracing.Tracer {
	if nil == tracer {
		if env.TraceMod() != 0 {
			switch env.TraceType() {
			case TRACETYPEZIPKIN:
				tracer = newZkTracer(env.AppName(), env.LocalIP(), env.TraceAddr(), env.TraceMod())
				log.Printf("[app] tracer success")
			case TRACETYPEJAEGER:
				tracer = newJTracer(env.AppName(), env.TraceAddr(), env.TraceMod())
				log.Printf("[app] tracer success")
			}
		}
	}

	return tracer
}

// ReloadTracer 重载一下tracer
func ReloadTracer() {
	if nil != tracer {
		if env.TraceMod() != 0 {
			switch env.TraceType() {
			case TRACETYPEZIPKIN:
				tracer = newZkTracer(env.AppName(), env.LocalIP(), env.TraceAddr(), env.TraceMod())
				log.Printf("[app] tracer success")
			case TRACETYPEJAEGER:
				tracer = newJTracer(env.AppName(), env.TraceAddr(), env.TraceMod())
				log.Printf("[app] tracer success")
			}
		}
	}
}
