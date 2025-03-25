package wsgateway

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"github.com/woxQAQ/gim/internal/wsgateway/base"
	"github.com/woxQAQ/gim/internal/wsgateway/codec"
	"github.com/woxQAQ/gim/internal/wsgateway/handler"
	"github.com/woxQAQ/gim/internal/wsgateway/user"
	"github.com/woxQAQ/gim/pkg/db"
	"github.com/woxQAQ/gim/pkg/logger"
	"github.com/woxQAQ/gim/pkg/mq"
)

var _ Gateway = &WSGateway{}

// Gateway 定义消息网关接口.
type Gateway interface {
	// Start 启动网关服务.
	Start(ctx context.Context) error

	// Stop 停止网关服务
	Stop() error

	// Broadcast 广播消息给所有在线用户
	Broadcast(base.IMessage) []error

	// SendToAllPlatforms 向指定用户的所有平台发送消息
	SendToAllPlatforms(string, base.IMessage) []error

	// SendToPlatform 向指定用户的指定平台发送消息
	SendToPlatform(userID string, platformID int32, msg base.IMessage) error

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
		l, err := logger.NewLogger(&logger.Config{
			Level: "info",
		})
		if err != nil {
			return nil, err
		}
		g.logger = l.With(logger.String("domain", string(logger.DomainWSGateway)))
	}

	if g.encoder == nil {
		g.encoder = codec.NewJSONEncoder()
	}

	if g.compressor == nil {
		g.compressor = codec.NewGzipCompressor()
	}

	ms := stores.NewMessageStore(db.GetDB())

	// 初始化MQ
	mqFactory := mq.NewMemoryMQFactory(100)
	producer, err := mqFactory.NewProducer(&mq.Config{})
	if err != nil {
		return nil, err
	}

	consumer, err := mqFactory.NewConsumer(&mq.Config{})
	if err != nil {
		return nil, err
	}

	consumer.Subscribe("message_forward", func(ctx context.Context, msg *mq.Message) error {
		// 反序列化消息

	})

	// 初始化消息处理链
	g.messageChain = handler.NewMessageChain(g.ctx,
		g.userManager, ms,
		g.encoder, producer,
	)

	return g, nil
}

// Start 实现Gateway接口的Start方法.
func (g *WSGateway) Start(ctx context.Context) error {
	g.logger.Info("Starting WebSocket gateway service")
	// 使用传入的context创建一个可取消的context
	ctx, cancel := context.WithCancel(ctx)
	g.ctx = ctx
	g.cancel = cancel

	// 监听context取消信号
	go func() {
		<-ctx.Done()
		_ = g.Stop()
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
func (g *WSGateway) Broadcast(msg base.IMessage) []error {
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
func (g *WSGateway) SendToAllPlatforms(userID string, msg base.IMessage) []error {
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
func (g *WSGateway) SendToPlatform(userID string, platformID int32, msg base.IMessage) error {
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
