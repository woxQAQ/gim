package server

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"testing"
)

type clientHandler struct {
	gnet.EventHandler
}

func (ch *clientHandler) OnBoot(g gnet.Engine) (action gnet.Action) {
	logging.Infof("client OnBoot\n")
	return gnet.None
}

func (ch *clientHandler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("client OnOpen\n")
	return
}
func (ch *clientHandler) OnTraffic(c gnet.Conn) (action gnet.Action) {
	logging.Infof("client OnTraffic\n")
	return
}

func TestNewClient(t *testing.T) {

	client, err := gnet.NewClient(&clientHandler{})
	if err != nil {
		t.Error(err)
	}
	_, err = client.Dial("tcp", "127.0.0.1:8088")
	if err != nil {
		t.Error(err)
	}
	client.Start()
	gnet.Run()
}
