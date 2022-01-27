package laya

import (
	"github.com/layasugar/laya/glogs"
	"net/http"
)

// LoggerContext 包含链路, 日志, 告警
type LoggerContext interface {
	InfoF(template string, args ...interface{})
	WarnF(template string, args ...interface{})
	ErrorF(template string, args ...interface{})
	Field(key string, value interface{}) glogs.Field

	// Alarm 告警
	Alarm(msg interface{})

	// StopSpan StartSpan StartSpanParent 开启关闭链路子span
	StopSpan()
	StartSpan()
	StartSpanParent()
}

func (ctx *LogContext) InfoF(template string, args ...interface{}) {
	glogs.Info(ctx.req, template, args...)
}

func (ctx *LogContext) WarnF(template string, args ...interface{}) {
	glogs.Warn(ctx.req, template, args...)
}

func (ctx *LogContext) ErrorF(template string, args ...interface{}) {
	glogs.Error(ctx.req, template, args...)
}

func (ctx *LogContext) Field(key string, value interface{}) glogs.Field {
	return glogs.String(key, value)
}

// Alarm 通知
func (ctx *LogContext) Alarm(msg interface{}) {}

// StopSpan StartSpan StartSpanParent 开启关闭链路子span
func (ctx *LogContext) StopSpan()        {}
func (ctx *LogContext) StartSpan()       {}
func (ctx *LogContext) StartSpanParent() {}

// LogContext logger
type LogContext struct {
	req               *http.Request
	traceId           string
	clientIP          string
	RspHTTPStatusCode int
	ErrMsg            string
}

var _ LoggerContext = &LogContext{}

// NewLogContext new obj
func NewLogContext(req *http.Request) *LogContext {
	ctx := &LogContext{
		req: req,
	}
	return ctx
}

// GetTraceId 得到TraceID
func (ctx *LogContext) GetTraceId() string {
	return ctx.traceId
}

// GetClientIP 得到clientIP
func (ctx *LogContext) GetClientIP() string {
	return ctx.clientIP
}
