package user

import "github.com/woxQAQ/gim/internal/wsgateway/base"

// AddConn 实现 IUserManager 接口.
func (um *Manager) AddConn(userID string, platformID int32, conn base.LongConn) error {
	// 在锁内收集需要的状态信息
	var oldState base.ConnectionState
	var newState base.ConnectionState

	{
		um.mutex.Lock()
		up, exists := um.users[userID]
		if !exists {
			up = NewUserPlatform(userID)
			um.users[userID] = up
		}

		up.mutex.Lock()
		// 获取旧连接的状态（如果存在）
		oldState = base.Disconnected
		if oldConn, exists := up.Conns[platformID]; exists {
			oldState = oldConn.State()
		}

		up.Conns[platformID] = conn
		newState = conn.State()

		up.mutex.Unlock()
		um.mutex.Unlock()
	}

	// 锁释放后进行状态通知
	um.notifyStateChange(userID, platformID, oldState, newState)

	// 设置连接断开回调
	conn.OnDisconnect(func(err error) {
		um.notifyStateChange(userID, platformID, base.Connected, base.Disconnected)
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
func (um *Manager) GetConn(userID string, platformID int32) (base.LongConn, error) {
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
		if conn.State() == base.Connected {
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
