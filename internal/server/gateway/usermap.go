package gateway

import (
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type ClientMap struct {
	sync.Map
}

var MapInstance *ClientMap

var once sync.Once

func init() {
	once.Do(func() {
		MapInstance = &ClientMap{}
	})
}

func GetConnId(c gnet.Conn) string {
	return c.RemoteAddr().String()
}

func (m *ClientMap) Set(key string, value *Client) {
	if c, ok := m.Load(key); ok {
		clients := c.([]*Client)
		clients = append(clients, value)
		m.Store(key, clients)
	}
	var clients []*Client
	clients = append(clients, value)
	m.Store(key, clients)
}
