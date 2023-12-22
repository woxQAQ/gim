package gateway

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/errors"
	"github.com/woxQAQ/gim/internal/server"
	"github.com/woxQAQ/gim/internal/server/message"
	"sync/atomic"
)

type GwServer struct {
	*server.Server
	gatewayId         string
	clientMap         *clientMap
	connToTransferMap *connMap
	connToAuthMap     *connMap
}

func (s *GwServer) ConnToTransfer(gsClient *gnet.Client) {
	// todo
	tsConn, err := gsClient.Dial("tcp", "127.0.0.1:8089")
	logging.Infof("dialing to ts server...\n")
	if err != nil {
		logging.Fatalf("failed to dial: %v", err)
	}
	logging.Infof("ts Server connected\n")
	// 将与转发层的连接保存起来
	s.connToTransferMap.Set(GetConnId(tsConn), &tsConn)
}

func gatewayId(addr string) string {
	return addr
}

func NewGatewayServer(network string, addr string, multicore bool) *GwServer {
	return &GwServer{
		server.NewServer(network, addr, multicore),
		gatewayId(addr),
		clientMapInstance,
		connMapInstance,
		connAuthMapInstance,
	}
}

func (s *GwServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running gateway servers on %s with multi-core=%t\n",
		fmt.Sprintf("%s://%s", s.Network, s.Addr), s.Multicore)
	// 创建网关层客户端
	gsClient, err := gnet.NewClient(s)
	if err != nil {
		panic(err)
	}
	// 需要建立与转发层的连接，
	// todo 转发层集群，需要多条连接
	s.ConnToTransfer(gsClient)
	s.Eng = eng
	return
}

func (s *GwServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("new connection: %s\n", c.RemoteAddr().String())
	atomic.AddInt32(&s.Connected, 1)
	// todo 设置websocket上下文
	// c.SetContext(new(&wsCodec))
	out = []byte("connection establishing...,\n")
	action = gnet.None
	return
}

func (s *GwServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	logging.Infof("message arrived from client %s\n", c.RemoteAddr().String())

	// 0. 获取缓冲区内容大小
	size := c.InboundBuffered()
	buf := make([]byte, size)

	// 1. 从连接获取序列化的内容
	n, err := c.Read(buf)
	if err != nil {
		return gnet.Close
	}

	//logging.Infof("message arrived: %s\n", string(buf))
	// 2. 反序列化客户端请求
	req := &message.RequestBuffer{}
	if err = req.UnMarshalJSON(buf[:n]); err != nil {
		logging.Infof("[ERROR] error: %v\n", err)
		_, err = c.Write([]byte("unmarshal error\n"))
		if err != nil {
			logging.Infof("[ERROR] gateway Server %s write error: %v\n", err)
			return gnet.Close
		}
		err = c.Flush()
		if err != nil {
			logging.Errorf("[ERROR] gateway Server %s flush error: %v\n", err)
			return gnet.Close
		}
		return gnet.Close
	}

	// 3. 鉴权
	err = authHandler(req)
	if err != nil {
		logging.Infof("[ERROR] %v", err)
		return gnet.Close
	}
	s.clientMap.Set(req.GetUserId(), c)

	if req.Type() == message.ReqTemp {
		logging.Infof("Test is ok")
		c.Write([]byte("Test is ok"))
		return gnet.None
	}

	// 3. 处理请求
	err = s.OnRequest(req, nil)
	if err != nil {
		logging.Infof("[ERROR] %v", err)
		c.Write([]byte(err.Error()))
		return gnet.Close
	}

	return
}

func (s *GwServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
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
func (s *GwServer) OnRequest(req *message.RequestBuffer, c gnet.Conn) error {
	// 每个请求都要鉴权
	data := req.GetData()
	msg, ok := (*data)["message"].(string)
	if !ok {
		return errors.ErrTemp
	}
	logging.Infof("message: %s\n", msg)
	return nil
}
