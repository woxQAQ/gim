package response

import (
	"time"
)

// MessageResponse 消息响应
type MessageResponse struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	Type       int32     `json:"type"`
	Status     int32     `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// MessageHistoryResponse 消息历史记录响应
type MessageHistoryResponse struct {
	Messages  []*MessageResponse `json:"messages"`
	NextToken string             `json:"next_token,omitempty"`
}

// UnreadCountResponse 未读消息数量响应
type UnreadCountResponse struct {
	Count int `json:"count"`
}
