package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/woxQAQ/gim/config"
	"github.com/woxQAQ/gim/internal/server/router"
)

func StartWebServer() {
	r := gin.Default()
	router.RegisterGin(r)
	gin.SetMode(config.RunMode)
	r.Run(fmt.Sprintf(":%d", config.HttpPort))
}
