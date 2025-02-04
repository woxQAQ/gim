package user

import (
	"context"
	"sync"
	"time"

	"github.com/woxQAQ/gim/internal/wsgateway/types"
)

// IUserManager 定义用户管理器的接口.
type IUserManager interface {
	// Start 启动用户管理器
	Start(ctx context.Context) error

	// AddConn 添加用户的平台连接
	AddConn(userID string, platformID int32, conn types.LongConn) error
	// RemoveConn 移除用户的平台连接
	RemoveConn(userID string, platformID int32) error
	// GetConn 获取用户在指定平台的连接
	GetConn(userID string, platformID int32) (types.LongConn, error)
	// GetState 获取用户在各平台的在线状态
	GetState(userID string) (*State, error)
	// GetAll 获取所有用户的状态信息
	GetAll() ([]*State, error)
	// GetOnlineCount 获取当前在线用户数量
	GetOnlineCount() int
	// IsOnline 检查用户是否在线
	IsOnline(userID string) bool
	// BroadcastMessage 向所有用户的所有平台广播消息
	BroadcastMessage(msg types.Message) []error
	// SendMessage 向指定用户的所有平台发送消息
	SendMessage(userID string, msg types.Message) []error
	// SendPlatformMessage 向指定用户的指定平台发送消息
	SendPlatformMessage(userID string, platformID int32, msg types.Message) error
}

var _ IUserManager = &Manager{}

// UserManager 管理所有用户的连接.
type Manager struct {
	users     map[string]*Platform // 用户ID到用户平台管理器的映射
	mutex     sync.RWMutex         // 用于并发安全的读写锁
	observers []StateObserver      // 状态观察者列表
}

// Start 实现 IUserManager 接口.
func (um *Manager) Start(ctx context.Context) error {
	// 目前用户管理器不需要特殊的启动逻辑
	// 如果将来需要添加后台任务或资源初始化，可以在这里实现
	return nil
}

// NewUserManager 创建新的用户管理器实例.
func NewUserManager() *Manager {
	return &Manager{
		users:     make(map[string]*Platform),
		observers: make([]StateObserver, 0),
	}
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
func (um *Manager) notifyStateChange(userID string, platformID int32, oldState, newState types.ConnectionState) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	timestamp := time.Now()
	for _, observer := range um.observers {
		go observer.OnUserStateChange(userID, platformID, oldState, newState, timestamp)
	}
}

// AddConn 实现 IUserManager 接口.
func (um *Manager) AddConn(userID string, platformID int32, conn types.LongConn) error {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	up, exists := um.users[userID]
	if !exists {
		up = NewUserPlatform(userID)
		um.users[userID] = up
	}

	up.mutex.Lock()
	defer up.mutex.Unlock()

	// 获取旧连接的状态（如果存在）
	var oldState types.ConnectionState = types.Disconnected
	if oldConn, exists := up.Conns[platformID]; exists {
		oldState = oldConn.State()
	}

	up.Conns[platformID] = conn

	// 通知观察者状态变化.
	um.notifyStateChange(userID, platformID, oldState, conn.State())

	// 设置连接断开回调.
	conn.OnDisconnect(func(err error) {
		um.notifyStateChange(userID, platformID, types.Connected, types.Disconnected)
	})

	return nil
}

// RemoveConn 实现 IUserManager 接口.
func (um *Manager) RemoveConn(userID string, platformID int32) error {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	up, exists := um.users[userID]
	if !exists {
		return nil
	}

	up.mutex.Lock()
	defer up.mutex.Unlock()
	delete(up.Conns, platformID)

	// 如果用户没有任何平台连接，则删除该用户.
	if len(up.Conns) == 0 {
		delete(um.users, userID)
	}
	return nil
}

// GetConn 实现 IUserManager 接口.
func (um *Manager) GetConn(userID string, platformID int32) (types.LongConn, error) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	up, exists := um.users[userID]
	if !exists {
		return nil, nil
	}

	up.mutex.RLock()
	defer up.mutex.RUnlock()
	conn, exists := up.Conns[platformID]
	if !exists {
		return nil, nil
	}
	return conn, nil
}

// GetState 实现 IUserManager 接口.
func (um *Manager) GetState(userID string) (*State, error) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	state := &State{Id: userID}
	up, exists := um.users[userID]
	if !exists {
		return state, nil
	}

	up.mutex.RLock()
	defer up.mutex.RUnlock()
	for platformID, conn := range up.Conns {
		if conn.State() == types.Connected {
			state.OnlinePlatform = append(state.OnlinePlatform, platformID)
		} else {
			state.OfflinePlatform = append(state.OfflinePlatform, platformID)
		}
	}
	return state, nil
}

// GetOnlineCount 实现 IUserManager 接口.
func (um *Manager) GetOnlineCount() int {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	return len(um.users)
}

// IsOnline 实现 IUserManager 接口.
func (um *Manager) IsOnline(userID string) bool {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	_, exists := um.users[userID]
	return exists
}

// BroadcastMessage 实现 IUserManager 接口.
func (um *Manager) BroadcastMessage(msg types.Message) []error {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	var errors []error
	for _, up := range um.users {
		up.mutex.RLock()
		for _, conn := range up.Conns {
			if conn.State() == types.Connected {
				if err := conn.Send(msg); err != nil {
					errors = append(errors, err)
				}
			}
		}
		up.mutex.RUnlock()
	}
	return errors
}

// SendMessage 实现 IUserManager 接口.
func (um *Manager) SendMessage(userID string, msg types.Message) []error {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	var errors []error
	up, exists := um.users[userID]
	if !exists {
		return nil
	}

	up.mutex.RLock()
	defer up.mutex.RUnlock()
	for _, conn := range up.Conns {
		if conn.State() == types.Connected {
			if err := conn.Send(msg); err != nil {
				errors = append(errors, err)
			}
		}
	}
	return errors
}

// SendPlatformMessage 实现 IUserManager 接口.
func (um *Manager) SendPlatformMessage(userID string, platformID int32, msg types.Message) error {
	conn, err := um.GetConn(userID, platformID)
	if err != nil || conn == nil {
		return err
	}
	return conn.Send(msg)
}

// GetAll 实现 IUserManager 接口.
func (um *Manager) GetAll() ([]*State, error) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	userStates := make([]*State, 0, len(um.users))
	for userID := range um.users {
		state, err := um.GetState(userID)
		if err != nil {
			continue
		}
		userStates = append(userStates, state)
	}
	return userStates, nil
}
