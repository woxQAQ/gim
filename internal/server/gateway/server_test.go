package gateway

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	gsServer := &gatewayServer{
		Multicore: true,
		Network:   "tcp",
		Addr:      "127.0.0.1:8080",
		Pool:      goroutine.Default(),
	}
	err := gnet.Run(gsServer, gsServer.Network+"://"+gsServer.Addr)
	if err != nil {
		t.Error("gatewayServer start error: ", err)
		return
	}
}

type s struct {
	gnet.BuiltinEventEngine
}

func (s *s) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running server\n")
	return
}

func (s *s) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext("aaaaaa")
	return
}

func (s *s) OnTraffic(c gnet.Conn) (action gnet.Action) {
	if ctx, ok := c.Context().(string); ok {
		logging.Infof("message arrived: %s\n", ctx)
	} else {
		logging.Infof("message arrived: %s\n", c.RemoteAddr().String())
	}
	return
}

func TestContext(t *testing.T) {
	s := &s{}
	err := gnet.Run(s, "tcp://127.0.0.1:8088")
	if err != nil {
		t.Errorf("server start error: %v\n", err)
	}
}

func TestContextClient(t *testing.T) {
	c, err := net.Dial("tcp", "127.0.0.1:8088")
	if err != nil {
		t.Errorf("client dial error: %v\n", err)
	}
	c.Write([]byte("hello"))
}
