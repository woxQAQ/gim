package transfer

import (
	"sync"

	"github.com/panjf2000/gnet/v2"
)

type tsClient struct {
	kafkaWriters *sync.Map
	messagePool  *sync.Pool
	gnet.EventHandler
}

func (ts *tsClient) OnTraffic(c gnet.Conn) (action gnet.Action) {

	return
}
