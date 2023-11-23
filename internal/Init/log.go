package Init

import (
	"go.uber.org/zap"
	"log"
)

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Zap logger init failed", err.Error())
	}
	zap.ReplaceGlobals(logger)
}
