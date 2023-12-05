package db

import (
	"errors"
	"fmt"
	"gIM/internal/global"
	"gIM/internal/models"

	vad "github.com/asaskevich/govalidator"
)

// 获取好友列表
// fetchFriends 查找 user 的所有好友
func fetchFriends(user models.UserBasic) ([]models.UserBasic, error) {
	// 查找 UserID 的所有好友
	FriendIds := user.Friends
	if len(FriendIds) == 0 {
		return nil, errors.New("好友列表为空")
	}
	friends := make([]models.UserBasic, 0, len(FriendIds))
	err := global.DB.Find(&friends, FriendIds).Error
	if err != nil {
		return nil, fmt.Errorf("查询好友失败: %w", err)
	}

	return friends, nil
}

// FriendListById 使用用户 Id 检索好友列表
func GetFriendListByUserId(userId uint) ([]models.UserBasic, error) {
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

func CreateFriend(user models.UserBasic, friendId uint) (err error) {
	// 验证结构完整性
	if ok, err := vad.ValidateStruct(user); !ok {
		return fmt.Errorf("用户结构不完整")
	} else if err != nil {
		return err
	}

	if user.ID == friendId {
		return errors.New("不能添加自己为好友")
	}
	if tx := global.DB.Where("user_id = ? and friend_id = ?", user.ID, friendId).
		First(&models.Relation{}); tx.RowsAffected != 0 {
		return errors.New("好友已添加过")
	}
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	relations := []models.Relation{
		{
			UserId:   user.ID,
			FriendId: friendId,
			Status:   2,
		},
		{
			UserId:   friendId,
			FriendId: user.ID,
			Status:   2,
		},
	}
	user.Friends = append(user.Friends, friendId)
	UpdateUser(user)
	if err = tx.Updates(&user).Error; err != nil {
		return fmt.Errorf("添加好友失败: %w", err)
	}
	if err = tx.Create(&relations).Error; err != nil {
		return fmt.Errorf("添加好友失败: %w", err)
	}

	tx.Commit()
	return nil
}

func DeleteFriend(user models.UserBasic, friendId uint) error {
	// 验证结构完整性
	if ok, err := vad.ValidateStruct(user); !ok {
		return fmt.Errorf("用户结构不完整")
	} else if err != nil {
		return err
	}
	return nil
}
