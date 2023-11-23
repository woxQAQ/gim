package router

import "github.com/gin-gonic/gin"
import "gIM/internal/server/auth"

func RegisterGin(router *gin.Engine) {
	router.POST("/login", auth.Login)
}
