package config

import (
	"time"

	"github.com/go-fuego/fuego"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"github.com/woxQAQ/gim/pkg/constants"
	"github.com/woxQAQ/gim/pkg/logger"
	"github.com/woxQAQ/gim/pkg/middleware"
)

func SetupApiServer(l logger.Logger) *fuego.Server {
	addr := viper.GetString(constants.ApiPath)
	jsonFilePath := viper.GetString(constants.OpenapiFilePath)
	disableSwagger := viper.GetBool(constants.EnableOpenapiSpec)
	swaggerUrl := viper.GetString(constants.OpenapiRoute)

	server := fuego.NewServer(fuego.WithAddr(addr),
		fuego.WithoutLogger(),
		fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			DisableSwagger:   disableSwagger,
			SwaggerUrl:       swaggerUrl,
			JsonFilePath:     jsonFilePath,
			PrettyFormatJson: true,
		}),
		fuego.WithoutAutoGroupTags(),
		fuego.WithAutoAuth(func(user, password string) (jwt.Claims, error) {
			return jwt.MapClaims{
				"user_id":  user,
				"password": password,
				"exp":      time.Now().Add(24 * time.Hour).Unix(),
			}, nil
		}),
	)
	// 设置全局中间件
	fuego.Use(server, middleware.Recovery(l))

	// 设置日志中间件
	fuego.Use(server, middleware.Logger(l))

	return server
}
