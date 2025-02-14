package config

import (
	"github.com/spf13/viper"

	"github.com/woxQAQ/gim/pkg/constants"
	"github.com/woxQAQ/gim/pkg/db"
	"github.com/woxQAQ/gim/pkg/logger"
)

func SetupDatabase(l logger.Logger) {
	if err := db.Init(&db.Config{
		Logger:       l,
		DatabasePath: viper.GetString(constants.DBPath),
	}); err != nil {
		l.Error(err.Error())
		panic(err)
	}
}
