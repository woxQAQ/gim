package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name     string `valid:"required"`
	Password string `valid:"required"`
	Gender   string `gorm:"column:gender;default:male;type:varchar(6) comment'性别'"`
	Phone    string `valid:"match(1^[3~9]{1}\\d{9}$)"`
	Email    string `valid:"email"`
	// Bio 是个人简介
	Bio string
	// Identify 是身份证号
	Identify string
	// Avatar 是头像
	Avatar        string
	ClientIp      string `valid:"ipv4"`
	ClientPort    int
	Salt          string `valid:"required"`
	Online        bool
	LoginTime     time.Time `gorm:"column:login_time"`
	HeartBeatTime time.Time `gorm:"column:heart_beat_time"`
	LogOutTime    time.Time `gorm:"column:logout_time"`
	// Friends 是好友列表
	Friends map[uint]bool
}

func (UserBasic) TableName() string {
	return "user_basic"
}
