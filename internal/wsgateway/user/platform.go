package user

import (
	"sync"

	"github.com/woxQAQ/gim/internal/wsgateway/types"
)

// UserPlatform 管理单个用户的多平台连接.
type Platform struct {
	UserID string                   // 用户唯一标识
	Conns  map[int32]types.LongConn // 平台ID到连接的映射
	mutex  sync.RWMutex             // 用于并发安全的读写锁
}

// NewUserPlatform 创建新的用户平台管理实例.
func NewUserPlatform(userID string) *Platform {
	return &Platform{
		UserID: userID,
		Conns:  make(map[int32]types.LongConn),
	}
}
