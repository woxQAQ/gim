package requests

import (
	"github.com/vmihailenco/msgpack/v5"
)

type AuthenticateReq struct {
	UserName string
	Token    string
}

func (req *AuthenticateReq) Marshal() ([]byte, error) {
	encoded, err := msgpack.Marshal(req)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func (req *AuthenticateReq) Unmarshal(encoded []byte) error {
	decoded := AuthenticateReq{}
	err := msgpack.Unmarshal(encoded, &decoded)
	if err != nil {
		return err
	}
	*req = decoded
	return nil
}
