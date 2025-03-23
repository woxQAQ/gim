package handler

import (
	"errors"

	"github.com/panjf2000/gnet/v2/pkg/buffer/linkedlist"

	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"github.com/woxQAQ/gim/internal/models"
	"github.com/woxQAQ/gim/internal/types"
	"github.com/woxQAQ/gim/internal/wsgateway/codec"
	"github.com/woxQAQ/gim/internal/wsgateway/user"
)

// ForwardHandler 消息转发处理器
type ForwardHandler struct {
	BaseHandler
	userManager user.IUserManager
	codec.Encoder
	linkedlist.Buffer
}

// NewForwardHandler 创建消息转发处理器
func NewForwardHandler(userManager user.IUserManager) *ForwardHandler {
	return &ForwardHandler{
		userManager: userManager,
		Buffer:      linkedlist.Buffer{},
	}
}

// Handle 实现消息转发逻辑
func (h *ForwardHandler) Handle(data []byte) (bool, error) {
	msg := new(types.Message)
	err := h.Decode(data, msg)
	if err != nil {
		return false, err
	}
	h.Buffer.Append(data)
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
	encoder      codec.Encoder
}

// NewStoreHandler 创建消息存储处理器
func NewStoreHandler(messageStore *stores.MessageStore, encoder codec.Encoder) *StoreHandler {
	return &StoreHandler{messageStore: messageStore, encoder: encoder}
}

// Handle 实现消息存储逻辑
func (h *StoreHandler) Handle(data []byte) (bool, error) {
	msg := new(types.Message)
	err := h.encoder.Decode(data, msg)
	if err != nil {
		return false, err
	}
	// 检查消息存储器是否已初始化
	if h.messageStore == nil {
		return false, errors.New("message store is not initialized")
	}

	// 将消息转换为数据库模型
	message := &models.Message{}
	message.FromTypes(msg)

	// 保存消息到数据库
	err = h.messageStore.CreateMessage(message)
	if err != nil {
		return false, err
	}

	// 继续处理链
	return true, nil
}

// NewMessageChain 创建默认的消息处理链
func NewMessageChain(userManager user.IUserManager, ms *stores.MessageStore, encoder codec.Encoder) *Chain {
	chain := NewChain()

	// 添加消息转发处理器
	chain.AddHandler(NewForwardHandler(userManager))

	// 添加消息存储处理器
	chain.AddHandler(NewStoreHandler(ms, encoder))

	return chain
}
