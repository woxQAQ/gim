package db

import (
	"context"
	"fmt"
	"strings"
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

// nolint
const (
	colorSQL     = "\033[36m" // 青色
	colorSQLFunc = "\033[33m" // 黄色
	colorReset   = "\033[0m"
)

// nolint
// highlightSQL 实现SQL语法高亮
func highlightSQL(sql string) string {
	keywords := map[string]bool{
		"SELECT": true, "FROM": true, "WHERE": true, "INSERT": true,
		"INTO": true, "UPDATE": true, "DELETE": true, "CREATE": true,
		"TABLE": true, "INDEX": true, "ON": true, "VALUES": true,
		"SET": true, "AND": true, "OR": true, "NOT": true,
	}

	// 分割保留大小写
	parts := strings.FieldsFunc(sql, func(r rune) bool {
		return r == ' ' || r == '(' || r == ')'
	})

	builder := strings.Builder{}
	for _, part := range parts {
		upperPart := strings.ToUpper(part)
		if keywords[upperPart] {
			builder.WriteString(colorSQLFunc)
			builder.WriteString(part)
			builder.WriteString(colorReset)
		} else {
			builder.WriteString(colorSQL)
			builder.WriteString(part)
			builder.WriteString(colorReset)
		}
		builder.WriteRune(' ')
	}
	return strings.TrimSpace(builder.String())
}
