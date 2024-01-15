package server

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/gnet/v2"
	"github.com/woxQAQ/gim/internal/server/message"
	redismanager "github.com/woxQAQ/gim/internal/server/redis"
)

type Server struct {
	gnet.BuiltinEventEngine
	Eng         gnet.Engine
	RedisConn   *redis.Client
	Client      *gnet.Client
	RequestPool *sync.Pool
	Network     string
	Connected   int32
	Multicore   bool
}

func NewServer(network string, multicore bool) *Server {
	rds := redismanager.InitRedis()
	return &Server{
		Network:   network,
		Multicore: multicore,
		RedisConn: rds,
		RequestPool: &sync.Pool{
			New: func() interface{} {
				return new(message.RequestBuffer)
			},
		},
	}
}
