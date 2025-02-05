package types

import "github.com/gorilla/websocket"

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
