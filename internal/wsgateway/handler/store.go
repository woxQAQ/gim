package handler

import (
	"errors"

	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"github.com/woxQAQ/gim/internal/models"
	"github.com/woxQAQ/gim/internal/types"
	"github.com/woxQAQ/gim/internal/wsgateway/codec"
)

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
