package models

import (
	"gorm.io/gorm"
)

// friend godoc
// 好友模块表设计如下
// ## 好友基础表
// - 好友应该是双向的
// - 好友存在权限
// 		- 允许与好友聊天，这是当然，不允许聊天就是拉黑
// 		- 查看个人资料的权限
//		- 查看朋友圈/动态的权限
// - 给好友的备注
// - 是否特别关心
// - 好友分组
// todo 群聊是否要单独另建表？

type FriendBasic struct {
	gorm.Model

	// 用户Id
	UserId uint

	// FriendId 朋友的用户ID
	FriendId uint

	// Permissions = 0 允许聊天，但不允许查看个人资料
	// Permission = 1 允许聊天，查看个人资料，不允许看朋友圈/动态
	// Permission = 2 全允许
	Permissions byte

	// Alias 表示对该朋友的备注
	Alias string

	// IsFocusOn 表示是否关注该用户
	IsFocusOn bool

	// Group 表示了好友所属于的好友组
	// 在QQ中我们有朋友分组，而微信中有朋友圈tag，都属于分组
	Group string
}

func (f *FriendBasic) FriTableName() string {
	return "friends"
}
