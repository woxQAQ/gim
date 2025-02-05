package wsgateway

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"github.com/woxQAQ/gim/internal/types"
	"github.com/woxQAQ/gim/internal/wsgateway/codec"
	"github.com/woxQAQ/gim/internal/wsgateway/handler"
	"github.com/woxQAQ/gim/internal/wsgateway/user"
	"github.com/woxQAQ/gim/pkg/db"
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
	Broadcast(msg types.Message) []error

	// SendToAllPlatforms 向指定用户的所有平台发送消息
	SendToAllPlatforms(userID string, msg types.Message) []error

	// SendToPlatform 向指定用户的指定平台发送消息
	SendToPlatform(userID string, platformID int32, msg types.Message) error

	// GetOnlineCount 获取当前在线用户数量
	GetOnlineCount() int

	// IsUserOnline 检查用户是否在线
	IsUserOnline(userID string) bool

	// GetUserHeartbeatStatus 获取指定用户在指定平台的最后心跳时间
	GetUserHeartbeatStatus(userID string, platformID int32) (time.Time, error)
}

// WSGateway 实现Gateway接口的WebSocket网关.
type WSGateway struct {
	// WebSocket配置
	upgrader websocket.Upgrader

	// 用户连接管理
	userManager user.IUserManager

	// 心跳检测配置
	heartbeatInterval time.Duration
	heartbeatTimeout  time.Duration

	// 消息编解码和压缩
	compressor codec.Compressor
	encoder    codec.Encoder

	// 消息处理链
	messageChain *handler.Chain

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

// WithCompressor 设置WSGateway的压缩器.
func WithCompressor(c codec.Compressor) Option {
	return func(g *WSGateway) {
		g.compressor = c
	}
}

// WithEncoder 设置WSGateway的编码器.
func WithEncoder(e codec.Encoder) Option {
	return func(g *WSGateway) {
		g.encoder = e
	}
}

func WithHeartbeat(interval, timeout time.Duration) Option {
	return func(g *WSGateway) {
		g.heartbeatInterval = interval
		g.heartbeatTimeout = timeout
	}
}

// NewWSGateway 创建新的WebSocket网关实例.
func NewWSGateway(opts ...Option) (*WSGateway, error) {
	g := &WSGateway{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		userManager:       user.NewUserManager(),
		heartbeatInterval: 30 * time.Second,
		heartbeatTimeout:  60 * time.Second,
		closedChan:        make(chan struct{}),
	}

	// 应用选项
	for _, opt := range opts {
		opt(g)
	}

	if g.logger == nil {
		// 使用默认配置创建logger，内部已包含fallback机制
		l, err := logger.NewLogger(logger.DomainWSGateway, &logger.Config{
			Level: "info",
		})
		if err != nil {
			return nil, err
		}
		g.logger = l
	}

	if g.encoder == nil {
		g.encoder = codec.NewJSONEncoder()
	}

	if g.compressor == nil {
		g.compressor = codec.NewGzipCompressor()
	}

	ms := stores.NewMessageStore(db.GetDB())

	// 初始化消息处理链
	g.messageChain = handler.NewMessageChain(g.userManager, ms)

	return g, nil
}

// Start 实现Gateway接口的Start方法.
func (g *WSGateway) Start(ctx context.Context) error {
	g.logger.Info("Starting WebSocket gateway service")
	// 使用传入的context创建一个可取消的context
	ctx, cancel := context.WithCancel(ctx)
	g.ctx = ctx
	g.cancel = cancel

	// 启动用户管理器
	g.logger.Info("Initializing user manager")
	if err := g.userManager.Start(ctx); err != nil {
		g.logger.Error("Failed to start user manager", logger.Error(err))
		return err
	}

	// 监听context取消信号
	go func() {
		<-ctx.Done()
		g.Stop()
	}()

	g.logger.Info("WebSocket gateway service started successfully")
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
func (g *WSGateway) Broadcast(msg types.Message) []error {
	g.logger.Info("Broadcasting message to all online users")
	errs := g.userManager.BroadcastMessage(msg)
	if len(errs) > 0 {
		for _, err := range errs {
			g.logger.Error("Failed to broadcast message", logger.Error(err))
		}
	}
	return errs
}

// SendToUser 实现Gateway接口的SendToUser方法.
// SendToAllPlatforms 向指定用户的所有平台发送消息.
func (g *WSGateway) SendToAllPlatforms(userID string, msg types.Message) []error {
	g.logger.Info("Sending message to all platforms of user", logger.String("user_id", userID))
	errs := g.userManager.SendMessage(userID, msg)
	if len(errs) > 0 {
		for _, err := range errs {
			g.logger.Error("Failed to send message to user", logger.String("user_id", userID), logger.Error(err))
		}
	}
	return errs
}

// SendToPlatform 向指定用户的指定平台发送消息.
func (g *WSGateway) SendToPlatform(userID string, platformID int32, msg types.Message) error {
	g.logger.Info("Sending message to specific platform",
		logger.String("user_id", userID),
		logger.Int32("platform_id", platformID))
	return g.userManager.SendPlatformMessage(userID, platformID, msg)
}

// GetOnlineCount 实现Gateway接口的GetOnlineCount方法.
func (g *WSGateway) GetOnlineCount() int {
	// 获取所有用户状态
	states, err := g.userManager.GetAll()
	if err != nil {
		g.logger.Error("Failed to get user states", logger.Error(err))
		return 0
	}

	// 统计在线用户数量
	count := 0
	for _, state := range states {
		if state != nil && len(state.OnlinePlatform) > 0 {
			count++
		}
	}

	g.logger.Info("Getting online user count", logger.Int("user_count", count))
	return count
}

// IsUserOnline 实现Gateway接口的IsUserOnline方法.
func (g *WSGateway) IsUserOnline(userID string) bool {
	state, err := g.userManager.GetState(userID)
	if err != nil {
		g.logger.Error("Failed to get user state", logger.Error(err))
		return false
	}
	isOnline := state != nil && len(state.OnlinePlatform) > 0
	g.logger.Info("Checking user online status", logger.String("user_id", userID), logger.Bool("is_online", isOnline))
	return isOnline
}

// GetUserHeartbeatStatus 实现Gateway接口的GetUserHeartbeatStatus方法.
func (g *WSGateway) GetUserHeartbeatStatus(userID string, platformID int32) (time.Time, error) {
	g.logger.Info("Getting user heartbeat status",
		logger.String("user_id", userID),
		logger.Int32("platform_id", platformID))

	// 获取用户连接
	conn, err := g.userManager.GetConn(userID, platformID)
	if err != nil {
		g.logger.Error("Failed to get user connection", logger.Error(err))
		return time.Time{}, err
	}

	if conn == nil {
		return time.Time{}, fmt.Errorf("user %s platform %d not connected", userID, platformID)
	}

	// 获取最后心跳时间
	lastPingTime := conn.LastPingTime()
	g.logger.Info("Got user heartbeat status",
		logger.String("user_id", userID),
		logger.Int32("platform_id", platformID),
		logger.Time("last_ping_time", lastPingTime))

	return lastPingTime, nil
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
	wsConn.OnMessage(func(msg types.Message) {
		// 心跳消息特殊处理
		if msg.Header.Type == types.MessageTypeHeartbeat {
			g.logger.Debug("Received heartbeat from user",
				logger.String("user_id", userID),
				logger.Int32("platform_id", platformID))
			wsConn.UpdateLastPingTime(time.Now())
			return
		}

		// 使用责任链处理其他类型的消息
		if err := g.messageChain.Process(msg); err != nil {
			g.logger.Error("Failed to process message",
				logger.String("user_id", userID),
				logger.Error(err))
		}
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

	g.logger.Info("New WebSocket connection established",
		logger.String("user_id", userID),
		logger.Int32("platform_id", platformID),
	)
}
