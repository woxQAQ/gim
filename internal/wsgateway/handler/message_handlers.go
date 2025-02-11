package handler

import (
	"errors"

	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"github.com/woxQAQ/gim/internal/models"
	"github.com/woxQAQ/gim/internal/types"
	"github.com/woxQAQ/gim/internal/wsgateway/base"
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
func (h *ForwardHandler) Handle(msg base.IMessage) (bool, error) {
	// 根据消息类型和目标进行转发
	switch msg.GetType() {
	case types.MessageTypeText, types.MessageTypeImage,
		types.MessageTypeVideo, types.MessageTypeAudio,
		types.MessageTypeFile:
		// 如果消息有特定目标用户，则转发给目标用户
		if msg.GetTo() != "" {
			// 不再使用Platform字段，确保消息能够正确转发给目标用户
			err := h.userManager.SendMessage(msg.GetTo(), msg)
			if err != nil {
				return false, errors.Join(err...)
			}
		}
	}

	// 继续处理链
	return true, nil
}

// StoreHandler 消息存储处理器
type StoreHandler struct {
	BaseHandler
	messageStore *stores.MessageStore
}

// NewStoreHandler 创建消息存储处理器
func NewStoreHandler(messageStore *stores.MessageStore) *StoreHandler {
	return &StoreHandler{messageStore: messageStore}
}

// Handle 实现消息存储逻辑
func (h *StoreHandler) Handle(msg base.IMessage) (bool, error) {
	// 检查消息存储器是否已初始化
	if h.messageStore == nil {
		return false, errors.New("message store is not initialized")
	}

	// 将消息转换为数据库模型
	message := &models.Message{}
	message.FromTypes(msg)

	// 保存消息到数据库
	err := h.messageStore.CreateMessage(message)
	if err != nil {
		return false, err
	}

	// 继续处理链
	return true, nil
}

// NewMessageChain 创建默认的消息处理链
func NewMessageChain(userManager user.IUserManager, ms *stores.MessageStore) *Chain {
	chain := NewChain()

	// 添加消息转发处理器
	chain.AddHandler(NewForwardHandler(userManager))

	// 添加消息存储处理器
	chain.AddHandler(NewStoreHandler(ms))

	return chain
}
