package base

import "time"

type IMessageType interface {
	String() string
	Int() int
}

// IMessage 定义消息接口
type IMessage interface {
	// GetID 获取消息ID
	GetID() string
	// GetType 获取消息类型
	GetType() IMessageType
	// GetTimestamp 获取消息时间戳
	GetTimestamp() time.Time
	// GetFrom 获取发送者ID
	GetFrom() string
	// GetTo 获取接收者ID
	GetTo() string
	// GetPlatform 获取平台标识
	GetPlatform() int32
	// GetPayload 获取消息内容
	GetPayload() []byte
}

// 确保Message实现了IMessage接口
