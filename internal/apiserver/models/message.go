package models

import (
	"time"

	"gorm.io/gorm"
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

// MessageType 消息类型
type MessageType int8

const (
	// MessageTypeUnknown 未知类型
	MessageTypeUnknown MessageType = iota
	// MessageTypeText 文本消息
	MessageTypeText
	// MessageTypeImage 图片消息
	MessageTypeImage
	// MessageTypeVideo 视频消息
	MessageTypeVideo
	// MessageTypeAudio 音频消息
	MessageTypeAudio
	// MessageTypeFile 文件消息
	MessageTypeFile
	// MessageTypeCustom 自定义消息
	MessageTypeCustom
)

// Message 消息模型
type Message struct {
	ID        string        `gorm:"primaryKey;type:text"`
	Type      MessageType   `gorm:"type:smallint;not null;index"`
	Content   string        `gorm:"type:text;not null"`
	FromID    string        `gorm:"type:text;not null;index"`
	ToID      string        `gorm:"type:text;not null;index"`
	Status    MessageStatus `gorm:"type:smallint;not null;default:1;index"`
	Platform  int32         `gorm:"type:integer;not null;index"`
	CreatedAt time.Time     `gorm:"autoCreateTime;index"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`
}

func (m *Message) TableName() string {
	return "messages"
}

// BeforeCreate 在创建消息前生成消息ID
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	// TODO: 实现消息ID生成逻辑
	return nil
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
