package gateway

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/errorhandler"
	"github.com/woxQAQ/gim/internal/server"
	"github.com/woxQAQ/gim/internal/server/message"
	"sync/atomic"
)

type gatewayServer struct {
	*server.Server
	gatewayId string
	clientMap *ClientMap
	connMap   *ConnMap
}

func gatewayId(addr string) string {
	return addr
}

func newGatewayServer(network string, addr string, multicore bool) *gatewayServer {
	return &gatewayServer{
		server.NewServer(network, addr, multicore),
		gatewayId(addr),
		ClientMapInstance,
		ConnMapInstance,
	}
}

func (s *gatewayServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running %s on %s with multi-core=%t\n",
		s.ServerName(), fmt.Sprintf("%s://%s", s.Network, s.Addr), s.Multicore)
	// 创建网关层客户端
	gsClient, err := gnet.NewClient(s)
	if err != nil {
		panic(err)
	}
	// 需要建立与转发层的连接，
	// todo 转发层集群，需要多条连接
	tsConn, err := gsClient.Dial("tcp", "127.0.0.1:9090")
	logging.Infof("dialing ts server...\n")
	if err != nil {
		logging.Fatalf("failed to dial: %v", err)
	}
	logging.Infof("ts Server connected\n")
	// 将与转发层的连接保存起来
	s.connMap.Set(GetConnId(tsConn), &tsConn)

	s.Eng = eng
	return
}

func (s *gatewayServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	// 要从客户端连接中获取token来鉴权
	// todo 鉴权
	logging.Infof("new connection: %s\n", c.RemoteAddr().String())
	c.Write([]byte("Authenticate...\n"))
	c.Context()
	atomic.AddInt32(&s.Connected, 1)
	out = []byte("connection establishing...,\n")
	action = gnet.None
	return
}

func (s *gatewayServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	logging.Infof("message arrived from client %s\n", c.RemoteAddr().String())
	buf := make([]byte, 512)

	// 1. 从连接获取序列化的内容
	n, err := c.Read(buf)
	if err != nil {
		return gnet.Close
	}

	logging.Infof("message arrived: %s\n", string(buf))

	// 2. 反序列化客户端请求
	req := &message.CommonRequest{}
	if err = req.UnMarshal(buf[:n]); err != nil {
		logging.Infof("[ERROR] Client: %s unmarshal error: %v, %v\n",
			c.RemoteAddr().String(), err, req)
		c.Write([]byte("unmarshal error\n"))
		return gnet.Close
	}
	err = s.OnRequest(req, nil)
	if err != nil {
		logging.Infof("[ERROR] %v", err)
		c.Write([]byte(err.Error()))
		return gnet.Close
	}

	return
}

func (s *gatewayServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	atomic.AddInt32(&s.Connected, -1)
	if err != nil {
		logging.Warnf("connection :%s closed due to: %v\n", c.RemoteAddr().String(), err)
		return
	}
	logging.Infof("connection closed: %s\n", c.RemoteAddr().String())
	// todo how to delete
	s.clientMap.Delete(c)
	return
}

// OnRequest 用来处理客户端发来的请求
// msg 任意消息
// todo msg类型改为any
// c 客户端连接
func (s *gatewayServer) OnRequest(msg interface{}, c gnet.Conn) error {
	request, ok := msg.(*message.CommonRequest)
	if !ok {
		return errorhandler.ErrMessageNotRequest
	}
	switch request.Type() {
	case message.ReqAuthenticate:
		auth, err := request.ToAuthenticateRequest()
		if err != nil {
			return errorhandler.ErrCommonTransferTo
		}
		return authHandler(s, auth)
	default:
		return nil
	}
}
