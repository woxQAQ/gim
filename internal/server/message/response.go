package message

import (
	"encoding/json"
	"github.com/panjf2000/gnet/v2"
)

type ResponseData map[string]any

type Response struct {
	code_ int32
	data_ *ResponseData
	err   string
}

func NewResponse(code_ int32, data_ *ResponseData, err string) *Response {
	return &Response{
		code_: code_,
		data_: data_,
		err:   err,
	}
}

func (r *Response) Marshal() ([]byte, error) {
	// todo
	temp := struct {
		Code int32
		Data *ResponseData
		Err  string
	}{
		Code: r.code_,
		Data: r.data_,
		Err:  r.err,
	}
	return json.Marshal(temp)
}

func (r *Response) UnMarshal(bytes []byte) error {
	// todo
	temp := struct {
		Code int32
		Data *ResponseData
		Err  string
	}{}

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	r.code_ = temp.Code
	r.data_ = temp.Data
	r.err = temp.Err
	return nil
}

// MarshalAndWrite 将响应体编码成json，然后经过c发送到网关层
// c 本次请求的网关层连接
func (r *Response) MarshalAndWrite(c gnet.Conn) error {
	jsonData, err := r.Marshal()
	if err != nil {
		//todo 异常处理
		return err
	}
	_, err = c.Write(jsonData)
	if err != nil {
		return err
	}
	err = c.Flush()
	if err != nil {
		return err
	}
	return nil
}
