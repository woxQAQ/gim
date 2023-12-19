package gateway

import (
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/message"
)

func authHandler(s *gatewayServer, data *message.AuthenticateRequest) error {
	// 获取与转发层的连接

	toTransferConn, err := s.connMap.GetRandomConn()
	if err != nil {
		return err
	}

	// 编码数据
	jsonData, err := data.Marshal()
	if err != nil {
		return err
	}
	_, err = toTransferConn.Write(jsonData)
	if err != nil {
		return err
	}
	logging.Infof("waiting for response...\n")
	err = toTransferConn.Flush()
	if err != nil {
		return err
	}

	// 等待响应
	_, err = toTransferConn.Read(jsonData)
	if err != nil {
		return err
	}
	response := message.Response{}
	err = response.UnMarshal(jsonData)
	if err != nil {
		return err
	}

	return nil
}
