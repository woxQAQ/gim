package transfer

import (
	"github.com/panjf2000/gnet/v2"
	"testing"
)

func TestTransfer(t *testing.T) {
	tsServer := newTransferServer("tcp", "127.0.0.1:8089", true)
	err := gnet.Run(tsServer, tsServer.Addr)
	if err != nil {
		t.Error("server start error: ", err)
		return
	}
}
