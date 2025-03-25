package wsgateway

import (
	"net/http"
	"time"

	"github.com/woxQAQ/gim/pkg/logger"
)

// HandleNewConnection 处理新的WebSocket连接.
func (g *WSGateway) HandleNewConnection(w http.ResponseWriter, r *http.Request) {
	// 获取用户ID
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		g.logger.Error("Missing user_id in connection request")
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	// 获取平台ID，默认为1
	platformID := int32(1)

	// 升级HTTP连接为WebSocket连接
	conn, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		g.logger.Error("Failed to upgrade connection", logger.Error(err))
		http.Error(w, "could not upgrade connection", http.StatusInternalServerError)
		return
	}

	// 设置WebSocket连接的基本配置
	conn.SetReadLimit(512) // 设置最大消息大小
	if err = conn.SetReadDeadline(time.Now().Add(g.heartbeatTimeout)); err != nil {
		g.logger.Error("Failed to set read deadline", logger.Error(err))
		conn.Close()
		return
	}

	// 设置心跳处理
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(g.heartbeatTimeout))
	})

	// 创建WebSocket连接实例
	wsConn := NewWebSocketConn(conn, userID, platformID)
	// 设置压缩器和编码器
	wsConn.compressor = g.compressor
	wsConn.encoder = g.encoder

	// 设置连接回调
	wsConn.OnMessage(func(msgType int, data []byte) {
		// 心跳消息特殊处理
		g.messageHandler[msgType](wsConn, data)
	})

	// 启动连接
	if err := wsConn.Connect(g.ctx); err != nil {
		g.logger.Error("Failed to start connection", logger.Error(err))
		conn.Close()
		return
	}

	// 添加连接到用户管理器
	if err := g.userManager.AddConn(userID, platformID, wsConn); err != nil {
		g.logger.Error("Failed to add connection to user manager", logger.Error(err))
		http.Error(w, "failed to add connection", http.StatusInternalServerError)
		conn.Close()
		return
	}

	g.logger.Debug("New WebSocket connection established",
		logger.String("user_id", userID),
		logger.Int32("platform_id", platformID),
	)
}
