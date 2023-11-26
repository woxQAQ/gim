package router

import (
	"gIM/internal/middleware"
	"gIM/internal/server/users"
	"github.com/gin-gonic/gin"
)

func RegisterGin(router *gin.Engine) {
	User := router.Group("/user")
	{
		User.POST("/login", middleware.JWY(), users.Login)
		User.POST("/signup", middleware.JWY(), users.Signup)
		User.GET("/info", middleware.JWY(), users.InfoUser)
		User.POST("/update", middleware.JWY(), users.UpdateUser)
	}
}
