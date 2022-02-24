package rdbx

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/layasugar/laya/core/appx"
	"github.com/layasugar/laya/core/grpcx"
	"github.com/layasugar/laya/core/httpx"
	"github.com/layasugar/laya/core/tracex"
	"strings"

	"github.com/layasugar/laya/core/rdbstmt"
)

const (
	componentName = "go-redis"
)

func NewHook() *Hook {
	return &Hook{}
}

type Hook struct{}

func (h *Hook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	var traceCtx *tracex.TraceContext
	switch ctx.(type) {
	case *httpx.WebContext:
		if ctxIf, okInterface := ctx.(*httpx.WebContext); okInterface {
			traceCtx = ctxIf.TraceContext
		}
	case *grpcx.GrpcContext:
		if ctxIf, okInterface := ctx.(*grpcx.GrpcContext); okInterface {
			traceCtx = ctxIf.TraceContext
		}
	case *appx.Context:
		if ctxIf, okInterface := ctx.(*appx.Context); okInterface {
			traceCtx = ctxIf.TraceContext
		}
	}

	if nil != traceCtx {
		span := traceCtx.SpanStart(cmdToSpanName(cmd))
		ext.Component.Set(span, componentName)
		stmt := rdbstmt.NewStatement(cmd.Args())
		ext.DBStatement.Set(span, stmt.ShortString())
		newCtx := context.Background()
		return opentracing.ContextWithSpan(newCtx, span), nil
	}
	return ctx, nil
}

func (h *Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		defer span.Finish()

		isRedisNil := errors.Is(redis.Nil, cmd.Err())
		if cmd.Err() != nil && !isRedisNil {
			ext.LogError(span, cmd.Err())
		} else {
			miss := false
			if isRedisNil {
				miss = true
			}
			span.SetTag("miss", miss)
		}
	}
	return nil
}

func (h *Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	var traceCtx *tracex.TraceContext
	switch ctx.(type) {
	case *httpx.WebContext:
		if ctxIf, okInterface := ctx.(*httpx.WebContext); okInterface {
			traceCtx = ctxIf.TraceContext
		}
	case *grpcx.GrpcContext:
		if ctxIf, okInterface := ctx.(*grpcx.GrpcContext); okInterface {
			traceCtx = ctxIf.TraceContext
		}
	case *appx.Context:
		if ctxIf, okInterface := ctx.(*appx.Context); okInterface {
			traceCtx = ctxIf.TraceContext
		}
	}

	if nil != traceCtx {
		span := traceCtx.SpanStart(spanName("pipeline"))
		ext.Component.Set(span, componentName)
		var stmt rdbstmt.Statement
		for i := 0; i < len(cmds); i++ {
			stmt = rdbstmt.NewStatement(cmds[i].Args())
			span.SetTag(fmt.Sprintf("cmd:%d", i), stmt.ShortString())
		}
		newCtx := context.Background()
		return opentracing.ContextWithSpan(newCtx, span), nil
	}
	return ctx, nil
}

func (h *Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		defer span.Finish()

		for i := 0; i < len(cmds); i++ {
			isRedisNil := errors.Is(redis.Nil, cmds[i].Err())
			if cmds[i].Err() != nil && !isRedisNil {
				ext.LogError(span, cmds[i].Err())
			} else {
				miss := false
				if isRedisNil {
					miss = true
				}
				span.SetTag(fmt.Sprintf("miss:%d", i), miss)
			}
		}
	}
	return nil
}

const spanNamePrefix = "redis:"

func spanName(name string) string {
	return spanNamePrefix + name
}

func cmdToSpanName(cmd redis.Cmder) string {
	return spanName(fullName(cmd))
}

func fullName(cmd redis.Cmder) string {
	return strings.Replace(cmd.FullName(), " ", "_", -1)
}
