package bootstrap

import (
	"gIM/internal/server/router"
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	router.RegisterGin(r)

	_ = r.Run("127.0.0.1:8080")
}
