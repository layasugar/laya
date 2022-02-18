package dbx

import (
	"context"
	"database/sql"
	"errors"
	"github.com/layasugar/laya/core/appx"
	"github.com/layasugar/laya/core/grpcx"
	"github.com/layasugar/laya/core/httpx"
	"github.com/layasugar/laya/core/tracex"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gorm.io/gorm"
)

const (
	contextKey    = "otgorm:context"
	tSpanName     = "mysql"
	componentName = "gorm"
)

func Wrap(ctx context.Context, dbName ...string) *gorm.DB {
	var db *gorm.DB
	if len(dbName) > 0 {
		db = getGormDB(dbName[0])
	} else {
		db = getGormDB(defaultDbName)
	}

	return db.Set(contextKey, ctx)
}

func registerCallbacks(db *gorm.DB) {
	prefix := db.Dialector.Name() + ":"

	db.Callback().Create().Before("gorm:begin_transaction").Register("aotgorm_before_create", newBefore(prefix+"create"))
	db.Callback().Create().After("gorm:commit_or_rollback_transaction").Register("otgorm_after_create", newAfter())

	db.Callback().Update().Before("gorm:begin_transaction").Register("otgorm_before_update", newBefore(prefix+"update"))
	db.Callback().Update().After("gorm:commit_or_rollback_transaction").Register("otgorm_after_update", newAfter())

	db.Callback().Query().Before("gorm:query").Register("otgorm_before_query", newBefore(prefix+"query"))
	db.Callback().Query().After("gorm:after_query").Register("otgorm_after_query", newAfter())

	db.Callback().Delete().Before("gorm:begin_transaction").Register("otgorm_before_delete", newBefore(prefix+"delete"))
	db.Callback().Delete().After("gorm:commit_or_rollback_transaction").Register("otgorm_after_delete", newAfter())

	db.Callback().Row().Before("gorm:row").Register("otgorm_before_row", newBefore(prefix+"row"))
	db.Callback().Row().After("gorm:row").Register("otgorm_after_row", newAfter())

	db.Callback().Raw().Before("gorm:raw").Register("otgorm_before_raw", newBefore(prefix+"raw"))
	db.Callback().Raw().After("gorm:raw").Register("otgorm_after_raw", newAfter())
}

func newBefore(name string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if v, ok := db.Get(contextKey); ok {
			var traceCtx *tracex.TraceContext
			switch v.(type) {
			case *httpx.WebContext:
				if ctx, okInterface := v.(*httpx.WebContext); okInterface {
					traceCtx = ctx.TraceContext
				}
			case *grpcx.GrpcContext:
				if ctx, okInterface := v.(*grpcx.GrpcContext); okInterface {
					traceCtx = ctx.TraceContext
				}
			case *appx.Context:
				if ctx, okInterface := v.(*appx.Context); okInterface {
					traceCtx = ctx.TraceContext
				}
			}
			if traceCtx != nil {
				span := traceCtx.SpanStart(tSpanName)
				if nil != span {
					newCtx := context.Background()
					keepScene(db, newCtx)
					ext.Component.Set(span, componentName)
					setSpan(db, newCtx, span)
				}
			}
		}
	}
}

func newAfter() func(*gorm.DB) {
	return func(db *gorm.DB) {
		span, _ := getSpan(db)
		if nil != span {
			defer func() {
				span.Finish()
				restoreScene(db)
			}()
			ext.DBStatement.Set(span, db.Statement.SQL.String())
			if db.Error != nil {
				if !errors.Is(db.Error, gorm.ErrRecordNotFound) && !errors.Is(db.Error, sql.ErrNoRows) {
					ext.LogError(span, db.Error)
				}
			}
		}
	}
}

func setSpan(db *gorm.DB, ctx context.Context, span opentracing.Span) {
	db.Set(contextKey, opentracing.ContextWithSpan(ctx, span))
}

func getSpan(db *gorm.DB) (opentracing.Span, context.Context) {
	if v, ok := db.Get(contextKey); ok {
		if ctx, ok := v.(context.Context); ok {
			return opentracing.SpanFromContext(ctx), ctx
		}
	}
	return nil, nil
}

const contextSceneKey = "otgorm:context:scene:" + "v1.0.0"

func keepScene(db *gorm.DB, ctx context.Context) {
	db.Set(contextSceneKey, ctx)
}

func restoreScene(db *gorm.DB) {
	if v, ok := db.Get(contextSceneKey); ok {
		db.Set(contextKey, v)
	}
}
