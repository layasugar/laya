package logx

import (
	"github.com/layasugar/laya/core/alarmx"
	"github.com/layasugar/laya/genv"
)

// LoggerContext 日志
type LoggerContext interface {
	InfoF(template string, args ...interface{})
	WarnF(template string, args ...interface{})
	ErrorF(template string, args ...interface{})
	Field(key string, value interface{}) Field

	Alarm(title string, content string, data map[string]interface{})
}

func (ctx *LogContext) InfoF(template string, args ...interface{}) {
	Info(ctx.logId, template, args...)
}

func (ctx *LogContext) WarnF(template string, args ...interface{}) {
	Warn(ctx.logId, template, args...)
}

// ErrorF 打印程序错误日志
func (ctx *LogContext) ErrorF(template string, args ...interface{}) {
	msg, _ := dealWithArgs(template, args)
	ctx.Alarm("程序打印error", "请引起重视", map[string]interface{}{
		"data": msg,
	})
	Error(ctx.logId, template, args...)
}

func (ctx *LogContext) Field(key string, value interface{}) Field {
	return String(key, value)
}

func (ctx *LogContext) Alarm(title string, content string, data map[string]interface{}) {
	alarmx.Alarm(title, content, data)
}

// LogContext logger
type LogContext struct {
	logId    string
	clientIP string
}

var _ LoggerContext = &LogContext{}

// NewLogContext new obj
func NewLogContext(logId string) *LogContext {
	ctx := &LogContext{
		logId:    logId,
		clientIP: genv.LocalIP(),
	}
	return ctx
}

// GetLogId 得到LogId
func (ctx *LogContext) GetLogId() string {
	return ctx.logId
}

// GetClientIP 得到clientIP
func (ctx *LogContext) GetClientIP() string {
	return ctx.clientIP
}
