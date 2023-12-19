package message

import "encoding/json"

type Request interface {
	GetData() interface{}

	Type() int

	Marshal() ([]byte, error)

	UnMarshal([]byte) error
}

const (
	ReqAuthenticate = iota
)

type CommonRequest struct {
	type_ int
	data_ *map[string]string
}

func (c *CommonRequest) GetData() interface{} {
	//TODO implement me
	return c.data_
}

func (c *CommonRequest) Type() int {
	//TODO implement me
	return c.type_
}

func (c *CommonRequest) Marshal() ([]byte, error) {
	//TODO implement me
	temp := struct {
		Type int
		Data *map[string]string
	}{
		Type: c.type_,
		Data: c.data_,
	}
	return json.Marshal(temp)
}

func (c *CommonRequest) UnMarshal(bytes []byte) error {
	//TODO implement me
	temp := struct {
		Type int
		Data *map[string]string
	}{}
	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}
	c.type_ = temp.Type
	c.data_ = temp.Data
	return nil
}
