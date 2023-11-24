package router

import (
	"gIM/internal/server/auth"
	"github.com/gin-gonic/gin"
)

func RegisterGin(router *gin.Engine) {
	Auth := router.Group("/auth")
	{
		Auth.POST("/login", auth.Login)
		Auth.POST("/signup", auth.Signup)
	}
}
