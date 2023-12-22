package message

import (
	"encoding/json"
)

const (
	ReqCommon = iota
	ReqAuthenticate
	ReqTemp //用做测试
	ReqMessageSingle
	ReqMessageGroup
)

type RequestData map[string]any

// RequestBuffer 用于缓存请求，由于请求到来的时候，服务器并不知道请求的具体类型
// 故需要一个通用的缓存类来存储请求，在进行细分
type RequestBuffer struct {
	json.Marshaler
	json.Unmarshaler
	// type_ 表示请求类型
	type_ int

	// token_ 用于鉴权，每条消息都要携带token
	token_ string

	// userId 用于标识消息发送方
	userId string

	// data_ 用于存储具体的请求数据
	data_ *RequestData
}

func NewRequest(type_ int, data_ *RequestData, token_ string, userId string) *RequestBuffer {
	return &RequestBuffer{
		type_:  type_,
		data_:  data_,
		token_: token_,
		userId: userId,
	}
}

func (c *RequestBuffer) GetToken() string {
	return c.token_
}

func (c *RequestBuffer) GetUserId() string {
	return c.userId
}

func (c *RequestBuffer) GetData() *RequestData {
	//TODO implement me
	return c.data_
}

func (c *RequestBuffer) Type() int {
	//TODO implement me
	return c.type_
}

func (c *RequestBuffer) MarshalJSON() ([]byte, error) {
	//TODO implement me
	temp := struct {
		Type   int
		Data   *RequestData
		Token  string
		UserId string
	}{
		Type:   c.type_,
		Data:   c.data_,
		Token:  c.token_,
		UserId: c.userId,
	}
	return json.Marshal(temp)
}

func (c *RequestBuffer) UnMarshalJSON(bytes []byte) error {
	//TODO implement me
	temp := struct {
		Type   int
		Token  string
		Data   *RequestData
		UserId string
	}{}
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	c.type_ = temp.Type
	c.data_ = temp.Data
	c.userId = temp.UserId
	c.token_ = temp.Token
	return nil
}
