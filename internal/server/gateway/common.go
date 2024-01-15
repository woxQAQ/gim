package gateway

import (
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/woxQAQ/gim/config"
	"github.com/woxQAQ/gim/internal/server"
	"gopkg.in/yaml.v3"
)

var (
	// once is used for singleton
	once sync.Once

	// bufferPoolInstance is a singleton of bufferpool
	bufferPoolInstance *sync.Pool

	// websocketUpgrade is used for upgrade a http conn to websocket
	websocketUpgrade *websocket.Upgrader

	// goroutinePool is a singleton of goroutine.Pool
	goroutinePool *goroutine.Pool
)

// GwConfig is the config of a gateway server
type GwConfig struct {
	// TcpAddress is the gateway's Tcp server address
	TcpAddress string `yaml:"gateway_tcp_address"`

	// WebsocketAddress is the gateway's websocket server address
	WebsocketAddress string `yaml:"gateway_websocket_address"`

	// AuthAddress is the auth server's server address
	AuthAddress string `yaml:"auth_address"`

	// AuthURL is the url to valid token
	AuthURL string `yaml:"auth_url"`

	// TransferAddress is the Transfer's Address
	TransferAddress string `yaml:"transfer_address"`
}

// GwServer is a GatewayServer struct
type GwServer struct {
	// embed basic server struct
	*server.Server

	// clientMap is used to save the client
	// connecting to this server
	clientMap *clientMap

	// connToTransferMap is used to save the
	// connnection with transfer server
	connToTransferMap *connMap

	// WsMgr is the websocket connection manager
	// receive register and unregister conn, receive
	// message and send message to conn
	WsMgr *WebsocketConnManager

	// WsEngine is a gin Engine of a gin-based
	// websocket server
	WsEngine *gin.Engine
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

	// get temp buf from bufferpool
	buf := bufferPoolInstance.Get().([]byte)
	// clear buf
	buf = buf[:0]

	// read config
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

	// creat http using gin,used for upgrading
	// websocket connections
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
