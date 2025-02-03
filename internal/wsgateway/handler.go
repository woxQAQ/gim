package wsgateway

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketHandler 处理WebSocket连接的处理器
type WebSocketHandler struct {
	upgrader websocket.Upgrader
}

var (
	instance *WebSocketHandler
	once     sync.Once
)

// GetWebSocketHandler 获取WebSocketHandler的单例实例
func GetWebSocketHandler() *WebSocketHandler {
	once.Do(func() {
		instance = &WebSocketHandler{
			upgrader: websocket.Upgrader{
				// 设置读写缓冲区大小
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				// 允许所有来源的跨域请求
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		}
	})
	return instance
}

// HandleWebSocket 处理WebSocket连接请求
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取必要的信息
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	// 获取平台ID，默认为1
	platformID := int32(1)

	// 升级HTTP连接为WebSocket连接
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade connection", http.StatusInternalServerError)
		return
	}

	// 设置WebSocket连接的基本配置
	conn.SetReadLimit(512)                                 // 设置最大消息大小
	conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // 设置读取超时
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// 创建WebSocket连接实例
	wsConn := NewWebSocketConn(conn, userID, platformID)

	// 启动连接
	if err := wsConn.Connect(r.Context()); err != nil {
		conn.Close()
		return
	}
}
