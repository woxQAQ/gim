package Init

import (
	"log"

	"go.uber.org/zap"
)

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalln("Zap logger init failed", err.Error())
	}
	zap.ReplaceGlobals(logger)
}
