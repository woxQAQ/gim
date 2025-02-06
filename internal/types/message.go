package types

import (
	"time"

	"github.com/woxQAQ/gim/internal/wsgateway/base"
	"github.com/woxQAQ/gim/pkg/snowflake"
)

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

var _ base.IMessage = (*Message)(nil)

// 实现Message的接口方法
func (m *Message) GetID() string {
	return m.Header.ID
}

func (m *Message) GetType() base.IMessageType {
	return m.Header.Type
}

func (m *Message) GetTimestamp() time.Time {
	return m.Header.Timestamp
}

func (m *Message) GetFrom() string {
	return m.Header.From
}

func (m *Message) GetTo() string {
	return m.Header.To
}

func (m *Message) GetPlatform() int32 {
	return m.Header.Platform
}

func (m *Message) GetPayload() []byte {
	return m.Payload
}
