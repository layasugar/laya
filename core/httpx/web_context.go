package httpx

import (
	"github.com/gin-gonic/gin"
	"github.com/layasugar/laya/core/alarmer"
	"github.com/layasugar/laya/core/logger"
	"github.com/layasugar/laya/core/tracer"
)

// WebHandlerFunc http请求的处理者
type WebHandlerFunc func(*WebContext)

// WebContext http 的context
// WebContext 继承了 gin.Context, 并且扩展了日志功能
type WebContext struct {
	*gin.Context
	*logger.Context
	*tracer.TraceContext
	*alarmer.AlarmContext
}

const ginFlag = "__gin__gin"

// NewWebContext 创建 http contextx
func NewWebContext(ginContext *gin.Context) *WebContext {
	obj, existed := ginContext.Get(ginFlag)
	if existed {
		return obj.(*WebContext)
	}
	traceCtx := tracer.NewTraceContext(ginContext.Request.RequestURI, ginContext.Request.Header)

	tmp := &WebContext{
		Context:      ginContext,
		Context:      logger.NewContext(traceCtx.TraceID),
		TraceContext: traceCtx,
	}
	ginContext.Set(ginFlag, tmp)

	return tmp
}
