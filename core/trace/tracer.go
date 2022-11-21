// 链路追踪

package trace

import (
	"github.com/layasugar/laya/gcnf"
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
		if gcnf.TraceMod() != 0 {
			switch gcnf.TraceType() {
			case TRACETYPEZIPKIN:
				tracer = newZkTracer(gcnf.AppName(), gcnf.LocalIP(), gcnf.TraceAddr(), gcnf.TraceMod())
				log.Printf("[app] tracer success")
			case TRACETYPEJAEGER:
				tracer = newJTracer(gcnf.AppName(), gcnf.TraceAddr(), gcnf.TraceMod())
				log.Printf("[app] tracer success")
			}
		}
	}

	return tracer
}

// ReloadTracer 重载一下tracer
func ReloadTracer() {
	if nil != tracer {
		if gcnf.TraceMod() != 0 {
			switch gcnf.TraceType() {
			case TRACETYPEZIPKIN:
				tracer = newZkTracer(gcnf.AppName(), gcnf.LocalIP(), gcnf.TraceAddr(), gcnf.TraceMod())
				log.Printf("[app] tracer success")
			case TRACETYPEJAEGER:
				tracer = newJTracer(gcnf.AppName(), gcnf.TraceAddr(), gcnf.TraceMod())
				log.Printf("[app] tracer success")
			}
		}
	}
}
