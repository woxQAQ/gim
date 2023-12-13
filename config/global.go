package config

import (
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

var Cfg *ini.File

var RunMode string

var (
	HttpPort     int
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	JwtSecret    string
	RecallTimes  int
)

var (
	DBName string
	DBUser string
	DBPwd  string
	DBHost string
)

var DateTemp = "2006-01-01"

type UidT uint
