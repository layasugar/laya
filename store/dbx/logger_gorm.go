package dbx

import (
	"context"
	"fmt"
	"github.com/layasugar/laya/env"
	"github.com/layasugar/laya/store/cm"
	"github.com/layasugar/laya/tools"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

const (
	sqlTitle = "mysql"
)

func Default(level logger.LogLevel) logger.Interface {
	var config = logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      level,
		Colorful:      true,
	}
	var (
		infoStr      = "{\"line\": \"%s\", \"level\": \"[info]\", \"msg\": \"%s\"}"
		warnStr      = "{\"line\": \"%s\", \"level\": \"[warn]\", \"msg\": \"%s\"}"
		errStr       = "{\"line\": \"%s\", \"level\": \"[error]\", \"msg\": \"%s\"}"
		traceStr     = "{\"line\": \"%s\", \"耗时\": \"%.3fms\", \"rows\": \"%v\" ,\"sql\": \"%s\"}"
		traceWarnStr = "{\"line\": \"%s\", \"错误\": \"%s\", \"耗时\": \"%.3fms\", \"rows\": \"%v\", \"sql\": \"%s\"}"
		traceErrStr  = "{\"line\": \"%s\", \"slow\": \"%s\", \"耗时\": \"%.3fms\", \"rows\": \"%v\", \"sql\": \"%s\"}"
	)

	return &gormLogger{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type gormLogger struct {
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode logger mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		errInfo := fmt.Sprintf(msg, data)
		gormWriter(ctx, logger.LevelInfo, fmt.Sprintf(l.infoStr, utils.FileWithLineNum(), errInfo))
	}
}

// Warn print warn messages
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		errInfo := fmt.Sprintf(msg, data)
		gormWriter(ctx, logger.LevelWarn, fmt.Sprintf(l.infoStr, utils.FileWithLineNum(), errInfo))
	}
}

// Error print error messages
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		errInfo := fmt.Sprintf(msg, data)
		gormWriter(ctx, logger.LevelError, fmt.Sprintf(l.infoStr, utils.FileWithLineNum(), errInfo))
	}
}

// Trace print sql message
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			sql, rows := fc()
			if rows == -1 {
				msg := fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err.Error(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
				gormWriter(ctx, logger.LevelError, msg)
			} else {
				msg := fmt.Sprintf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
				gormWriter(ctx, logger.LevelError, msg)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				msg := fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
				gormWriter(ctx, logger.LevelWarn, msg)
			} else {
				msg := fmt.Sprintf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
				gormWriter(ctx, logger.LevelWarn, msg)
			}
		case l.LogLevel >= logger.Info:
			sql, rows := fc()
			if rows == -1 {
				msg := fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
				gormWriter(ctx, logger.LevelInfo, msg)
			} else {
				msg := fmt.Sprintf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
				gormWriter(ctx, logger.LevelInfo, msg)
			}
		}
	}
}

func gormWriter(ctx context.Context, level, msg string, fields ...zap.Field) {
	requestId := cm.ParseLogIdByCtx(ctx)
	fields = append(fields, zap.String(tools.RequestIdKey, requestId), zap.String("title", sqlTitle))

	switch level {
	case logger.LevelInfo:
		if env.MysqlLog() {
			logger.GetSugar().Info(msg, fields...)
		}
	case logger.LevelWarn:
		if env.MysqlLog() {
			logger.GetSugar().Warn(msg, fields...)
		}
	case logger.LevelError:
		logger.GetSugar().Error(msg, fields...)
	}

	return
}
