package gcnf

func TraceType() string {
	if gcf.IsSet(_appTraceType) {
		return gcf.GetString(_appTraceType)
	}
	return defaultTraceType
}

func TraceAddr() string {
	if gcf.IsSet(_appTraceAddr) {
		return gcf.GetString(_appTraceAddr)
	}
	return defaultTraceAddr
}

func TraceMod() float64 {
	if gcf.IsSet(_appTraceMod) {
		return gcf.GetFloat64(_appTraceMod)
	}
	return defaultTraceMod
}
