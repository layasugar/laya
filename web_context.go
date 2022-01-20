package laya

import (
	"github.com/gin-gonic/gin"
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
		LogContext: NewLogContext(ginContext.Request),
	}
	ginContext.Set(ginFlag, tmp)
	return tmp
}
