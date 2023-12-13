package db

import (
	"errors"
	"fmt"
	"github.com/woxQAQ/gim/config"

	vad "github.com/asaskevich/govalidator"
	"github.com/woxQAQ/gim/internal/models"
)

// 查找

// 获取好友列表
// fetchFriends 查找 user 的所有好友
func fetchFriends(user models.UserBasic) ([]models.UserBasic, error) {
	// 查找 UserID 的所有好友
	FriendIds := user.Friends
	if len(FriendIds) == 0 {
		return nil, errors.New("好友列表为空")
	}
	friends := make([]models.UserBasic, 0, len(FriendIds))
	err := config.DB.Find(&friends, FriendIds).Error
	if err != nil {
		return nil, fmt.Errorf("查询好友失败: %w", err)
	}

	return friends, nil
}

func FetchFriendsByIds(friendIds []uint) ([]models.UserBasic, error) {
	friends := make([]models.UserBasic, 0, len(friendIds))
	err := config.DB.Find(&friends, friendIds).Error
	if err != nil {
		return nil, fmt.Errorf("查询好友失败: %w", err)
	}

	return friends, nil
}

// FetchFriendListByUserId 使用用户 Id 检索好友列表
func FetchFriendListByUserId(userId uint) ([]models.UserBasic, error) {
	// 查找 UserID 的所有好友
	user, err := QueryById(userId)
	if err != nil {
		return nil, err
	}

	friends, err := fetchFriends(user)
	if err != nil {
		return nil, err
	}
	return friends, err
}

// FriendListByUser 使用用户结构检索好友列表
// 登录后，我们的前端会将 user 的内容记录，而后我们根据 user 结构便可加载好友列表
func FriendListByUser(user models.UserBasic) ([]models.UserBasic, error) {
	ok, err := vad.ValidateStruct(user)
	if !ok {
		return nil, fmt.Errorf("用户结构不完整")
	}
	if err != nil {
		return nil, err
	}
	friends, err := fetchFriends(user)
	if err != nil {
		return nil, err
	}
	return friends, err
}

func FetchRequestById(userId uint) ([]models.Relation, error) {
	relations := make([]models.Relation, 0)
	err := config.DB.Where(map[string]interface{}{
		"UserId": userId,
		"Status": models.Sending,
	}).Find(&relations).Error
	//err := global.DB.Find(&relations, "UserId = ?", userId, "Status = ?", models.Pending).Error
	if err != nil {
		return nil, err
	}
	return relations, nil
}

/* 创建 */

// CreateRelation 用于向数据库中添加新的 Relation 行
// 当用户发起好友请求，就会向数据库中添加新的关系
func CreateRelation(userId uint, friendId uint) (relation models.Relation, err error) {
	relations := []models.Relation{
		{
			UserId:   userId,
			FriendId: friendId,
			Status:   models.Sending,
		},
		{
			UserId:   friendId,
			FriendId: userId,
			Status:   models.Pending,
		},
	}
	// 开启事务
	tx := config.DB.Begin()
	defer closeTransactions(tx, err)

	// 插入元素
	err = tx.Create(&relations).Error
	relation = relations[0]
	return
}

// MappingFriendToList 用于将好友映射到用户的好友集合中
// 当用户通过验证后，就会将好友映射到用户的好友集合中
func MappingFriendToList(user *models.UserBasic, friendId uint) (err error) {
	// 检查朋友是否已存在
	if _, ok := user.Friends[friendId]; ok {
		err = errors.New("用户已存在")
		return
	}
	user.Friends[friendId] = true

	tx := config.DB.Begin()
	defer closeTransactions(tx, err)

	err = tx.Model(&user).Update("Friends", user.Friends).Error
	return
}

// 修改

func UpdateFriendStatus(relation *models.Relation, status models.Status) (err error) {
	tx := config.DB.Begin()
	defer closeTransactions(tx, err)
	err = tx.Model(&relation).Update("Status", status).Error
	return
}

func UpdateFriendAlias(relation *models.Relation, status models.Status, alias string) (err error) {
	tx := config.DB.Begin()
	defer closeTransactions(tx, err)
	err = tx.Model(&relation).Update("Alias", alias).Error
	return
}

func CancelRequest(relation *models.Relation) (err error) {
	if relation.Status != models.Pending {
		err = errors.New("请求已处理")
		return
	}
	tx := config.DB.Begin()
	defer closeTransactions(tx, err)
	err = tx.Delete(&relation).Error
	return
}

func DeleteFriend(user *models.UserBasic, relation *models.Relation) (err error) {
	// 开启事务
	friendId := relation.FriendId
	friends := user.Friends
	tx := config.DB.Begin()
	defer closeTransactions(tx, err)
	if err = tx.Delete(&relation).Error; err != nil {
		return
	}
	// 删除好友
	delete(friends, friendId)
	// 更新用户
	err = tx.Model(&user).Update("Friends", friends).Error
	return
}
