package types

import "time"

type User struct {
	// 基本信息
	ID         string    `json:"id" gorm:"primaryKey"`
	PlatformId string    `json:"platform_id" gorm:"uniqueIndex"`
	Nickname   string    `json:"nickname"`
	Avatar     string    `json:"avatar"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// 认证信息
	Password string `json:"-" gorm:"not null"` // 密码哈希，json序列化时忽略

	// 状态信息
	LastOnline   time.Time `json:"last_online"`
	OnlineStatus int       `json:"online_status"` // 0:离线 1:在线

	// 元数据
	Metadata map[string]interface{} `json:"metadata" gorm:"serializer:json"`
}
