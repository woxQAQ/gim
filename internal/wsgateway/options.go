package wsgateway

import (
	"time"

	"github.com/woxQAQ/gim/internal/wsgateway/codec"
	"github.com/woxQAQ/gim/pkg/logger"
)

// Option 定义WSGateway的配置选项函数类型.
type Option func(*WSGateway)

// WithLogger 设置WSGateway的logger.
func WithLogger(l logger.Logger) Option {
	return func(g *WSGateway) {
		g.logger = l
	}
}

// WithCompressor 设置WSGateway的压缩器.
func WithCompressor(c codec.Compressor) Option {
	return func(g *WSGateway) {
		g.compressor = c
	}
}

// WithEncoder 设置WSGateway的编码器.
func WithEncoder(e codec.Encoder) Option {
	return func(g *WSGateway) {
		g.encoder = e
	}
}

func WithHeartbeat(interval, timeout time.Duration) Option {
	return func(g *WSGateway) {
		g.heartbeatInterval = interval
		g.heartbeatTimeout = timeout
	}
}
