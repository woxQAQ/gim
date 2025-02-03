package user

import (
	"sync"
	"time"

	"github.com/woxQAQ/gim/pkg/types"
)

// IUserManager 定义用户管理器的接口
type IUserManager interface {
	// AddUserConn 添加用户的平台连接
	AddUserConn(userID string, platformID int32, conn types.LongConn) error
	// RemoveUserConn 移除用户的平台连接
	RemoveUserConn(userID string, platformID int32) error
	// GetUserConn 获取用户在指定平台的连接
	GetUserConn(userID string, platformID int32) (types.LongConn, error)
	// GetUserState 获取用户在各平台的在线状态
	GetUserState(userID string) (*UserState, error)
	// GetAllUsers 获取所有用户的状态信息
	GetAllUsers() ([]*UserState, error)
	// BroadcastMessage 向所有用户的所有平台广播消息
	BroadcastMessage(messageType int, data []byte) error
	// SendMessage 向指定用户的所有平台发送消息
	SendMessage(userID string, messageType int, data []byte) error
	// SendPlatformMessage 向指定用户的指定平台发送消息
	SendPlatformMessage(userID string, platformID int32, messageType int, data []byte) error
}

// UserManager 管理所有用户的连接
type UserManager struct {
	users     map[string]*UserPlatform // 用户ID到用户平台管理器的映射
	mutex     sync.RWMutex             // 用于并发安全的读写锁
	observers []UserStateObserver      // 状态观察者列表
}

// NewUserManager 创建新的用户管理器实例
func NewUserManager() *UserManager {
	return &UserManager{
		users:     make(map[string]*UserPlatform),
		observers: make([]UserStateObserver, 0),
	}
}

// AddObserver 添加状态观察者
func (um *UserManager) AddObserver(observer UserStateObserver) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	um.observers = append(um.observers, observer)
}

// RemoveObserver 移除状态观察者
func (um *UserManager) RemoveObserver(observer UserStateObserver) {
	um.mutex.Lock()
	defer um.mutex.Unlock()
	for i, obs := range um.observers {
		if obs == observer {
			um.observers = append(um.observers[:i], um.observers[i+1:]...)
			break
		}
	}
}

// notifyStateChange 通知所有观察者状态变化
func (um *UserManager) notifyStateChange(userID string, platformID int32, oldState, newState types.ConnectionState) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()
	timestamp := time.Now()
	for _, observer := range um.observers {
		go observer.OnUserStateChange(userID, platformID, oldState, newState, timestamp)
	}
}

// AddUserConn 实现 IUserManager 接口
func (um *UserManager) AddUserConn(userID string, platformID int32, conn types.LongConn) error {
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

	// 通知观察者状态变化
	um.notifyStateChange(userID, platformID, oldState, conn.State())

	// 设置连接断开回调
	conn.OnDisconnect(func(err error) {
		um.notifyStateChange(userID, platformID, types.Connected, types.Disconnected)
	})

	return nil
}

// RemoveUserConn 实现 IUserManager 接口
func (um *UserManager) RemoveUserConn(userID string, platformID int32) error {
	um.mutex.Lock()
	defer um.mutex.Unlock()

	up, exists := um.users[userID]
	if !exists {
		return nil
	}

	up.mutex.Lock()
	defer up.mutex.Unlock()
	delete(up.Conns, platformID)

	// 如果用户没有任何平台连接，则删除该用户
	if len(up.Conns) == 0 {
		delete(um.users, userID)
	}
	return nil
}

// GetUserConn 实现 IUserManager 接口
func (um *UserManager) GetUserConn(userID string, platformID int32) (types.LongConn, error) {
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

// GetUserState 实现 IUserManager 接口
func (um *UserManager) GetUserState(userID string) (*UserState, error) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	state := &UserState{Id: userID}
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

// BroadcastMessage 实现 IUserManager 接口
func (um *UserManager) BroadcastMessage(messageType int, data []byte) error {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	for _, up := range um.users {
		up.mutex.RLock()
		for _, conn := range up.Conns {
			if conn.State() == types.Connected {
				_ = conn.Send(types.Message{Type: messageType, Payload: data})
			}
		}
		up.mutex.RUnlock()
	}
	return nil
}

// SendMessage 实现 IUserManager 接口
func (um *UserManager) SendMessage(userID string, messageType int, data []byte) error {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	up, exists := um.users[userID]
	if !exists {
		return nil
	}

	up.mutex.RLock()
	defer up.mutex.RUnlock()
	for _, conn := range up.Conns {
		if conn.State() == types.Connected {
			_ = conn.Send(types.Message{Type: messageType, Payload: data})
		}
	}
	return nil
}

// SendPlatformMessage 实现 IUserManager 接口
func (um *UserManager) SendPlatformMessage(userID string, platformID int32, messageType int, data []byte) error {
	conn, err := um.GetUserConn(userID, platformID)
	if err != nil || conn == nil {
		return err
	}
	return conn.Send(types.Message{Type: messageType, Payload: data})
}

// GetAllUsers 实现 IUserManager 接口
func (um *UserManager) GetAllUsers() ([]*UserState, error) {
	um.mutex.RLock()
	defer um.mutex.RUnlock()

	userStates := make([]*UserState, 0, len(um.users))
	for userID := range um.users {
		state, err := um.GetUserState(userID)
		if err != nil {
			continue
		}
		userStates = append(userStates, state)
	}
	return userStates, nil
}
