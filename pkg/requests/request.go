package requests

type request interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
