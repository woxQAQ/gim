package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/woxQAQ/gim/internal/protobuf/proto_pb"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type WebsocketConnManager struct {
	// Map 用于保存注册到管理器的websocket连接
	Map websocketSessionMap

	// 来自客户端的消息
	messageChan chan []byte

	registerChan chan *wsSession

	unregisterChan chan *wsSession
}

func NewWsMgr() *WebsocketConnManager {
	return &WebsocketConnManager{
		Map: websocketSessionMap{
			sync.Map{},
		},
		messageChan: make(chan []byte, 1024),
	}
}

func (wm *WebsocketConnManager) Run() {
	for {
		select {
		case client := <-wm.registerChan:
			wm.Map.set(client.conn.RemoteAddr().String(), *client)
		case client := <-wm.unregisterChan:
			if _, ok := wm.Map.Load(client.conn.RemoteAddr().String()); ok {
				wm.Map.Delete(client.conn.RemoteAddr().String())
				close(client.messageChan)
			}
		}
	}
}

func upgradeWebsocket(wsMgr *WebsocketConnManager, c *gin.Context) {

	conn, err := websocketUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zap.S().Infof("websocket upgrade error: %v", err)
		c.ProtoBuf(http.StatusBadRequest, &proto_pb.WsResponse{
			Code:    -1,
			Message: "websocket Upgrade Error!",
			Data: &proto_pb.WsResponse_UpgradeError{
				UpgradeError: &proto_pb.UpgradeError{
					Error: err.Error(),
				},
			},
		})
	}
	requestChannel := make(chan []byte, websocketUpgrade.WriteBufferSize)
	session := wsSession{
		wsMgr,
		conn,
		requestChannel,
	}
	wsMgr.registerChan <- &session
}
