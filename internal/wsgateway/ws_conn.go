package wsgateway

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/woxQAQ/gim/internal/wsgateway/codec"
	"github.com/woxQAQ/gim/internal/wsgateway/types"
)

// WebSocketConn 实现LongConn接口的WebSocket连接.
type WebSocketConn struct {
	id         string
	platformID int32
	conn       *websocket.Conn
	state      types.ConnectionState
	stateMu    sync.RWMutex

	lastPingTime time.Time
	pingMu       sync.RWMutex

	// 消息处理
	compressor codec.Compressor
	encoder    codec.Encoder

	// 回调函数
	onMessage    func(types.Message)
	onDisconnect func(error)
	onError      func(error)

	// 用于优雅关闭
	closeOnce sync.Once
	closeChan chan struct{}
}

// NewWebSocketConn 创建新的WebSocket连接实例
func NewWebSocketConn(conn *websocket.Conn, id string, platformID int32) *WebSocketConn {
	return &WebSocketConn{
		id:           id,
		platformID:   platformID,
		conn:         conn,
		state:        types.Disconnected,
		lastPingTime: time.Now(),
		closeChan:    make(chan struct{}),
		onError:      func(err error) { fmt.Printf("WebSocket连接错误 [ID: %s]: %v\n", id, err) },
		onDisconnect: func(err error) { fmt.Printf("WebSocket连接断开 [ID: %s]: %v\n", id, err) },
	}
}

// Connect 实现LongConn接口的Connect方法
func (w *WebSocketConn) Connect(ctx context.Context) error {
	// WebSocket连接已经在HTTP升级时建立，这里只需要启动消息读取循环
	w.setConnectionState(types.Connected)
	go w.readPump()
	return nil
}

// Disconnect 实现LongConn接口的Disconnect方法
func (w *WebSocketConn) Disconnect(err error) error {
	w.closeOnce.Do(func() {
		w.setConnectionState(types.Closing)
		close(w.closeChan)

		// 关闭WebSocket连接
		if w.conn != nil {
			w.conn.Close()
		}

		w.setConnectionState(types.Disconnected)

		// 触发断开连接回调
		if w.onDisconnect != nil {
			w.onDisconnect(err)
		}
	})
	return nil
}

// Send 实现LongConn接口的Send方法
func (w *WebSocketConn) Send(msg types.Message) error {
	if w.State() != types.Connected {
		return errors.New("connection is not established")
	}

	return w.conn.WriteMessage(msg.Header.Type.MapMessageType(), msg.Payload)
}

// Receive 实现LongConn接口的Receive方法
func (w *WebSocketConn) Receive() (types.Message, error) {
	msgType, data, err := w.conn.ReadMessage()
	if err != nil {
		return types.Message{}, err
	}

	return types.Message{
		Header: types.MessageHeader{
			Type:      types.MessageType(msgType),
			Timestamp: time.Now(),
			From:      w.id,
			Platform:  w.platformID,
		},
		Payload: data,
	}, nil
}

// State 实现LongConn接口的State方法
func (w *WebSocketConn) State() types.ConnectionState {
	w.stateMu.RLock()
	defer w.stateMu.RUnlock()
	return w.state
}

// LastPingTime 实现LongConn接口的LastPingTime方法
func (w *WebSocketConn) LastPingTime() time.Time {
	w.pingMu.RLock()
	defer w.pingMu.RUnlock()
	return w.lastPingTime
}

// UpdateLastPingTime 实现LongConn接口的UpdateLastPingTime方法
func (w *WebSocketConn) UpdateLastPingTime(t time.Time) {
	w.pingMu.Lock()
	defer w.pingMu.Unlock()
	w.lastPingTime = t
}

// ID 实现LongConn接口的ID方法
func (w *WebSocketConn) ID() string {
	return w.id
}

// PlatformID 实现LongConn接口的PlatformID方法
func (w *WebSocketConn) PlatformID() int32 {
	return w.platformID
}

// OnMessage 实现LongConn接口的OnMessage方法
func (w *WebSocketConn) OnMessage(handler func(types.Message)) {
	w.onMessage = handler
}

// OnDisconnect 实现LongConn接口的OnDisconnect方法
func (w *WebSocketConn) OnDisconnect(handler func(error)) {
	w.onDisconnect = handler
}

// OnError 实现LongConn接口的OnError方法
func (w *WebSocketConn) OnError(handler func(error)) {
	w.onError = handler
}

// 内部方法

// setConnectionState 设置连接状态
func (w *WebSocketConn) setConnectionState(state types.ConnectionState) {
	w.stateMu.Lock()
	defer w.stateMu.Unlock()
	w.state = state
}

// readPump 持续读取WebSocket消息
func (w *WebSocketConn) readPump() {
	defer func() {
		if err := w.Disconnect(nil); err != nil {
			if w.onError != nil {
				w.onError(err)
			}
		}
	}()

	// 启动心跳检测
	go w.heartbeatChecker()

	for {
		select {
		case <-w.closeChan:
			return
		default:
			msg, err := w.Receive()
			if err != nil {
				if w.onError != nil {
					w.onError(err)
				}
				return
			}

			if w.onMessage != nil {
				w.onMessage(msg)
			}
		}
	}
}

// heartbeatChecker 心跳检测协程
func (w *WebSocketConn) heartbeatChecker() {
	ticker := time.NewTicker(30 * time.Second) // 每30秒检查一次
	defer ticker.Stop()

	var handleHeartbeatError = func(err error) {
		if w.onError != nil {
			w.onError(err)
		}
		_ = w.Disconnect(err) // 忽略Disconnect可能返回的错误，因为连接已经出现问题
	}

	for {
		select {
		case <-w.closeChan:
			return
		case <-ticker.C:
			// 发送ping消息
			if err := w.conn.WriteControl(websocket.PingMessage,
				[]byte("\n"), time.Now().Add(10*time.Second),
			); err != nil {
				handleHeartbeatError(err)
				return
			}

			// 检查最后一次心跳时间
			if time.Since(w.LastPingTime()) > 60*time.Second {
				handleHeartbeatError(errors.New("heartbeat timeout"))
				return
			}
		}
	}
}
