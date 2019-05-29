package tcp

import (
	"net"

	"../cache"
	"../cluster"
)

// Server ...
type Server struct {
	cache.Cache
	cluster.Node
}

// Listen ...
func (s *Server) Listen() {
	l, e := net.Listen("tcp", ":6810")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.process(c)
	}
}

// New ...
func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
