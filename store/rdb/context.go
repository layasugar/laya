package rdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/layasugar/laya/env"
	"github.com/layasugar/laya/store/cm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
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
	if env.RedisTrace() {
		span := cm.ParseSpanByCtx(ctx, cmdToSpanName(cmd))

		if nil != span {
			ext.Component.Set(span, componentName)
			stmt := rdbstmt.NewStatement(cmd.Args())
			ext.DBStatement.Set(span, stmt.ShortString())
			newCtx := context.Background()
			return opentracing.ContextWithSpan(newCtx, span), nil
		}
	}
	return ctx, nil
}

func (h *Hook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if env.RedisTrace() {
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
	}
	return nil
}

func (h *Hook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if env.RedisTrace() {
		span := cm.ParseSpanByCtx(ctx, spanName("pipeline"))
		if nil != span {
			ext.Component.Set(span, componentName)
			var stmt rdbstmt.Statement
			for i := 0; i < len(cmds); i++ {
				stmt = rdbstmt.NewStatement(cmds[i].Args())
				span.SetTag(fmt.Sprintf("cmd:%d", i), stmt.ShortString())
			}
			newCtx := context.Background()
			return opentracing.ContextWithSpan(newCtx, span), nil
		}
	}
	return ctx, nil
}

func (h *Hook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if env.RedisTrace() {
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
