package models

import (
	"gorm.io/gorm"
)

// friend godoc
// 好友模块表设计如下
// ## 好友基础表
// - 好友应该是双向的
// - 好友存在权限
//   - 允许与好友聊天，这是当然，不允许聊天就是拉黑
//   - 查看个人资料的权限
//   - 查看朋友圈/动态的权限
//
// - 给好友的备注
//   - 好友备注
//
// - 是否特别关心
// - 好友分组
// - 好友状态
// todo 群聊是否要单独另建表？
type Status int

const (
	Accepted Status = iota
	Pending
	Rejected
	Blocked
	Focused
)

type Relation struct {
	gorm.Model

	// UserId 表示这个关系所属于的用户。
	// 注意关系是双份的，也就是说，必然存在一份UserId和FriendId与当前信息相反的数据
	UserId uint `gorm:"index"`
	// FriendId 朋友的用户ID
	FriendId uint `gorm:"index"`

	// status 表示当前关系的状态，存在以下几种状态
	// 1. accepted 已接受，是好友状态
	// 2. pending 待处理
	// 3. rejected 拒绝
	// 值得注意的是，我们会存储双份的好友信息，1，2，3，这三种状态一定是对应的
	// 4. blocking 已拉黑
	// 5. focused 特别关注
	Status Status

	// 表示用户 UserId 对 用户 FriendId 的备注
	Alias string
}
