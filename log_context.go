package laya

import (
	"github.com/layasugar/laya/glogs"
	"net/http"
)

// LoggerContext 包含链路, 日志, 告警
type LoggerContext interface {
	Info(template string, args ...interface{})
	Warn(template string, args ...interface{})
	Error(template string, args ...interface{})

	// Alarm 告警
	Alarm(msg interface{})

	// StopSpan StartSpan StartSpanParent 开启关闭链路子span
	StopSpan()
	StartSpan()
	StartSpanParent()
}

func (ctx *LogContext) Info(template string, args ...interface{}) {
	glogs.Info(ctx.req, template, args...)
}

func (ctx *LogContext) Warn(template string, args ...interface{}) {
	glogs.Warn(ctx.req, template, args...)
}

func (ctx *LogContext) Error(template string, args ...interface{}) {
	glogs.Error(ctx.req, template, args...)
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
