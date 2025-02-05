package stores

import (
	"github.com/woxQAQ/gim/internal/models"

	"gorm.io/gorm"
)

// MessageStore 处理消息相关的数据库操作
type MessageStore struct {
	db *gorm.DB
}

// NewMessageStore 创建MessageStore实例
func NewMessageStore(db *gorm.DB) *MessageStore {
	return &MessageStore{db: db}
}

// CreateMessage 创建新消息
func (s *MessageStore) CreateMessage(message *models.Message) error {
	return s.db.Create(message).Error
}

// GetMessagesByUserID 获取用户的消息历史记录
func (s *MessageStore) GetMessagesByUserID(userID string, limit int, lastID string) ([]*models.Message, error) {
	var messages []*models.Message
	query := s.db.Where("sender_id = ? OR receiver_id = ?", userID, userID)

	if lastID != "" {
		query = query.Where("id < ?", lastID)
	}

	err := query.Order("created_at desc").Limit(limit).Find(&messages).Error
	return messages, err
}
