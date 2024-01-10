package server

import (
	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	"github.com/woxQAQ/gim/internal/server/message"
	redismanager "github.com/woxQAQ/gim/internal/server/redis"
	"sync"
)

type Server struct {
	gnet.BuiltinEventEngine
	Eng         gnet.Engine
	Multicore   bool
	Network     string
	Connected   int32
	Pool        *goroutine.Pool
	RedisConn   *redis.Client
	Client      *gnet.Client
	RequestPool *sync.Pool
}

func NewServer(network string, multicore bool) *Server {
	rds := redismanager.InitRedis()
	return &Server{
		Network:   network,
		Multicore: multicore,
		RedisConn: rds,
		Pool:      goroutine.Default(),
		RequestPool: &sync.Pool{
			New: func() interface{} {
				return new(message.RequestBuffer)
			},
		},
	}
}
