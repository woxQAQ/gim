package types

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/woxQAQ/gim/pkg/snowflake"
)

// MessageType 定义消息类型
type MessageType int32

const (
	// 系统消息类型
	MessageTypeUnknown   MessageType = iota
	MessageTypeHeartbeat             // 心跳消息
	MessageTypeSystem                // 系统消息

	// 业务消息类型
	MessageTypeText   // 文本消息
	MessageTypeImage  // 图片消息
	MessageTypeVideo  // 视频消息
	MessageTypeAudio  // 音频消息
	MessageTypeFile   // 文件消息
	MessageTypeCustom // 自定义消息
)

// mapMessageType 将内部消息类型映射到WebSocket消息类型
func (w MessageType) MapMessageType() int {
	switch w {
	case MessageTypeText, MessageTypeSystem, MessageTypeHeartbeat:
		return websocket.TextMessage
	case MessageTypeImage, MessageTypeVideo, MessageTypeAudio, MessageTypeFile, MessageTypeCustom:
		return websocket.BinaryMessage
	default:
		return websocket.TextMessage
	}
}

func (w MessageType) String() string {
	switch w {
	case MessageTypeText:
		return "text"
	case MessageTypeImage:
		return "image"
	case MessageTypeVideo:
		return "video"
	case MessageTypeAudio:
		return "audio"
	case MessageTypeFile:
		return "file"
	case MessageTypeCustom:
		return "custom"
	case MessageTypeSystem:
		return "system"
	case MessageTypeHeartbeat:
		return "heartbeat"
	default:
		return "unknown"
	}
}

// MessageHeader 定义消息头部结构
type MessageHeader struct {
	ID        string      // 消息唯一标识
	Type      MessageType // 消息类型
	Timestamp time.Time   // 消息时间戳
	From      string      // 发送者ID
	To        string      // 接收者ID
	Platform  int32       // 平台标识
}

// Message 定义新的消息结构
type Message struct {
	Header  MessageHeader // 消息头部
	Payload []byte        // 消息内容
}

// NewMessage 创建新的消息实例
func NewMessage(msgType MessageType, from, to string, platform int32, payload []byte) *Message {
	return &Message{
		Header: MessageHeader{
			ID:        generateMessageID(), // 这里需要实现一个生成消息ID的函数
			Type:      msgType,
			Timestamp: time.Now(),
			From:      from,
			To:        to,
			Platform:  platform,
		},
		Payload: payload,
	}
}

// generateMessageID 生成消息唯一标识
func generateMessageID() string {
	return snowflake.GenerateID()
}
