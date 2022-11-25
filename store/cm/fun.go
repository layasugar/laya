package cm

import (
	"context"

	"github.com/layasugar/laya"
	"github.com/opentracing/opentracing-go"
)

// ParseSpanByCtx 公共方法, 从ctx中获取
func ParseSpanByCtx(ctx context.Context, spanName string) opentracing.Span {
	return ctx.SpanStart(spanName)
}

// ParseLogIdByCtx 从context中解析出logId
func ParseLogIdByCtx(ctx context.Context) string {
	if webCtx, okInterface := ctx.(*laya.Context); okInterface {
		return webCtx.LogID()
	}
	return ""
}
