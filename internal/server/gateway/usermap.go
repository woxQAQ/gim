package gateway

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"math/rand"
	"sync"
	"time"
)

type ClientMap struct {
	sync.Map
}

type ConnMap struct {
	sync.Map
}
type transferConnIdList struct {
	mu   sync.Mutex
	list []string
}

var ClientMapInstance *ClientMap
var ConnMapInstance *ConnMap
var TransferConnIdListInstance *transferConnIdList
var once sync.Once

func init() {
	once.Do(func() {
		ClientMapInstance = &ClientMap{}
		ConnMapInstance = &ConnMap{}
		TransferConnIdListInstance = &transferConnIdList{
			list: make([]string, 0),
			mu:   sync.Mutex{},
		}
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

func (m *ClientMap) IsExist(conn *gnet.Conn) bool {
	_, ok := m.Load(GetConnId(*conn))
	return ok
}

func (m *ConnMap) Set(key string, value *gnet.Conn) {
	// todo 校验key
	if _, ok := m.Load(key); ok {
		return
	}
	m.Store(key, value)
	TransferConnIdListInstance.Add(key)
}

func (m *ConnMap) GetRandomConn() (gnet.Conn, error) {
	rand.NewSource(time.Now().UnixNano())
	randIndex := rand.Intn(len(TransferConnIdListInstance.list))
	randKey, err := TransferConnIdListInstance.Get(randIndex)
	if err != nil {
		return nil, err
	}
	conn, ok := m.Load(randKey)
	if !ok {
		return nil, fmt.Errorf("get random connection failed, maybe there are some fault when you connect to transfer server")
	}
	return conn.(gnet.Conn), nil
}

func (l *transferConnIdList) Add(id string) {
	l.mu.Lock()
	l.list = append(l.list, id)
	l.mu.Unlock()
}

func (l *transferConnIdList) Get(index int) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if index < 0 || index >= len(l.list) {
		return "", fmt.Errorf("index out of range")
	}
	return l.list[index], nil
}
