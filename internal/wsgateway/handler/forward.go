package handler

import (
	"context"
	"fmt"

	"github.com/woxQAQ/gim/internal/types"
	"github.com/woxQAQ/gim/internal/wsgateway/codec"
	"github.com/woxQAQ/gim/internal/wsgateway/user"
	"github.com/woxQAQ/gim/pkg/mq"
)

// ForwardHandler 消息转发处理器
type ForwardHandler struct {
	BaseHandler
	userManager user.IUserManager
	producer    mq.Producer
	encoder     codec.Encoder
	topic       string
}

// NewForwardHandler 创建消息转发处理器
func NewForwardHandler(userManager user.IUserManager, producer mq.Producer, encoder codec.Encoder) *ForwardHandler {
	return &ForwardHandler{
		userManager: userManager,
		producer:    producer,
		encoder:     encoder,
		topic:       "message_forward", // 可配置的转发主题
	}
}

// Handle 实现消息转发逻辑
func (h *ForwardHandler) Handle(data []byte) (bool, error) {
	msg := new(types.Message)
	if err := h.encoder.Decode(data, msg); err != nil {
		return false, err
	}

	// 检查消息类型是否需要转发
	switch msg.GetType() {
	case types.MessageTypeText, types.MessageTypeImage,
		types.MessageTypeVideo, types.MessageTypeAudio,
		types.MessageTypeFile:
		// 将消息发送到MQ
		if msg.GetTo() != "" {
			if err := h.publishToMQ(msg); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (h *ForwardHandler) publishToMQ(msg *types.Message) error {
	// 序列化消息
	data, err := h.encoder.Encode(msg)
	if err != nil {
		return err
	}

	// 发布到MQ
	return h.producer.Publish(context.Background(), &mq.Message{
		Topic:   h.topic,
		Key:     fmt.Sprintf("message:%s:%s", msg.GetFrom(), msg.GetTo()), // 使用接收者ID作为key
		Value:   data,
		Headers: map[string]string{"type": msg.GetType().String()},
	})
}
