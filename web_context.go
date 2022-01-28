package laya

import (
	"github.com/gin-gonic/gin"
	"github.com/layasugar/laya/glogs"
	"github.com/layasugar/laya/gutils"
	uuid "github.com/satori/go.uuid"
)

// WebHandlerFunc http请求的处理者
type WebHandlerFunc func(*WebContext)

// WebContext http 的context
// WebContext 继承了 gin.Context， 并且扩展了日志功能
type WebContext struct {
	*gin.Context
	*LogContext
}

const ginFlag = "__gin__gin"

// NewWebContext 创建 http context
func NewWebContext(ginContext *gin.Context) *WebContext {
	obj, existed := ginContext.Get(ginFlag)
	if existed {
		return obj.(*WebContext)
	}

	tmp := &WebContext{
		Context:    ginContext,
		LogContext: NewLogContext(ginContext.Request, getTraceId(ginContext)),
	}
	ginContext.Set(ginFlag, tmp)
	return tmp
}

// 解析trace_id request_id > x-b3-traceid > uber-trace-id
func getTraceId(ginContext *gin.Context) string {
	// 尝试获取request_id
	traceId := ginContext.GetHeader(glogs.RequestIdKey)
	if traceId != "" {
		return traceId
	}

	// 尝试获取x-b3-traceid
	traceId = ginContext.GetHeader(glogs.ZipkinHeaderKey)
	if traceId != "" {
		ginContext.Request.Header.Set(glogs.RequestIdKey, traceId)
		return traceId
	}

	// 尝试获取uber-trace-id
	traceId = ginContext.GetHeader(glogs.JaegerHeaderKey)
	if traceId != "" {
		ginContext.Request.Header.Set(glogs.RequestIdKey, traceId)
		return traceId
	}

	traceId = gutils.Md5(uuid.NewV4().String())
	ginContext.Request.Header.Set(glogs.RequestIdKey, traceId)
	return traceId
}
