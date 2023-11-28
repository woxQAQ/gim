package models

import (
	"time"
)

type MessageBasic struct {
	ID       uint      `gorm:"primarykey" json:"id"`
	Type     int       `json:"msg_type"`
	Content  string    `json:"content"`
	CreateAt time.Time `json:"create_at"`
}

type MessageIndex struct {
	ID       uint `gorm:"primarykey" json:"id"`
	FromID   uint `json:"send_id"`
	TarID    uint `json:"recv_id"`
	IsSender bool
	MsgId    uint
}
