package redis

import (
	"time"
)

type Session struct {
	Id     string
	UserId string
	//Values         map[any]any
	CreateTime     time.Time
	ExpiredTime    time.Time
	LastAccessTime time.Time
	ConnID         string
	//UserAgent      string
	//Ip             string
}
