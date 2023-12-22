package main

import (
	"github.com/woxQAQ/gim/config"
	"github.com/woxQAQ/gim/internal/server"
)

func main() {
	config.InitConfig()
	server.StartWebServer()
}
