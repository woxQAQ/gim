package gateway

import (
	"fmt"
	"github.com/panjf2000/gnet/pkg/logging"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
	"github.com/panjf2000/gnet/v2"
	"github.com/woxQAQ/gim/internal/server/auth"
	"github.com/woxQAQ/gim/pkg/requests"
	"sync/atomic"
)

type Server struct {
	gnet.BuiltinEventEngine
	eng          gnet.Engine
	multicore    bool
	network      string
	addr         string
	connected    int32
	disconnected int32
	pool         *goroutine.Pool
}

func (s *Server) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running server on %s with multi-core=%t\n",
		fmt.Sprintf("%s://%s", s.network, s.addr), s.multicore)
	s.eng = eng
	s.pool = goroutine.Default()
	return
}

func (s *Server) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	// 要从客户端连接中获取token来鉴权
	// todo
	logging.Infof("new connection: %s\n", c.RemoteAddr().String())
	c.SetContext(new(requests.AuthenticateReq))
	atomic.AddInt32(&s.connected, 1)
	out = []byte("connection establishing...\n")
	action = gnet.None
	return
}

func (s *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	logging.Infof("message arrived: %s\n", c.RemoteAddr().String())
	buf := make([]byte, 512)
	_, err := c.Read(buf)
	if err != nil {
		c.Close()
		return
	}
	fmt.Println(buf)
	req := &requests.AuthenticateReq{}
	err = req.Unmarshal(buf)
	if err != nil {
		return gnet.Close
	}
	if !auth.IsAuthenticated(req) {
		c.Write([]byte("authentication failed\n"))
		return gnet.Close
	}
	logging.Infof("authentication success: %s\n", req.UserName)
	client := NewClient(req.UserName, req.Token, c)
	MapInstance.Set(req.UserName, client)
	c.Write([]byte("welcome\n"))
	return
}

func (s *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	s.connected--
	logging.Infof("connection closed: %s\n", c.RemoteAddr().String())
	// todo how to delete
	//MapInstance.Delete(c)
	return
}
