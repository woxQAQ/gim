package message

import (
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"testing"
)

func TestAuth(t *testing.T) {
	req := AuthenticateRequest{
		type_: ReqAuthenticate,
		data_: &AuthenticateData{
			Token:   "token",
			UserId:  "username",
			UserPwd: "password",
		},
	}
	data, _ := req.Marshal()
	request := AuthenticateRequest{}
	_ = request.UnMarshal(data)
	logging.Infof("request: %v\n", request)
}

func TestCommon(t *testing.T) {
	req := RequestBuffer{
		type_: ReqAuthenticate,
		data_: &map[string]string{
			"token":    "token",
			"username": "username",
			"password": "password",
		},
	}
	data, _ := req.Marshal()
	request := RequestBuffer{}
	_ = request.UnMarshal(data)
	logging.Infof("request: %v\n", request)
}
