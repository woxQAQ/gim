package db

import (
	"gIM/internal/global"
	"gIM/internal/models"
	"go.uber.org/zap"
)

// TODO 需求分析：用户模块的基本crud操作
// 1. 根据用户名和密码检索用户
// 2. 添加用户
// 3. 删除用户
// 4. 更新用户信息
// 5. 根据用户名查询用户

// QueryByNameAndPwd 根据用户名和密码检索用户
func QueryByNameAndPwd(name string, password string) (*models.UserBasic, error) {
	var User models.UserBasic
	if tx := global.DB.Where(models.UserBasic{Name: name, Password: password}).First(&User); tx.Error != nil {
		return nil, tx.Error
	}
	return &User, nil
}

// CreateUser 用来添加用户
func CreateUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Create(&user)
	if tx.Error != nil {
		zap.S().Info("新建用户失败")
		return nil, tx.Error
	}
	return &user, nil
}

func DeleteUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Delete(&user)
	if tx.Error != nil {
		zap.S().Info("删除用户失败")
		return nil, tx.Error
	}
	return &user, nil
}

func UpdateUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Model(&user).Updates(models.UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Gender:   user.Gender,
	})
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &user, nil
}

func QueryByUserName(name string) (*models.UserBasic, error) {
	var User models.UserBasic
	if tx := global.DB.Where(models.UserBasic{Name: name}).First(&User); tx.Error != nil {
		return nil, tx.Error
	}
	return &User, nil
}
