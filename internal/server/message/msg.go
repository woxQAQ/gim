package message

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Node struct {
	Conn *websocket.Conn
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	WriteBufferPool: &sync.Pool{},
}
