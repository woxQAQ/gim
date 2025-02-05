package handler

import (
	"github.com/woxQAQ/gim/internal/wsgateway/types"
	"github.com/woxQAQ/gim/internal/wsgateway/user"
)

// ForwardHandler 消息转发处理器
type ForwardHandler struct {
	BaseHandler
	userManager user.IUserManager
}

// NewForwardHandler 创建消息转发处理器
func NewForwardHandler(userManager user.IUserManager) *ForwardHandler {
	return &ForwardHandler{
		userManager: userManager,
	}
}

// Handle 实现消息转发逻辑
func (h *ForwardHandler) Handle(msg types.Message) (bool, error) {
	// 根据消息类型和目标进行转发
	switch msg.Header.Type {
	case types.MessageTypeText, types.MessageTypeImage,
		types.MessageTypeVideo, types.MessageTypeAudio,
		types.MessageTypeFile:
		// 如果消息有特定目标用户，则转发给目标用户
		if msg.Header.To != "" {
			err := h.userManager.SendPlatformMessage(msg.Header.To, msg.Header.Platform, msg)
			if err != nil {
				return false, err
			}
		}
	}

	// 继续处理链
	return true, nil
}

// StoreHandler 消息存储处理器
type StoreHandler struct {
	BaseHandler
}

// NewStoreHandler 创建消息存储处理器
func NewStoreHandler() *StoreHandler {
	return &StoreHandler{}
}

// Handle 实现消息存储逻辑
func (h *StoreHandler) Handle(msg types.Message) (bool, error) {
	// TODO: 实现消息存储逻辑
	// 这里可以添加将消息保存到数据库或其他存储系统的逻辑

	// 继续处理链
	return true, nil
}

// NewMessageChain 创建默认的消息处理链
func NewMessageChain(userManager user.IUserManager) *Chain {
	chain := NewChain()

	// 添加消息转发处理器
	chain.AddHandler(NewForwardHandler(userManager))

	// 添加消息存储处理器
	chain.AddHandler(NewStoreHandler())

	return chain
}
