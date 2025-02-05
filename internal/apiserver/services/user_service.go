package services

import (
	"errors"
	"time"

	"github.com/woxQAQ/gim/internal/apiserver/models"
	"github.com/woxQAQ/gim/internal/apiserver/stores"
	"github.com/woxQAQ/gim/internal/apiserver/types/response"
	"github.com/woxQAQ/gim/pkg/snowflake"
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

// Register 处理用户注册的业务逻辑
func (s *UserService) Register(user *models.User) error {
	// 生成用户ID
	user.ID = snowflake.GenerateID()

	// 检查用户名是否已存在
	exists, err := s.userStore.CheckUsernameExists(user.Username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if user.Email != "" {
		exists, err = s.userStore.CheckEmailExists(user.Email)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("邮箱已被使用")
		}
	}

	// 创建用户
	return s.userStore.CreateUser(user)
}

// Login 处理用户登录的业务逻辑
func (s *UserService) Login(username, password string) (*response.UserResponse, error) {
	// 根据用户名获取用户
	user, err := s.userStore.GetUserByUsername(username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !user.ValidatePassword(password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	user.LastLogin = time.Now()
	if err := s.userStore.UpdateUser(user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// GetUserByID 获取用户信息的业务逻辑
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userStore.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
