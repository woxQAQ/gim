package message

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/woxQAQ/gim/internal/errorhandler"
)

type AuthenticateData struct {
	Token    string `json:"token"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type AuthenticateRequest struct {
	type_ int
	data_ *AuthenticateData
}

func NewAuthReq(type_ int, data_ *AuthenticateData) *AuthenticateRequest {
	return &AuthenticateRequest{
		type_: type_,
		data_: data_,
	}
}

func (a *AuthenticateRequest) GetData() interface{} {
	return a.data_
}

func (a *AuthenticateRequest) Type() int {
	return a.type_
}

func (a *AuthenticateRequest) Marshal() ([]byte, error) {
	temp := struct {
		Type int
		Data *AuthenticateData
	}{
		Type: a.type_,
		Data: a.data_,
	}
	return json.Marshal(temp)
}

func (a *AuthenticateRequest) UnMarshal(bytes []byte) error {
	temp := struct {
		Type int
		Data *AuthenticateData
	}{}
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	a.type_ = temp.Type
	a.data_ = temp.Data
	return nil
}

func (c *CommonRequest) ToAuthenticateRequest() (*AuthenticateRequest, error) {
	// todo 判断字段满足AuthenticateRequest
	var data *AuthenticateData
	err := mapstructure.Decode(c.data_, &data)
	if err != nil {
		return nil, errorhandler.ErrCommonTransferTo
	}
	return NewAuthReq(c.type_, data), nil
}
