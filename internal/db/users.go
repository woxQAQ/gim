package db

import (
	"errors"
	"fmt"

	"github.com/woxQAQ/gim/internal/global"
	"gorm.io/gorm"
	"github.com/woxQAQ/gim/internal/models"
)

// TODO 需求分析：用户模块的基本crud操作
// 1. 根据用户名和密码检索用户
// 2. 添加用户
// 3. 删除用户
// 4. 更新用户信息
// 5. 根据用户名查询用户

// QueryByNameAndPwd 根据用户名和密码检索用户
func QueryByNameAndPwd(name string, password string) (models.UserBasic, error) {
	var User models.UserBasic
	if tx := global.DB.Where(&models.UserBasic{Name: name, Password: password}).First(&User); tx.Error != nil {
		return User, tx.Error
	}
	return User, nil
}

// CreateUser 用来添加用户
func CreateUser(user models.UserBasic) (err error) {
	tx := global.DB.Begin()
	defer closeTransactions(tx, err)
	err = global.DB.Create(&user).Error
	if tx.Error != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}
	return nil
}

func DeleteUser(user models.UserBasic) (err error) {
	tx := global.DB.Begin()
	defer closeTransactions(tx, err)
	err = tx.Delete(&user).Error
	if err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}
	return nil
}

func UpdateUser(user models.UserBasic) (err error) {
	tx := global.DB.Begin()
	defer closeTransactions(tx, err)
	err = tx.Model(&user).Updates(&user).Error
	if err != nil {
		return fmt.Errorf("更新用户失败: %w", tx.Error)
	}
	return nil
}

func QueryByUserName(name string) (*models.UserBasic, error) {
	var User models.UserBasic
	if tx := global.DB.Where(map[string]interface{}{"Name": name}).First(&User); tx.Error != nil {
		return nil, tx.Error
	}
	return &User, nil
}

func UserExist(userId uint) (bool,error){
	var User models.UserBasic
	if err := global.DB.Where("ID = ?", userId).First(&User).Error; err != nil {
		return false, err
	}
	
	return true, nil
}

func QueryById(userId uint) (models.UserBasic, error) {
	User := models.UserBasic{}
	if tx := global.DB.Where("id = ?", userId).First(&User); tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return User, errors.New("用户不存在")
		}
		return User, tx.Error
	}
	return User, nil
}
