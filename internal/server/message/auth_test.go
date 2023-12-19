package message

import (
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"testing"
)

func TestAuth(t *testing.T) {
	req := AuthenticateRequest{
		type_: ReqAuthenticate,
		data_: &AuthenticateData{
			Token:    "token",
			UserName: "username",
			Password: "password",
		},
	}
	data, _ := req.Marshal()
	request := AuthenticateRequest{}
	_ = request.UnMarshal(data)
	logging.Infof("request: %v\n", request)
}

func TestCommon(t *testing.T) {
	req := CommonRequest{
		type_: ReqAuthenticate,
		data_: &map[string]string{
			"token":    "token",
			"username": "username",
			"password": "password",
		},
	}
	data, _ := req.Marshal()
	request := CommonRequest{}
	_ = request.UnMarshal(data)
	logging.Infof("request: %v\n", request)
}
