package main

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/gateway"
)

func main() {
	gsServer := gateway.NewGatewayServer("tcp", "127.0.0.1:8088", true)
	err := gnet.Run(gsServer, gsServer.Network+"://"+gsServer.Addr)
	if err != nil {
		logging.Errorf("gatewayServer start error: %v", err)
		return
	}
}
