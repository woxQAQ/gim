package transfer

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/woxQAQ/gim/internal/server"
	"github.com/woxQAQ/gim/internal/server/message"
	"sync/atomic"
)

type TsServer struct {
	*server.Server
	transferId     string
	gatewayConnMap *connMap
}

func transferId(addr string) string {
	return addr
}

func NewTransferServer(network string, addr string, multicore bool) *TsServer {
	return &TsServer{
		server.NewServer(network, addr, multicore),
		transferId(addr),
		connMapInstance,
	}
}

func (s *TsServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running server on %s with multi-core=%t\n",
		fmt.Sprintf("%s://%s", s.Network, s.Addr), s.Multicore)
	s.Eng = eng
	s.Pool = goroutine.Default()
	client, err := gnet.NewClient(s)
	if err != nil {
		panic(err)
	}
	s.Client = client
	return
}

func (s *TsServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("gateway %s has been connected", c.RemoteAddr().String())
	out = []byte(fmt.Sprintf("gateway %s has been connected, so it's time to transfer your messages\n", c.RemoteAddr().String()))
	s.gatewayConnMap.Set(getConnId(c), &c)
	return
}

func (s *TsServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	logging.Infof("message arrived from gateway %s\n", c.RemoteAddr().String())

	size := c.InboundBuffered()
	// todo 对象复用
	buf := make([]byte, size)

	_, err := c.Read(buf)
	if err != nil {
		logging.Infof("[ERROR] transfer: %s read error: %v\n", c.LocalAddr().String(), err.Error())
		return gnet.Close
	}

	req := s.messagePool.Get().(message.RequestBuffer)
	if err := req.UnMarshalJSON(buf); err != nil {
		logging.Infof("[ERROR] Gateway: %s unmarshal error: %v, %v\n",
			c.RemoteAddr().String(), err, req)
		c.Write([]byte("unmarshal error\n"))
		return gnet.None
	}
	err = s.OnRequest(&req, c)
	if err != nil {
		logging.Infof("[ERROR] Request handler error: %v", err)
		return gnet.None
	}
	s.messagePool.Put(req)
	return
}

// OnRequest 用来处理网关层发来的请求
func (s *TsServer) OnRequest(req *message.RequestBuffer, c gnet.Conn) error {
	switch req.Type() {
	case message.ReqTestTran:
		reqData := req.GetData()
	default:
		return nil
	}
}

func (s *TsServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("connection :%s closed due to: %v\n", c.RemoteAddr().String(), err)
		return
	}
	atomic.AddInt32(&s.Connected, -1)
	logging.Infof("connection closed: %s\n", c.RemoteAddr().String())
	return
}
