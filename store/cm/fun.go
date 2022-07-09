package cm

import (
	"context"
	"github.com/layasugar/laya/service"
	"github.com/opentracing/opentracing-go"
)

// ParseSpanByCtx 公共方法, 从ctx中获取
func ParseSpanByCtx(ctx service.Context, spanName string) opentracing.Span {
	return ctx.SpanStart(spanName)
}

// ParseLogIdByCtx 从context中解析出logId
func ParseLogIdByCtx(ctx context.Context) string {
	if webCtx, okInterface := ctx.(*service.Context); okInterface {
		return webCtx.LogID()
	}
	return ""
}
