package test

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"testing"
)

type server struct {
	gnet.BuiltinEventEngine
}

func TestWsClient(t *testing.T) {
	_, _, hs, err := ws.Dial(context.Background(), "ws://127.0.0.1:8888")
	fmt.Println(hs)
	if err != nil {
		t.Error(err)
	}
}

func TestGnet(t *testing.T) {
	s := &server{}
	err := gnet.Run(s, "tcp://127.0.0.1:8888")
	if err != nil {
		logging.Error(err)
		return
	}
}

func (s *server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("connection arrived\n")
	_, err := ws.Upgrade(c)
	if err != nil {
		fmt.Println(err.Error())
		out = ([]byte)(err.Error())
	}
	return
}
