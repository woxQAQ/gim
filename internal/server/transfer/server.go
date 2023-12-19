package transfer

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/woxQAQ/gim/internal/errorhandler"
	"github.com/woxQAQ/gim/internal/server"
	"github.com/woxQAQ/gim/internal/server/message"
	"sync/atomic"
)

type transferServer struct {
	*server.Server
	transferId string
}

func (s *transferServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running server on %s with multi-core=%t\n",
		fmt.Sprintf("%s://%s", s.Network, s.Addr), s.Multicore)
	s.Eng = eng
	s.Pool = goroutine.Default()
	return
}

func (s *transferServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("gateway %s has been connected", c.RemoteAddr().String())
	out = []byte(fmt.Sprintf("gateway %s has been connected, so it's time to transfer your messages\n", c.RemoteAddr().String()))
	return
}

func (s *transferServer) OnTraffic(c gnet.Conn) (action gnet.Action) {

	logging.Infof("message arrived from gateway %s\n", c.RemoteAddr().String())
	buf := make([]byte, 512)
	c.Read(buf)
	req := &message.CommonRequest{}
	if err := req.UnMarshal(buf); err != nil {
		logging.Infof("[ERROR] Gateway: %s unmarshal error: %v, %v\n",
			c.RemoteAddr().String(), err, req)
		c.Write([]byte("unmarshal error\n"))
		return gnet.None
	}
	err := s.OnRequest(req, c)
	if err != nil {

	}
	return
}

// OnRequest 用来处理网关层发来的请求
func (s *transferServer) OnRequest(msg interface{}, c gnet.Conn) error {
	request, ok := msg.(*message.CommonRequest)
	if !ok {
		return errorhandler.ErrMessageNotRequest
	}
	switch request.Type() {
	case message.ReqAuthenticate:
		auth, err := request.ToAuthenticateRequest()
		if err != nil {
			// todo 重试？
			return err
		}
		// 处理请求
		response, err := authHandler(auth)
		if err != nil {
			return err
		}
		err = response.MarshalAndWrite(c)
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

func (s *transferServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("connection :%s closed due to: %v\n", c.RemoteAddr().String(), err)
		return
	}
	atomic.AddInt32(&s.Connected, -1)
	logging.Infof("connection closed: %s\n", c.RemoteAddr().String())
	return
}
