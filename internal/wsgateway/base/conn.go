package base

import (
	"context"
	"time"
)

// ConnectionState 定义连接状态。
type ConnectionState int

const (
	Disconnected ConnectionState = iota
	Connecting
	Connected
	Closing
)

// LongConn 定义长连接接口。
type LongConn interface {
	// Connect 建立连接
	Connect(ctx context.Context) error

	// Disconnect 主动断开连接
	Disconnect(err error) error

	// Send 发送消息
	Send(int, []byte) error

	// Receive 接收消息
	Receive() (int, []byte, error)

	// State 获取当前连接状态
	State() ConnectionState

	// LastPingTime 获取最后一次心跳时间
	LastPingTime() time.Time

	// UpdateLastPingTime 更新最后一次心跳时间
	UpdateLastPingTime(t time.Time)

	// ID 获取连接唯一标识
	ID() string

	// PlatformID 获取平台标识
	PlatformID() int32

	// OnMessage 设置消息处理回调
	OnMessage(handler func(int, []byte))

	// OnDisconnect 设置连接断开回调
	OnDisconnect(handler func(error))

	// OnError 设置错误处理回调
	OnError(handler func(error))
}
