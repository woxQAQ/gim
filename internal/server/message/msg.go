package message

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Node struct {
	Conn     *websocket.Conn
	MsgQueen chan []byte
}

var NodeMap = make(map[uint64]*Node, 0)

var maplock sync.RWMutex

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  0,
	WriteBufferSize: 0,
	WriteBufferPool: &sync.Pool{},
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ConnEstablish(ctx *gin.Context) {
	// 获取UserId？
	/// todo 此处需要考虑修改
	userId, err := strconv.ParseUint(ctx.Query("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": "id获取失败",
		})
		log.Println("id获取失败")
		return
	}

	// 升级为websocket
	conn, err := Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("websocket连接建立失败")
		return
	}

	log.Println("webSocket连接建立")
	defer conn.Close()
	node := &Node{
		Conn:     conn,
		MsgQueen: make(chan []byte, 50),
	}

	maplock.Lock()
	NodeMap[userId] = node
	maplock.Unlock()

}

func Chat(w http.ResponseWriter, r *http.Request) {
}
