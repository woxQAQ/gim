package Init

import (
	"fmt"
	"gIM/internal/global"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func Initdb() {
	User := "root"
	Password := "woaisj870621"
	Host := "127.0.0.1"
	Port := "3306"
	dbName := "gIM"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", User, Password, Host, Port, dbName)

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
}
