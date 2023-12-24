package server

import (
	"github.com/go-redis/redis/v8"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
	redismanager "github.com/woxQAQ/gim/internal/server/redis"
	"sync"
)

type Server struct {
	gnet.BuiltinEventEngine
	Eng         gnet.Engine
	Multicore   bool
	Network     string
	Addr        string
	Connected   int32
	Pool        *goroutine.Pool
	RedisConn   *redis.Client
	Client      *gnet.Client
	MessagePool *sync.Pool
}

func NewServer(network string, addr string, multicore bool) *Server {
	rds := redismanager.InitRedis()
	return &Server{
		Network:   network,
		Addr:      addr,
		Multicore: multicore,
		RedisConn: rds,
		Pool:      goroutine.Default(),
		MessagePool: &sync.Pool{
			New: func() interface{} {
				return new(gnet.Message)
			}},
	}
}
