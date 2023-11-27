package bootstrap

import (
	"fmt"
	"gIM/internal/global"
	"gIM/internal/server/router"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	router.RegisterGin(r)

	r.Run(fmt.Sprintf(":%d", global.HttpPort))
}
