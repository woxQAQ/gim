package main

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/gateway"
)

func main() {
	gsServer := gateway.NewGatewayServer("tcp", true)
	gsServer.WsEngine.Run(gsServer.WebsocketAddress)
	wsMgr := gateway.NewWsMgr()
	wsMgr.Run()
	err := gnet.Run(gsServer, gsServer.Network+"://"+gsServer.TcpAddress)
	if err != nil {
		logging.Errorf("gatewayServer start error: %v", err)
		return
	}
}
