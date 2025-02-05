package stores

import (
	"errors"

	"github.com/woxQAQ/gim/internal/apiserver/models"
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
func (s *UserStore) GetUserByID(id string) error {
	var user models.User
	result := s.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return result.Error
	}
	return nil
}
