package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        string    `gorm:"primaryKey;type:text"`
	Username  string    `gorm:"type:text;not null;unique"`
	Password  string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
