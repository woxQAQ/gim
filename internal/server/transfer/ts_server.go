package transfer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/segmentio/kafka-go"
	"github.com/woxQAQ/gim/config"
	"github.com/woxQAQ/gim/internal/protobuf/proto_pb"
	"github.com/woxQAQ/gim/internal/server"
	"github.com/woxQAQ/gim/internal/server/message"
	"gopkg.in/yaml.v3"
	"os"
	"sync/atomic"
)

type TsServer struct {
	*server.Server
	transferId       string
	gatewayConnMap   *connMap
	kafkaConnections *kafka.Writer
}
type kafkaConfig struct {
	Address string   `yaml:"kafka_address"`
	Topics  []string `yaml:"topics"`
}

func connKafka() (*kafka.Writer, error) {
	// 读取kafka配置
	data, err := os.ReadFile(config.KafkaFilePath)
	if err != nil {
		return nil, err
	}
	var configs kafkaConfig
	err = yaml.Unmarshal(data, &configs)
	if err != nil {
		return nil, err
	}

	// 连接kafka
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(configs.Address),
		Balancer: &kafka.LeastBytes{},
	}

	return kafkaWriter, nil
}

func transferId(addr string) string {
	return addr
}

func NewTransferServer(network string, addr string, multicore bool) (*TsServer, error) {
	kafkaConn, err := connKafka()
	if err != nil {
		return nil, err
	}
	return &TsServer{
		server.NewServer(network, addr, multicore),
		transferId(addr),
		connMapInstance,
		kafkaConn,
	}, nil
}

func (s *TsServer) OnBoot(eng gnet.Engine) (action gnet.Action) {
	logging.Infof("running server on %s with multi-core=%t\n",
		fmt.Sprintf("%s://%s", s.Network, s.Addr), s.Multicore)
	s.Eng = eng
	s.Pool = goroutine.Default()
	client, err := gnet.NewClient(&tsClient{
		messagePoll: bufferPoolInstance,
	})
	if err != nil {
		panic(err)
	}
	err = client.Start()
	if err != nil {
		panic(err)
	}
	// todo 连接 kafka
	s.Client = client
	return
}

func (s *TsServer) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	logging.Infof("gateway %s has been connected", c.RemoteAddr().String())
	out = []byte(fmt.Sprintf("gateway %s has been connected, "+
		"so it's time to transfer your messages\n", c.RemoteAddr().String()))
	s.gatewayConnMap.Set(getConnId(c), &c)
	return
}

func (s *TsServer) OnTraffic(c gnet.Conn) (action gnet.Action) {
	logging.Infof("message arrived from gateway %s\n", c.RemoteAddr().String())

	buf := bufferPoolInstance.Get().(*bytes.Buffer)
	_, err := c.Read(buf.Bytes())
	if err != nil {
		logging.Infof("[ERROR] transfer: %s read error: %v\n",
			c.LocalAddr().String(), err.Error())
		return gnet.Close
	}

	req := s.RequestPool.Get().(message.RequestBuffer)
	if err := req.UnMarshalJSON(buf.Bytes()); err != nil {
		logging.Infof("[ERROR] Gateway: %s unmarshal error: %v, %v\n",
			c.RemoteAddr().String(), err, req)
		return gnet.None
	}
	// 使用完毕，放回
	// 1.
	// 分发层会将请求发给kafka，再经由kafka发给业务层。
	// 对于网关层发来的消息体的处理，分发层服务器仅做一个转发功能
	// 对于业务层发来的的消息体，则是由分发层客户端进行处理
	// 分发层客户端的主要功能也是进行转发
	// 当kafka向分发层客户端发送处理完的消息，
	// 分发层客户端需要能够知道需要向哪个网关层发送消息

	// 2.
	// 很显然，分发层发送给kafka的“待处理”消息一定携带 send_id 或 receive_id，`
	// kafka一定能知道send_id 或 receive_id。
	// 理所当然，kafka发送给分发层的消息，一定也指定消息所要分发给的 user_id
	// 即，分发层一定要能知道 user_id 对应的客户端所连接的网关层服务器是哪个

	// 于是，此处要记录的是，从网关层发过来的消息是属于哪个客户端的？
	// 需要将 userSession 与 gateway connid 映射起来

	s.RedisConn.Set(context.Background(), req.GetUserId(), getConnId(c), 0)

	// lazy load
	getTopicToReqMap()

	// 其实，所有信息都应该直接全发给kafka的
	// 首先根据请求类型，找到对应的topic
	// 其次，获得消息，一个萝卜一个坑
	// 任何消息发送的时候，前端都要校验内容与类型是否是一一对应的
	// 此处理论上不应该再校验
	// todo 发送给kafka
	buf.Reset()
	reqTopic := topicToReqMap[req.Type()]
	switch req.Type() {
	case message.ReqSingleMessage:
	}
	err = s.kafkaConnections.WriteMessages(context.Background(),
		kafka.Message{
			Topic: reqTopic,
			Value:
		},
	)
	s.RequestPool.Put(req)

	return
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
