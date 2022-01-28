package laya

import (
	"github.com/layasugar/laya/genv"
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

	// StartSpan StopSpan StartSpanP StartSpanR 开启,关闭,通过上级span开启span, 通过request开启span
	StartSpan() glogs.Span
	StopSpan(span glogs.Span)
	StartSpanP(span glogs.Span, name string) glogs.Span
	StartSpanR(name string) glogs.Span
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

// StartSpan StopSpan StartSpanP StartSpanR 开启,关闭,通过上级span开启span, 通过request开启span
func (ctx *LogContext) StartSpan() glogs.Span    { return glogs.StartSpan(genv.AppName()) }
func (ctx *LogContext) StopSpan(span glogs.Span) { glogs.StopSpan(span) }
func (ctx *LogContext) StartSpanP(span glogs.Span, name string) glogs.Span {
	return glogs.StartSpanP(span.Context(), name)
}
func (ctx *LogContext) StartSpanR(name string) glogs.Span { return glogs.StartSpanR(ctx.req, name) }

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
func NewLogContext(req *http.Request, traceId string) *LogContext {
	ctx := &LogContext{
		req:      req,
		traceId:  traceId,
		clientIP: genv.LocalIP(),
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
