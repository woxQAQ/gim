package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Domain 定义日志所属的领域。
type Domain string

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
)

const (
	DomainWSGateway Domain = "wsgateway" // WebSocket网关领域
	DomainAPIServer Domain = "apiserver" // API服务器领域
	DomainDatabase  Domain = "database"  // 数据库领域
)

// Config 日志配置.
type Config struct {
	Level    string // 日志级别：debug, info, warn, error
	FilePath string // 日志文件路径，为空时仅输出到控制台
}

// defaultConfig 默认日志配置.
var defaultConfig = Config{
	Level:    "info",
	FilePath: "",
}

// Logger 统一日志接口.
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger

	// adapt with amqp logger
	Printf(format string, v ...interface{})
	Enable()
	Disable()
}

// Field 日志字段.
type Field = zap.Field

// 提供创建Field的便捷方法.
var (
	String  = zap.String
	Int     = zap.Int
	Int32   = zap.Int32
	Int64   = zap.Int64
	Time    = zap.Time
	Float64 = zap.Float64
	Bool    = zap.Bool
	Any     = zap.Any
	Error   = zap.Error
)

// logger 实现Logger接口.
type logger struct {
	zapLogger *zap.Logger
	domain    Domain
	disabled  bool
}

// Enable 启用日志输出
func (l *logger) Enable() {
	l.disabled = false
}

// Disable 禁用日志输出
func (l *logger) Disable() {
	l.disabled = true
}

// Debug 实现Logger接口
func (l *logger) Debug(msg string, fields ...Field) {
	if !l.disabled {
		l.zapLogger.Debug(msg, fields...)
	}
}

// Info 实现Logger接口
func (l *logger) Info(msg string, fields ...Field) {
	if !l.disabled {
		l.zapLogger.Info(msg, fields...)
	}
}

// Warn 实现Logger接口
func (l *logger) Warn(msg string, fields ...Field) {
	if !l.disabled {
		l.zapLogger.Warn(msg, fields...)
	}
}

// Error 实现Logger接口
func (l *logger) Error(msg string, fields ...Field) {
	if !l.disabled {
		l.zapLogger.Error(msg, fields...)
	}
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.zapLogger.Sugar().Infof(format, v...)
}

// With 实现Logger接口
func (l *logger) With(fields ...Field) Logger {
	return &logger{
		zapLogger: l.zapLogger.With(fields...),
		domain:    l.domain,
		disabled:  l.disabled,
	}
}

// NewLogger 创建指定领域的日志记录器.
func NewLogger(cfg *Config, opts ...zap.Option) (Logger, error) {
	// 使用默认配置
	if cfg == nil {
		cfg = &defaultConfig
	}

	// 配置日志输出.
	var cores []zapcore.Core

	// 添加控制台输出.
	// 自定义日志编码配置
	encConfig := zap.NewDevelopmentEncoderConfig()
	encConfig.LevelKey = "level"
	encConfig.TimeKey = "timestamp"
	encConfig.CallerKey = "caller"
	encConfig.MessageKey = "message"
	encConfig.EncodeCaller = customCallerEncoder
	encConfig.EncodeLevel = customLevelEncoder
	encConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encConfig),
		zapcore.AddSync(os.Stdout),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.DebugLevel
		}),
	)
	cores = append(cores, consoleCore)

	// 如果指定了文件路径，添加文件输出.
	if cfg.FilePath != "" {
		// 创建日志目录
		fileEncConfig := zap.NewProductionEncoderConfig()
		fileEncConfig.TimeKey = "timestamp"
		fileEncConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		if err := os.MkdirAll(filepath.Dir(cfg.FilePath), 0o755); err != nil {
			return nil, fmt.Errorf("create log directory failed: %w", err)
		}

		// 配置文件输出.
		// 使用相同的编码配置
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(fileEncConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename: cfg.FilePath,
				MaxSize:  100, // 默认100MB
			}),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.InfoLevel
			}),
		)
		cores = append(cores, fileCore)
	}

	// 创建多输出核心.
	core := zapcore.NewTee(cores...)

	// 创建logger.
	zapLogger := zap.New(core, zap.AddCaller())
	for _, opt := range opts {
		zapLogger = zapLogger.WithOptions(opt)
	}
	l := &logger{
		zapLogger: zapLogger,
	}

	return l, nil
}

func customLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var color string
	switch l {
	case zapcore.DebugLevel:
		color = colorBlue
	case zapcore.InfoLevel:
		color = colorGreen
	case zapcore.WarnLevel:
		color = colorYellow
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		color = colorRed
	default:
		color = colorReset
	}
	enc.AppendString(fmt.Sprintf("%s%-5s%s", color, l.CapitalString(), colorReset))
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%s%s%s", colorCyan, caller.TrimmedPath(), colorReset))
}
