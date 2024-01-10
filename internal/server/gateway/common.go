package gateway

import "sync"

var (
	once               sync.Once
	bufferPoolInstance *sync.Pool
)

func init() {
	once.Do(func() {
		bufferPoolInstance = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 1024)
			},
		}
	})
}
