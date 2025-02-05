package stores

import (
	"errors"

	"github.com/woxQAQ/gim/internal/models"

	"gorm.io/gorm"
)

// UserStore 处理用户相关的数据库操作
type UserStore struct {
	db *gorm.DB
}

// NewUserStore 创建UserStore实例
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// GetUserByID 根据用户ID获取用户信息
func (s *UserStore) GetUserByID(id string) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户信息
func (s *UserStore) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, "username = ?", username)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

// CheckUsernameExists 检查用户名是否已存在
func (s *UserStore) CheckUsernameExists(username string) (bool, error) {
	var count int64
	result := s.db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// CheckEmailExists 检查邮箱是否已存在
func (s *UserStore) CheckEmailExists(email string) (bool, error) {
	var count int64
	result := s.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// CreateUser 创建新用户
func (s *UserStore) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

// UpdateUser 更新用户信息
func (s *UserStore) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}
