package env

import "github.com/layasugar/laya/gcf"

const (
	defaultTraceType         = ""
	defaultTraceAddr         = ""
	defaultTraceMod  float64 = 0
)

func TraceType() string {
	if gcf.IsSet("app.trace.type") {
		return gcf.GetString("app.trace.type")
	}
	return defaultTraceType
}

func TraceAddr() string {
	if gcf.IsSet("app.trace.addr") {
		return gcf.GetString("app.trace.addr")
	}
	return defaultTraceAddr
}

func TraceMod() float64 {
	if gcf.IsSet("app.trace.mod") {
		return gcf.GetFloat64("app.trace.mod")
	}
	return defaultTraceMod
}
