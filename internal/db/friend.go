package db

import (
	"errors"
	"gIM/internal/global"
	"gIM/internal/models"
	"go.uber.org/zap"
)

func FriendList(userId uint) (*[]models.UserBasic, error) {
	Friend := make([]models.FriendBasic, 0)
	// 查找 UserID 的所有
	if tx := global.DB.Where(&models.FriendBasic{UserId: userId}).First(&Friend); tx.Error != nil {
		zap.S().Info("未找到好友信息")
		return nil, errors.New("未找到好友信息")
	}

	FriId := make([]uint, 0)
	for _, t := range Friend {
		FriId = append(FriId, t.FriendId)
	}

	var user = make([]models.UserBasic, 0)
	if tx := global.DB.Where("ID in ?", FriId).Find(&user); tx.Error != nil {
		zap.S().Info("未查找到好友")
		return nil, errors.New("未查找到好友")
	}

	return &user, nil
}
