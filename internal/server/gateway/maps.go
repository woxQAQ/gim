package gateway

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"math/rand"
	"sync"
	"time"
)

// userSession 保存用户的连接,键为userid，值为userConn
type userSession struct {
	conn        gnet.Conn
	userId      string
	token       string
	expiredTime time.Time
}

// clientMap 用于保存鉴权成功后的客户端与网关的连接
type clientMap struct {
	sync.Map
}

type connMap struct {
	sync.Map
	// s 负责存储所有的 connID，为了实现负载均衡（随机调度）
	s safeSlice
}
type safeSlice struct {
	mu   sync.RWMutex
	list []string
}

// clientMapInstance ClientMap的唯一实例，用于存储网关层与客户端的连接
// 需要实现；当转发层传递消息给网关层时，网关层需要将消息转发给客户端
var clientMapInstance *clientMap

// connMapInstance 是 ConnMap的唯一实例，用于存储网关层和转发层的连接
var connMapInstance *connMap

// 单例模式
func init() {
	once.Do(func() {
		clientMapInstance = &clientMap{}
		connMapInstance = &connMap{
			s: safeSlice{
				mu:   sync.RWMutex{},
				list: make([]string, 0),
			},
		}
	})
}

func GetConnId(c gnet.Conn) string {
	return c.RemoteAddr().String()
}

func (m *clientMap) Set(connId string, value *userSession) {
	if _, ok := m.Load(connId); ok {
		// 已经存在
		return
	}
	m.Store(connId, value)
}

func (m *clientMap) IsExist(conn *gnet.Conn) bool {
	_, ok := m.Load(GetConnId(*conn))
	return ok
}

func (m *clientMap) Get(conn *gnet.Conn) *userSession {
	if v, ok := m.Load(GetConnId(*conn)); ok {
		return v.(*userSession)
	}
	return nil
}

// connMap

func (m *connMap) Set(key string, value *gnet.Conn) {
	// todo 校验key
	if _, ok := m.Load(key); ok {
		return
	}
	m.Store(key, value)
	m.s.Add(key)
}

// GetRandomConn 随机选取一个连接，进行初步的负载均衡（随机调度）
// todo 更好的负载均衡算法
func (m *connMap) GetRandomConn() (gnet.Conn, error) {
	rand.NewSource(time.Now().UnixNano())
	randIndex := rand.Intn(len(m.s.list))
	randKey, err := m.s.Get(randIndex)
	if err != nil {
		return nil, err
	}
	conn, ok := m.Load(randKey)
	if !ok {
		return nil, fmt.Errorf("get random connection failed, maybe there are some fault ")
	}
	return conn.(gnet.Conn), nil
}

// Add 往 s 中添加一个元素,并发安全
func (l *safeSlice) Add(id string) {
	l.mu.Lock()
	l.list = append(l.list, id)
	l.mu.Unlock()
}

// Get 从 s 中获取一个元素,并发安全
func (l *safeSlice) Get(index int) (string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// 防止index溢出
	if index < 0 || index >= len(l.list) {
		return "", fmt.Errorf("index out of range")
	}

	return l.list[index], nil
}
