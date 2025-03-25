package user

import (
	"sync"
	"time"

	"github.com/woxQAQ/gim/internal/wsgateway/base"
	"github.com/woxQAQ/gim/pkg/workerpool"
)

type IUserConnManager interface {
	// AddConn 添加用户的平台连接
	AddConn(userID string, platformID int32, conn base.LongConn) error
	// RemoveConn 移除用户的平台连接
	RemoveConn(userID string, platformID int32) error
	// GetConn 获取用户在指定平台的连接
	GetConn(userID string, platformID int32) (base.LongConn, error)
	// GetState 获取用户在各平台的在线状态
	GetState(userID string) (*State, error)
	// GetAll 获取所有用户的状态信息
	GetAll() ([]*State, error)
	// GetOnlineCount 获取当前在线用户数量
	GetOnlineCount() int
	// IsOnline 检查用户是否在线
	IsOnline(userID string) bool
}

type IMessageSender interface {
	// BroadcastMessage 向所有用户的所有平台广播消息
	BroadcastMessage(msg base.IMessage) []error
	// SendMessage 向指定用户的所有平台发送消息
	SendMessage(userID string, msg base.IMessage) []error
	// SendPlatformMessage 向指定用户的指定平台发送消息
	SendPlatformMessage(userID string, platformID int32, msg base.IMessage) error
}

// IUserManager 定义用户管理器的接口.
type IUserManager interface {
	IUserConnManager
	IMessageSender
}

var _ IUserManager = &Manager{}

// UserManager 管理所有用户的连接.
type Manager struct {
	users     map[string]*Platform // 用户ID到用户平台管理器的映射
	mutex     sync.RWMutex         // 用于并发安全的读写锁
	observers []StateObserver      // 状态观察者列表
}

// NewUserManager 创建新的用户管理器实例.
func NewUserManager() *Manager {
	m := &Manager{
		users:     make(map[string]*Platform),
		observers: make([]StateObserver, 0),
	}

	return m
}

// AddObserver 添加状态观察者.
func (um *Manager) AddObserver(observer StateObserver) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.observers = append(um.observers, observer)
}

// RemoveObserver 移除状态观察者.
func (um *Manager) RemoveObserver(observer StateObserver) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	for i, obs := range um.observers {
		if obs == observer {
			um.observers = append(um.observers[:i], um.observers[i+1:]...)
			break
		}
	}
}

// notifyStateChange 通知所有观察者状态变化.
func (um *Manager) notifyStateChange(userID string, platformID int32, oldState, newState base.ConnectionState) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	timestamp := time.Now()
	for _, observer := range um.observers {
		// 使用协程池处理观察者通知
		observer := observer // 创建副本以避免闭包问题
		workerpool.GetInstance().Submit(func() {
			observer.OnUserStateChange(userID, platformID, oldState, newState, timestamp)
		})
	}
}
