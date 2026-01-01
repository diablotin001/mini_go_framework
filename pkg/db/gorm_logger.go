package db

import (
    "context"
    "time"
    "mini_go/logger"
    "go.uber.org/zap"
    gormlog "gorm.io/gorm/logger"
)

type ZapGormLogger struct {
    SlowThreshold time.Duration
    LogLevel      gormlog.LogLevel
}

func NewZapGormLogger(slowMS int, level gormlog.LogLevel) *ZapGormLogger {
    return &ZapGormLogger{SlowThreshold: time.Duration(slowMS) * time.Millisecond, LogLevel: level}
}

func (l *ZapGormLogger) LogMode(level gormlog.LogLevel) gormlog.Interface { l.LogLevel = level; return l }
func (l *ZapGormLogger) Info(ctx context.Context, msg string, args ...interface{})  { logger.Log.Sugar().Infof(msg, args...) }
func (l *ZapGormLogger) Warn(ctx context.Context, msg string, args ...interface{})  { logger.Log.Sugar().Warnf(msg, args...) }
func (l *ZapGormLogger) Error(ctx context.Context, msg string, args ...interface{}) { logger.Log.Sugar().Errorf(msg, args...) }

func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
    elapsed := time.Since(begin)
    sql, rows := fc()
    fields := []zap.Field{zap.Duration("elapsed", elapsed), zap.String("sql", sql), zap.Int64("rows", rows)}
    switch {
    case err != nil:
        logger.Log.Error("SQL Error", append(fields, zap.Error(err))...)
    case l.SlowThreshold > 0 && elapsed > l.SlowThreshold:
        logger.Log.Warn("Slow SQL", fields...)
    default:
        logger.Log.Info("SQL", fields...)
    }
}
