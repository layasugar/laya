package gutils

const (
	XForwardedFor = "X-Forwarded-For" // 获取真实ip
	XRealIP       = "X-Real-IP"       // 获取真实ip
)

var IgnoreRoutes = []string{"/", "/ready", "/health", "/reload", "/metrics"}
