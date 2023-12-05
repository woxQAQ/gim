package Init

import (
	"fmt"
	"gIM/internal/global"
	"gIM/internal/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initdb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.DBUser, global.DBPwd, global.DBHost, global.DBName)

	Logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: Logger,
	})
	if err != nil {
		panic(err)
	}
	err = global.DB.AutoMigrate(&models.UserBasic{}, &models.Relation{})
	if err != nil {
		panic(err)
	}
}
