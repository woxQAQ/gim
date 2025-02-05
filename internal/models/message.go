package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/woxQAQ/gim/internal/apiserver/types/response"
	"github.com/woxQAQ/gim/internal/types"
)

// MessageStatus 消息状态
type MessageStatus int8

const (
	// MessageStatusUnknown 未知状态
	MessageStatusUnknown MessageStatus = iota
	// MessageStatusSent 已发送
	MessageStatusSent
	// MessageStatusDelivered 已投递
	MessageStatusDelivered
	// MessageStatusRead 已读
	MessageStatusRead
	// MessageStatusFailed 发送失败
	MessageStatusFailed
)

// Message 消息模型
type Message struct {
	ID        string            `gorm:"primaryKey;type:text"`
	Type      types.MessageType `gorm:"type:smallint;not null;index"`
	Content   string            `gorm:"type:text;not null"`
	FromID    string            `gorm:"type:text;not null;index"`
	ToID      string            `gorm:"type:text;not null;index"`
	Status    MessageStatus     `gorm:"type:smallint;not null;default:1;index"`
	Platform  int32             `gorm:"type:integer;not null;index"`
	CreatedAt time.Time         `gorm:"autoCreateTime;index"`
	UpdatedAt time.Time         `gorm:"autoUpdateTime"`
}

func (m *Message) TableName() string {
	return "messages"
}

// BeforeCreate 在创建消息前生成消息ID
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	// TODO: 实现消息ID生成逻辑
	return nil
}

// ToResponse 将Message转换为MessageResponse
func (m *Message) ToResponse() *response.MessageResponse {
	return &response.MessageResponse{
		ID:         m.ID,
		SenderID:   m.FromID,
		ReceiverID: m.ToID,
		Content:    m.Content,
		Type:       int32(m.Type),
		Status:     int32(m.Status),
		CreatedAt:  m.CreatedAt,
	}
}

func (m *Message) FromTypes(msg types.Message) {
	m.ID = msg.Header.ID
	m.FromID = msg.Header.From
	m.ToID = msg.Header.To
	m.Content = string(msg.Payload)
	m.Platform = msg.Header.Platform
	m.CreatedAt = msg.Header.Timestamp
}

// MessageAttachment 消息附件模型
type MessageAttachment struct {
	ID        string    `gorm:"primaryKey;type:text"`
	MessageID string    `gorm:"type:text;not null;index"`
	Type      string    `gorm:"type:text;not null"` // 附件类型：image, video, audio, file
	URL       string    `gorm:"type:text;not null"` // 附件URL
	Size      int64     `gorm:"type:bigint"`        // 附件大小（字节）
	MimeType  string    `gorm:"type:text"`          // MIME类型
	Metadata  string    `gorm:"type:text"`          // 附件元数据（JSON格式）
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (ma *MessageAttachment) TableName() string {
	return "message_attachments"
}
