package transfer

import (
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"github.com/panjf2000/gnet/v2"
	"testing"
)

func TestTransfer(t *testing.T) {
	tsServer := &TRServer{
		Multicore: true,
		Network:   "tcp",
		Addr:      "127.0.0.1:9090",
		Pool:      goroutine.Default(),
	}
	err := gnet.Run(tsServer, tsServer.Addr)
	if err != nil {
		t.Error("server start error: ", err)
		return
	}
}
