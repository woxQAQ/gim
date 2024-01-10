package main

import (
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/gateway"
)

func main() {
	gsServer := gateway.NewGatewayServer("tcp", true)
	go func() {
		r := gin.Default()
		r.POST("/ws", func(c *gin.Context) {

		})
		err := r.Run(gsServer.WebsocketAddress)
		if err != nil {
			logging.Errorf("gatewayServer start error: %v", err)
			return
		}
	}()
	err := gnet.Run(gsServer, gsServer.Network+"://"+gsServer.TransferAddress)
	if err != nil {
		logging.Errorf("gatewayServer start error: %v", err)
		return
	}
}
