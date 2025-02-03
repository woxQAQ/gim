package wsgateway

import (
	"github.com/woxQAQ/gim/pkg/types"
)

// LongConn 定义长连接接口.
type LongConn = types.LongConn

// Message 定义消息的基本结构.
type Message = types.Message

// ConnectionState 定义连接状态.
type ConnectionState = types.ConnectionState

// 连接状态常量.
const (
	Disconnected = types.Disconnected
	Connecting   = types.Connecting
	Connected    = types.Connected
	Closing      = types.Closing
)
