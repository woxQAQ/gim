package transfer

import (
	"github.com/panjf2000/gnet/v2"
	"sync"
)

type connMap struct {
	sync.Map
}

var connMapInstance *connMap

var once sync.Once

func init() {
	once.Do(func() {
		connMapInstance = &connMap{}
	})
}

func getConnId(conn gnet.Conn) string {
	return conn.RemoteAddr().String()
}

func (m *connMap) Set(key string, value *gnet.Conn) {
	// todo 校验key
	if _, ok := m.Load(key); ok {
		return
	}
	m.Store(key, value)
}
