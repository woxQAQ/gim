package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/woxQAQ/gim/config"
	"github.com/woxQAQ/gim/internal/server"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

var (
	once               sync.Once
	bufferPoolInstance *sync.Pool
	websocketUpgrade   *websocket.Upgrader
	goroutinePool      *goroutine.Pool
)

type GwConfig struct {
	TcpAddress       string `yaml:"gateway_tcp_address"`
	WebsocketAddress string `yaml:"gateway_websocket_address"`
	AuthAddress      string `yaml:"auth_address"`
	AuthURL          string `yaml:"auth_url"`
	TransferAddress  string `yaml:"transfer_address"`
}

type GwServer struct {
	*server.Server
	clientMap         *clientMap
	connToTransferMap *connMap
	WsMgr             *WebsocketConnManager
	WsEngine          *gin.Engine
	*GwConfig
}

func init() {
	once.Do(func() {
		bufferPoolInstance = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 1024)
			},
		}
		websocketUpgrade = &websocket.Upgrader{
			WriteBufferSize: 1024,
			ReadBufferSize:  1024,
			WriteBufferPool: &sync.Pool{},
		}
		goroutinePool, _ = ants.NewPool(
			1024,
			ants.WithPreAlloc(true),
			ants.WithNonblocking(true))
	})
}

func NewGatewayServer(network string, multicore bool, wsMgr *WebsocketConnManager) *GwServer {
	gwconfig := GwConfig{}
	buf := bufferPoolInstance.Get().([]byte)
	// 清空 buf
	buf = buf[:0]
	buf, err := os.ReadFile(config.GatewayConfigPath)
	if err != nil {
		logging.Errorf("NewGatewayServer Error: os.ReadFile Error: %v\n", err.Error())
		panic(err)
	}
	err = yaml.Unmarshal(buf, gwconfig)
	if err != nil {
		logging.Errorf("NewGatewayServer Error: yaml.Unmarshal Error: %v\n", err.Error())
		panic(err)
	}
	bufferPoolInstance.Put(buf)

	r := gin.Default()
	r.POST("/ws", func(c *gin.Context) {
		upgradeWebsocket(wsMgr, c)
	})

	return &GwServer{
		server.NewServer(network, multicore),
		clientMapInstance,
		connMapInstance,
		wsMgr,
		r,
		&gwconfig,
	}
}
