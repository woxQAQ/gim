package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/woxQAQ/gim/pkg/auth"
)

// User 用户模型
type User struct {
	ID        string    `gorm:"primaryKey;type:text"`
	Username  string    `gorm:"type:text;not null;unique"`
	Password  string    `gorm:"type:text;not null"`
	Nickname  string    `gorm:"type:text;default:''"`
	Avatar    string    `gorm:"type:text;default:''"`
	Gender    int8      `gorm:"type:smallint;default:0"` // 0: 未知, 1: 男, 2: 女
	Phone     string    `gorm:"type:text;index"`
	Email     string    `gorm:"type:text;index"`
	Status    int8      `gorm:"type:smallint;default:1;index"` // 1: 正常, 0: 禁用
	Bio       string    `gorm:"type:text;default:''"`
	LastLogin time.Time `gorm:"index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}

// BeforeSave 在保存用户数据前对密码进行加密
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 只有在密码被修改时才进行加密
	if u.Password != "" {
		hashedPassword, err := auth.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return nil
}

// ValidatePassword 验证密码是否正确
func (u *User) ValidatePassword(password string) bool {
	return auth.ValidatePassword(password, u.Password)
}
