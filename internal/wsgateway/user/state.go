package user

import (
	"time"

	"github.com/woxQAQ/gim/internal/wsgateway/types"
)

// State 定义用户在各平台的在线状态.
type State struct {
	Id              string
	OnlinePlatform  []int32
	OfflinePlatform []int32
}

// StateObserver 定义用户状态观察者接口.
type StateObserver interface {
	// OnUserStateChange 当用户状态发生变化时调用.
	OnUserStateChange(userID string, platformID int32, oldState, newState types.ConnectionState, timestamp time.Time)
}
