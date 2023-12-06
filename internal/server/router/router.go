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
		User := v1.Group("/user")
		{
			User.POST("/login", users.Login)
			User.POST("/signup", users.Signup)
			User.GET("/:name", users.InfoUser)
			User.POST("/update", jwt.JWY(), users.UpdateUser)
		}
	}

}
