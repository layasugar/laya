package laya

import (
	"github.com/gin-gonic/gin"
	"github.com/layasugar/laya/core/alarmx"
	"github.com/layasugar/laya/core/logx"
	"github.com/layasugar/laya/core/tracex"
	"github.com/layasugar/laya/gtools"
	uuid "github.com/satori/go.uuid"
)

// WebHandlerFunc http请求的处理者
type WebHandlerFunc func(*WebContext)

// WebContext http 的context
// WebContext 继承了 gin.Context, 并且扩展了日志功能
type WebContext struct {
	*gin.Context
	*logx.LogContext
	*tracex.TraceContext
	*alarmx.AlarmContext
}

const ginFlag = "__gin__gin"

// NewWebContext 创建 http context
func NewWebContext(ginContext *gin.Context) *WebContext {
	obj, existed := ginContext.Get(ginFlag)
	if existed {
		return obj.(*WebContext)
	}

	logId := ginContext.GetHeader(gtools.RequestIdKey)
	if logId == "" {
		logId = gtools.Md5(uuid.NewV4().String())
		ginContext.Request.Header.Set(gtools.RequestIdKey, logId)
		ginContext.Set(gtools.RequestIdKey, logId)
	}

	tmp := &WebContext{
		Context:      ginContext,
		LogContext:   logx.NewLogContext(logId),
		TraceContext: tracex.NewTraceContext(ginContext.Request.RequestURI, ginContext.Request.Header),
	}
	ginContext.Set(ginFlag, tmp)

	return tmp
}
