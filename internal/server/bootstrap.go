package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/woxQAQ/gim/internal/global"
	"github.com/woxQAQ/gim/internal/server/router"
)

func Start() {
	r := gin.Default()

	router.RegisterGin(r)
	gin.SetMode(global.RunMode)
	r.Run(fmt.Sprintf(":%d", global.HttpPort))
}
