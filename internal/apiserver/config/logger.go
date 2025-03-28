package config

import (
	"github.com/spf13/viper"
	"github.com/woxQAQ/gim/pkg/constants"
	"github.com/woxQAQ/gim/pkg/logger"
	"go.uber.org/zap"
)

func SetupLogger() logger.Logger {
	logLevel := viper.GetString(constants.LogLevel)
	logFile := viper.GetString(constants.LogFilePath)

	l, err := logger.NewLogger(&logger.Config{
		Level:    logLevel,
		FilePath: logFile,
	}, zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return l
}
