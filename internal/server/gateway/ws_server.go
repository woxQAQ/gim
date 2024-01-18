package gateway

import (
	"context"
	"fmt"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/message"
	"sync"
	"sync/atomic"
	"time"
)

func (s *GwServer) connToTransfer() {
	// todo 连接多个转发层,读取配置文件
	// 此处建立的连接是 网关客户端与转发层服务器的连接,注意区分网关服务器与客户端
	tsConn, err := s.Client.Dial("tcp", s.TransferAddress)
	logging.Infof("dialing to ts server...\n")
	if err != nil {
		logging.Fatalf("failed to dial: %v", err)
	}
	logging.Infof("ts Server connected\n")
	// 将与转发层的连接保存起来
	s.connToTransferMap.Set(GetConnId(tsConn), &tsConn)
}

func (s *GwServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running gateway servers on %s with multi-core=%t\n",
		fmt.Sprintf("%s://%s", s.Network, s.TcpAddress), s.Multicore)
	// 创建网关层客户端
	// todo 网关客户端编程
	gsClient, err := gnet.NewClient(&gwClient{
		responsePool: &sync.Pool{
			New: func() interface{} {
				return new(message.Response)
			},
		},
	})
	if err != nil {
		panic(err)
	}
	err = gsClient.Start()
	if err != nil {
		panic(err)
	}
	// 需要建立与转发层的连接，

	// todo 转发层集群，需要多条连接
	s.Client = gsClient
	s.Eng = eng
	s.connToTransfer()
	return
}

func (s *GwServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("new connection: %s\n", c.RemoteAddr().String())
	atomic.AddInt32(&s.Connected, 1)
	// todo 设置websocket上下文
	// c.SetContext(new(&wsCodec))
	c.SetContext(new(wsCodec))
	out = []byte("connection establishing...,\n")
	action = gnet.None
	return
}

func (s *GwServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ws := c.Context().(*wsCodec)
	logging.Infof("message arrived from client %s\n", c.RemoteAddr().String())

	if ws.readBuffBytes(c) == gnet.Close {
		return gnet.Close
	}

	ok, action := ws.upgrade(c)
	if !ok {
		return
	}

	if ws.buf.Len() <= 0 {
		return
	}
	messages, err := ws.Decode(c)
	if err != nil || messages == nil {
		return
	}

	for _, message := range messages {
		msgLen := len(message.Payload)
		if msgLen > 128 {
			logging.Infof("conn[%v] receive [op=%v] [msg=%v..., len=%d]",
				c.RemoteAddr().String(), message.OpCode, string(message.Payload[:128]), len(message.Payload))
		} else {
			logging.Infof("conn[%v] receive [op=%v] [msg=%v, len=%d]",
				c.RemoteAddr().String(), message.OpCode, string(message.Payload), len(message.Payload))
		}
		err := wsutil.WriteServerMessage(c, message.OpCode, message.Payload)
		if err != nil {
			logging.Infof("conn[%v] [err=%v]", c.RemoteAddr().String(), err.Error())
			return gnet.Close
		}
	}
	// 0. 获取缓冲区大小
	buf := bufferPoolInstance.Get().([]byte)
	buf = buf[:0]

	// 1. 从连接获取序列化的内容
	n, err := c.Read(buf)
	if err != nil {
		return gnet.Close
	}

	//logging.Infof("message arrived: %s\n", string(buf))
	// 2. 反序列化客户端请求
	req := s.RequestPool.Get().(*message.RequestBuffer)
	if err = req.UnMarshalJSON(buf[:n]); err != nil {
		logging.Infof("[ERROR] error: %v\n", err)
		_, err = c.Write([]byte("unmarshal error\n"))
		if err != nil {
			logging.Infof("[ERROR] gateway Server %s write error: %v\n", err)
			return gnet.Close
		}
		return gnet.Close
	}

	// 3. 鉴权
	value := s.clientMap.Get(&c)
	// 如果 value不存在，或者 token 过期，则需要重新鉴权
	if value == nil || value.expiredTime.Unix() < time.Now().Unix() {
		// todo 鉴权进阶：心跳包，
		// 如果仅仅因为 token 过期就要重新鉴权, 则未免过于繁琐
		expiredAt, err := authHandler(req)
		if err != nil {
			logging.Infof("[ERROR] error: %v\n", err)
			_, err = c.Write([]byte("auth error\n"))
			if err != nil {
				logging.Infof("[ERROR] gateway Server %s write error: %v\n", err)
				return gnet.Close
			}
			return gnet.Close
		}

		// 填写结构体
		uc := userSession{
			conn:        c,
			userId:      req.GetUserId(),
			token:       req.GetToken(),
			expiredTime: expiredAt,
		}

		s.clientMap.Set(GetConnId(c), &uc)
	}

	// 4. 处理客户端请求,直接发给转发层即可
	if req.Type() == message.ReqTestGatewayConn {
		// 用作测试用例,不发往转发层
		logging.Infof("Test is ok")
		c.Write([]byte("Test is ok"))
		return gnet.None
	}

	// 随机获取一个连接
	tsConn, err := s.connToTransferMap.GetRandomConn()
	if err != nil {
		logging.Infof("[ERROR] %v", err)
		return gnet.None
	}

	_, err = tsConn.Write(buf)
	if err != nil {
		logging.Infof("[ERROR] %v", err)
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
	logging.Infof("connection %s closed\n", c.RemoteAddr().String())
	// todo how to delete?
	ctx := context.Background()
	s.RedisConn.Del(ctx, "UserSession")
	return
}
