package gateway

import (
	"bytes"
	"net/http"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"

	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/woxQAQ/gim/internal/protobuf/proto_pb"
	"go.uber.org/zap"
)

const (
	maxMessageSize = 1024

	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10

	defaultgoroutine = 1024
)

type wsSession struct {
	// mgr is the session's connmaneger
	mgr *WebsocketConnManager

	// conn websocket conn
	conn *websocket.Conn

	// messageChan is uesd for receive message that sent to client
	messageChan chan []byte
}

type WebsocketConnManager struct {
	// Map is used to save websocket conn
	Map websocketSessionMap

	// receivedChan receive message from client
	receivedChan chan []byte

	// registerChan is to receive session that
	// want to register to the Map
	// using when a conn arrived
	registerChan chan *wsSession

	// registerChan is to receive session that
	// want to unregister to the Map
	// using when a conn is going to leave
	unregisterChan chan *wsSession

	// goroutinePool is used to reuse goroutine
	goroutinePool *goroutine.Pool
}

// NewWsMgr is used to create a new wsMgr
func NewWsMgr() *WebsocketConnManager {
	pool, err := ants.NewPool(defaultgoroutine, ants.WithPreAlloc(true), ants.WithNonblocking(true))
	if err != nil {
		panic(err)
	}
	return &WebsocketConnManager{
		Map: websocketSessionMap{
			sync.Map{},
		},
		receivedChan:   make(chan []byte, 1024),
		registerChan:   make(chan *wsSession),
		unregisterChan: make(chan *wsSession),
		goroutinePool:  pool,
	}
}

// Run is used to start a WebsocketConnManager
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
		case message := <-wm.receivedChan:
			wm.Map.Range(func(connId, Session any) bool {
				client := Session.(*wsSession)
				select {
				case client.messageChan <- message:
				default:
					close(client.messageChan)
					wm.Map.Delete(connId)
				}
				return true
			})
		}
	}
}

// upgradeWebsocket is used to upgrade a http conn to a websocket conn
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
	messageChannel := make(chan []byte, websocketUpgrade.WriteBufferSize)
	session := wsSession{
		wsMgr,
		conn,
		messageChannel,
	}
	// register session
	wsMgr.registerChan <- &session

	// to avoid to many goroutine to be created,
	// I use goroutinePool to schedule goroutines
	wsMgr.goroutinePool.Submit(session.writeToConn)
	wsMgr.goroutinePool.Submit(session.readFromChan)
}

// readFromChan used for reading message from client
func (ws *wsSession) readFromChan() {
	defer func() {
		// send ws to WebsocketConnManager for unregister
		ws.mgr.unregisterChan <- ws
		ws.conn.Close()
	}()

	ws.conn.SetReadLimit(maxMessageSize)

	// loop for read message, and send to WebsocketConnManager
	for {
		// read from websocket conn
		_, message, err := ws.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Infof("error: %v", err.Error())
			}
			break
		}

		// remove '\n' to be space, and remove the message's edge's space
		// is used to process multiline text to be one line
		message = bytes.TrimSpace(bytes.Replace(message, []byte("\n"), []byte(" "), -1))
		// send message to the WebsocketConnManager
		ws.mgr.receivedChan <- message
	}
}

func (ws *wsSession) writeToConn() {
	ticker := time.NewTimer(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.conn.Close()
	}()
	for {
		select {
		case message, ok := <-ws.messageChan:
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				ws.conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}

			// get next writer
			w, err := ws.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// write message to the writer we get
			w.Write(message)

			n := len(ws.messageChan)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-ws.messageChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
