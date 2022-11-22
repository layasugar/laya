package gcnf

func TraceType() string {
	if IsSet(_appTraceType) {
		return GetString(_appTraceType)
	}
	return defaultTraceType
}

func TraceAddr() string {
	if IsSet(_appTraceAddr) {
		return GetString(_appTraceAddr)
	}
	return defaultTraceAddr
}

func TraceMod() float64 {
	if IsSet(_appTraceMod) {
		return GetFloat64(_appTraceMod)
	}
	return defaultTraceMod
}
