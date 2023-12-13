package Init

import (
	"fmt"
	"github.com/woxQAQ/gim/config"
	"log"
	"os"
	"time"

	"github.com/woxQAQ/gim/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initdb() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DBUser, config.DBPwd, config.DBHost, config.DBName)

	Logger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})
	var err error
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: Logger,
	})
	if err != nil {
		panic(err)
	}
	err = config.DB.AutoMigrate(&models.UserBasic{}, &models.Relation{})
	if err != nil {
		panic(err)
	}
}
