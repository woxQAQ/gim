package server

import (
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
)

type Server struct {
	gnet.BuiltinEventEngine
	Eng       gnet.Engine
	Multicore bool
	Network   string
	Addr      string
	Connected int32
	Type      int32
	Pool      *goroutine.Pool
	ReqHandler
}

type ReqHandler interface {
	// OnRequest is called when a request is received from client
	OnRequest(msg interface{}, c gnet.Conn) error
}

const (
	ServerTypeGateway = iota
	ServerTypeTransfer
)

func (s *Server) ServerName() string {
	switch s.Type {
	case ServerTypeGateway:
		return "gateway"
	case ServerTypeTransfer:
		return "transfer"
	default:
		return "unknown"
	}
}

func NewServer(network string, addr string, multicore bool) *Server {
	return &Server{
		Network:   network,
		Addr:      addr,
		Multicore: multicore,
		Pool:      goroutine.Default(),
	}
}
