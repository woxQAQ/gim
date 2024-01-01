package gateway

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/message"
	"sync"
)

type gwClient struct {
	responsePool *sync.Pool
	gnet.EventHandler
}

func (gc *gwClient) OnTraffic(c gnet.Conn) (action gnet.Action) {
	buf := bufferPoolInstance.Get().([]byte)
	_, err := c.Read(buf)
	if err != nil {
		logging.Errorf("gateway %s read error: %v\n", c.LocalAddr().String(), err)
		return
	}
	response := gc.responsePool.Get().(*message.Response)

	err = response.UnMarshal(buf)
	if err != nil {
		logging.Errorf("gateway %s unmarshal error: %v, %v\n",
			c.RemoteAddr().String(), err, response)
		return
	}
	gc.OnResponse(response)
	gc.responsePool.Put(buf)
	return
}

func (gc *gwClient) OnResponse(response *message.Response) {
	logging.Infof("response arrived: %v\n", response)
	return
}
