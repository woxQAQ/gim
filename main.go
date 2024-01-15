package main

import (
	"github.com/woxQAQ/gim/config"
	"github.com/woxQAQ/gim/internal/Init"
	"github.com/woxQAQ/gim/internal/server"
)

// @title API文档
// @version 1.0
// @description gIM服务器
// @host localhost:8964
// @license.name MIT

// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io/resources/open-api/
func main() {
	config.InitConfig()
	Init.Initdb()
	Init.InitLogger()
	server.StartWebServer()
}

