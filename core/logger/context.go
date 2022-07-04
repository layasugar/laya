package logger

// Logger 日志
type Logger interface {
	Info(template string, args ...interface{})
	Warn(template string, args ...interface{})
	Error(template string, args ...interface{})
	Field(key string, value interface{}) Field
}

func (ctx *Context) Info(template string, args ...interface{}) {
	Info(ctx.logID, template, args...)
}

func (ctx *Context) Warn(template string, args ...interface{}) {
	Warn(ctx.logID, template, args...)
}

// ErrorF 打印程序错误日志
func (ctx *Context) Error(template string, args ...interface{}) {
	//msg, _ := dealWithArgs(template, args)
	Error(ctx.logID, template, args...)
}

func (ctx *Context) Field(key string, value interface{}) Field {
	return String(key, value)
}

// Context logger
type Context struct {
	logID    string
	clientIP string
}

var _ Logger = &Context{}

// NewContext new obj
func NewContext(logID string, ip string) *Context {
	ctx := &Context{
		logID:    logID,
		clientIP: ip,
	}
	return ctx
}

// LogID 得到LogId
func (ctx *Context) LogID() string {
	return ctx.logID
}

// ClientIP 得到clientIP
func (ctx *Context) ClientIP() string {
	return ctx.clientIP
}
