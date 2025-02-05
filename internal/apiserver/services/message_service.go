package services

import (
	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"github.com/woxQAQ/gim/internal/apiserver/types/response"
)

// MessageService 处理消息相关的业务逻辑
type MessageService struct {
	messageStore *stores.MessageStore
}

// NewMessageService 创建MessageService实例
func NewMessageService(messageStore *stores.MessageStore) *MessageService {
	return &MessageService{
		messageStore: messageStore,
	}
}

// GetMessageHistory 获取消息历史记录
func (s *MessageService) GetMessageHistory(userID string, pageSize int, lastMessageID string) (*response.MessageHistoryResponse, error) {
	// 获取消息列表
	messages, err := s.messageStore.GetMessagesByUserID(userID, pageSize+1, lastMessageID)
	if err != nil {
		return nil, err
	}

	// 构建响应
	response := &response.MessageHistoryResponse{
		Messages: make([]*response.MessageResponse, 0, len(messages)),
	}

	// 如果获取到的消息数量超过pageSize，说明还有更多消息
	if len(messages) > pageSize {
		response.NextToken = messages[len(messages)-1].ID
		messages = messages[:pageSize] // 只返回请求的数量
	}

	// 转换消息格式
	for _, msg := range messages {
		response.Messages = append(response.Messages, msg.ToResponse())
	}

	return response, nil
}
