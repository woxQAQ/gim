package db

import (
	"context"
	"fmt"
	"time"

	"github.com/woxQAQ/gim/pkg/logger"

	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// GormLogger 实现gorm.Logger接口的适配器
type GormLogger struct {
	Logger logger.Logger
}

// NewGormLogger 创建一个新的GORM日志适配器
func NewGormLogger(l logger.Logger) *GormLogger {
	return &GormLogger{Logger: l}
}

// LogMode 实现gorm.Logger接口
func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

// Info 实现gorm.Logger接口
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, data...))
}

// Warn 实现gorm.Logger接口
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, data...))
}

// Error 实现gorm.Logger接口
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, data...))
}

// Trace 实现gorm.Logger接口
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 构建日志字段
	fields := []logger.Field{
		logger.String("sql", sql),
		logger.Int64("rows", rows),
		logger.Float64("elapsed_ms", float64(elapsed.Milliseconds())),
		logger.String("source", utils.FileWithLineNum()),
	}

	// 根据执行结果选择日志级别
	if err != nil {
		fields = append(fields, logger.Error(err))
		l.Logger.Error("SQL执行失败", fields...)
		return
	}

	// 根据执行时间选择日志级别
	if elapsed > time.Second {
		l.Logger.Warn("SQL执行时间过长", fields...)
	} else {
		l.Logger.Debug("SQL执行成功", fields...)
	}
}
