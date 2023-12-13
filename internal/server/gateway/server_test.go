package gateway

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"testing"
)

func TestServer(t *testing.T) {
	gsServer := &Server{
		multicore: true,
		network:   "tcp",
		addr:      "127.0.0.1:8080",
		pool:      goroutine.Default(),
	}
	gnet.Run(gsServer, gsServer.addr)

}
