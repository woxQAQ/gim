package server

import (
	"fmt"
	"gIM/internal/global"
	"gIM/internal/server/router"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	router.RegisterGin(r)
	gin.SetMode(global.RunMode)
	r.Run(fmt.Sprintf(":%d", global.HttpPort))
}
