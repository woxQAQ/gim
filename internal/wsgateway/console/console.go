package console

import (
	"github.com/woxQAQ/gim/internal/wsgateway"
	"github.com/woxQAQ/gim/internal/wsgateway/console/cmd"
	"github.com/woxQAQ/gim/pkg/logger"
)

// Console 提供WebSocket网关的控制台界面
type Console struct {
	gateway wsgateway.Gateway
	logger  logger.Logger
}

// NewConsole 创建新的控制台实例
func NewConsole(gateway wsgateway.Gateway, logger logger.Logger) *Console {
	return &Console{
		gateway: gateway,
		logger:  logger,
	}
}

// Start 启动控制台
func (c *Console) Start() {
	cmd.Execute(c.gateway, c.logger)
}

// Stop 停止控制台
func (c *Console) Stop() {
	// 在使用cobra框架后，不需要特殊的停止逻辑
}
