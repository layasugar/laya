package env

import "github.com/layasugar/laya/gcnf"

const (
	defaultTraceType         = ""
	defaultTraceAddr         = ""
	defaultTraceMod  float64 = 0

	_appTraceType = "app.trace.type"
	_appTraceAddr = "app.trace.addr"
	_appTraceMod  = "app.trace.mod"
)

func TraceType() string {
	if gcnf.IsSet(_appTraceType) {
		return gcnf.GetString(_appTraceType)
	}
	return defaultTraceType
}

func TraceAddr() string {
	if gcnf.IsSet(_appTraceAddr) {
		return gcnf.GetString(_appTraceAddr)
	}
	return defaultTraceAddr
}

func TraceMod() float64 {
	if gcnf.IsSet(_appTraceMod) {
		return gcnf.GetFloat64(_appTraceMod)
	}
	return defaultTraceMod
}
