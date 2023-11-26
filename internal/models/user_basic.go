package models

import (
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	Password      string
	Gender        string `gorm:"column:gender;default:male;type:varchar(6) comment'性别'"`
	Phone         string `valid:"match(1^[3~9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identify      string
	Avatar        string
	ClientIp      string `valid:"ipv4"`
	ClientPort    string
	Salt          string
	LoginTime     time.Time `gorm:"column:login_time"`
	HeartBeatTime time.Time `gorm:"column:heart_beat_time"`
	LogOutTime    time.Time `gorm:"column:logout_time"`
	IsLogOut      bool
	DeviceInfo    string
}
