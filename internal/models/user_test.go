package models

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Test_models(t *testing.T) {
	dsn := "root:woaisj870621@tcp(127.0.0.1:3306)/gIM?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&UserBasic{})
	if err != nil {
		panic(err)
	}
}
