package transfer

import (
	"bytes"
	"github.com/panjf2000/gnet/v2"
	"github.com/woxQAQ/gim/internal/server/message"
	"sync"
)

type connMap struct {
	sync.Map
}

var connMapInstance *connMap
var bufferPoolInstance *sync.Pool

var once sync.Once

func init() {
	once.Do(func() {
		connMapInstance = &connMap{}
		bufferPoolInstance = &sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		}
	})
}

var topicToReqMap map[int]string

func initTopicToReqMap() {
	topicToReqMap = map[int]string{}
	topicToReqMap[message.ReqSingleMessage] = "SINGLE_MESSAGE_TRANSFER_TO_LOGIC"
	topicToReqMap[message.ReqGroupMessage] = "GROUP_MESSAGE_TRANSFER_TO_LOGIC"
}

func getTopicToReqMap() {
	once.Do(initTopicToReqMap)
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
