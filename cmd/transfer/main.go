package main

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/transfer"
)

func main() {
	tsServer := transfer.NewTransferServer("tcp", "127.0.0.1:8089", true)
	err := gnet.Run(tsServer, tsServer.Addr)
	if err != nil {
		logging.Errorf("server start error: %v", err)
		return
	}
}
