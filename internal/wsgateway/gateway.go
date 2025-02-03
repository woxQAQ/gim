package wsgateway

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/woxQAQ/gim/internal/wsgateway/user"
	"github.com/woxQAQ/gim/pkg/logger"
)

var _ Gateway = &WSGateway{}

// Gateway 定义消息网关接口.
type Gateway interface {
	// Start 启动网关服务.
	Start(ctx context.Context) error

	// Stop 停止网关服务
	Stop() error

	// Broadcast 广播消息给所有在线用户
	Broadcast(msg Message) error

	// SendToUser 发送消息给指定用户
	SendToUser(userID string, msg Message) error

	// GetOnlineCount 获取当前在线用户数量
	GetOnlineCount() int

	// IsUserOnline 检查用户是否在线
	IsUserOnline(userID string) bool
}

// WSGateway 实现Gateway接口的WebSocket网关.
type WSGateway struct {
	// WebSocket处理器
	wsHandler *WebSocketHandler

	// 用户连接管理
	userManager user.IUserManager

	// 心跳检测配置
	heartbeatInterval time.Duration
	heartbeatTimeout  time.Duration

	// 日志记录器
	logger logger.Logger

	// 关闭控制
	ctx        context.Context
	cancel     context.CancelFunc
	closeOnce  sync.Once
	closedChan chan struct{}
}

// Option 定义WSGateway的配置选项函数类型.
type Option func(*WSGateway)

// WithLogger 设置WSGateway的logger.
func WithLogger(l logger.Logger) Option {
	return func(g *WSGateway) {
		g.logger = l
	}
}

// NewWSGateway 创建新的WebSocket网关实例.
func NewWSGateway(opts ...Option) (*WSGateway, error) {
	ctx, cancel := context.WithCancel(context.Background())
	g := &WSGateway{
		wsHandler:         GetWebSocketHandler(),
		userManager:       user.NewUserManager(),
		heartbeatInterval: 30 * time.Second,
		heartbeatTimeout:  60 * time.Second,
		ctx:               ctx,
		cancel:            cancel,
		closedChan:        make(chan struct{}),
	}

	// 应用选项
	for _, opt := range opts {
		opt(g)
	}

	if g.logger == nil {
		l, err := logger.NewLogger("ws_gateway", nil)
		if err != nil {
			// 使用默认的控制台logger作为fallback
			defaultLogger, fallbackErr := logger.NewLogger("ws_gateway", &logger.Config{
				Level: "info",
			})
			if fallbackErr != nil {
				return nil, fallbackErr
			}
			g.logger = defaultLogger
			g.logger.Warn("Failed to initialize configured logger, using default console logger", logger.Error(err))
		} else {
			g.logger = l
		}
	}

	return g, nil
}

// Start 实现Gateway接口的Start方法.
func (g *WSGateway) Start(ctx context.Context) error {
	g.logger.Info("Starting WebSocket gateway service")
	return nil
}

// Stop 实现Gateway接口的Stop方法.
func (g *WSGateway) Stop() error {
	g.closeOnce.Do(func() {
		// 取消上下文
		g.cancel()

		// 标记关闭完成
		close(g.closedChan)
		g.logger.Info("WebSocket gateway service stopped")
	})
	return nil
}

// Broadcast 实现Gateway接口的Broadcast方法.
func (g *WSGateway) Broadcast(msg Message) error {
	g.logger.Info("Broadcasting message to all online users")
	return g.userManager.BroadcastMessage(msg.Type, msg.Payload)
}

// SendToUser 实现Gateway接口的SendToUser方法.
func (g *WSGateway) SendToUser(userID string, msg Message) error {
	g.logger.Info("Sending message to user", logger.String("user_id", userID))
	return g.userManager.SendMessage(userID, msg.Type, msg.Payload)
}

// GetOnlineCount 实现Gateway接口的GetOnlineCount方法.
func (g *WSGateway) GetOnlineCount() int {
	count := 0
	// TODO: 实现获取在线用户数量的逻辑
	g.logger.Info("Getting online user count", logger.Int("user_count", count))
	return count
}

// IsUserOnline 实现Gateway接口的IsUserOnline方法.
func (g *WSGateway) IsUserOnline(userID string) bool {
	state, err := g.userManager.GetUserState(userID)
	if err != nil {
		g.logger.Error("Failed to get user state", logger.Error(err))
		return false
	}
	isOnline := state != nil && len(state.OnlinePlatform) > 0
	g.logger.Info("Checking user online status", logger.String("user_id", userID), logger.Bool("is_online", isOnline))
	return isOnline
}

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

	// 创建WebSocket连接实例
	wsConn := NewWebSocketConn(nil, userID, platformID)

	// 设置连接回调
	wsConn.OnMessage(func(msg Message) {
		// 处理接收到的消息
		// TODO: 实现消息处理逻辑
		g.logger.Info("Received message from user", logger.String("user_id", userID), logger.String("message", string(msg.Payload)))
	})

	// 添加连接到用户管理器
	if err := g.userManager.AddUserConn(userID, platformID, wsConn); err != nil {
		g.logger.Error("Failed to add connection to user manager", logger.Error(err))
		http.Error(w, "failed to add connection", http.StatusInternalServerError)
		return
	}

	g.logger.Info("New WebSocket connection established",
		logger.String("user_id", userID), logger.Int32("platform_id", platformID))

	// 处理WebSocket连接
	g.wsHandler.HandleWebSocket(w, r)
}
