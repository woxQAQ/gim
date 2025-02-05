package services

import (
	"github.com/woxQAQ/gim/internal/apiserver/stores"
)

// UserService 处理用户相关的业务逻辑
type UserService struct {
	userStore *stores.UserStore
}

// NewUserService 创建UserService实例
func NewUserService(userStore *stores.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

// GetUserByID 获取用户信息的业务逻辑
func (s *UserService) GetUserByID(id string) error {
	// TODO: 实现获取用户信息的业务逻辑
	return s.userStore.GetUserByID(id)
}
