package transfer

import "sync"

type connMap struct {
	sync.Map
}

var connMapInstance *connMap

var once sync.Once

func init() {
	once.Do(func() {
		connMapInstance = &connMap{}
	})
}
