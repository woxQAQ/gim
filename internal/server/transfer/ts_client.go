package transfer

import (
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type tsClient struct {
	messagePoll *sync.Pool
	gnet.EventHandler
}

func (ts *tsClient) OnTraffic(c gnet.Conn) (action gnet.Action) {

	return
}
