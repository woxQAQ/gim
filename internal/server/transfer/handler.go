package transfer

import (
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/woxQAQ/gim/internal/server/message"
)

func authHandler(req *message.AuthenticateRequest) (*message.Response, error) {
	// todo
	logging.Infof("authenticate success")
	data := req.GetData().(*message.AuthenticateData)
	response := message.NewResponse(0, &message.ResponseData{
		"message": "OK",
		"token":   getToken(data),
	}, "")
	return response, nil
}
