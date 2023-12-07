package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/woxQAQ/gim/docs"
	"github.com/woxQAQ/gim/internal/middleware/jwt"
	"github.com/woxQAQ/gim/internal/server/users"
)

func RegisterGin(router *gin.Engine) {
	// 获取 swagger 文档
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := router.Group("/v1")
	{
		// user api
		auth := v1.Group("/auth")
		{
			auth.POST("/login", users.Login)
			auth.POST("/signup", users.Signup)
		}
		User := v1.Group("/users")
		{
			User.GET("/:id", users.InfoUser)
			User.POST("/:id", jwt.JWY(), users.UpdateUser)
			User.DELETE("/:id", jwt.JWY(), users.DelUser)
		}
	}

}
