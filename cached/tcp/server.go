package tcp

import (
	"net"

	"../cache"
)

// Server ...
type Server struct {
	cache.Cache
}

// Listen ...
func (s *Server) Listen() {
	l, e := net.Listen("tcp", ":9000")
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
func New(c cache.Cache) *Server {
	return &Server{c}
}
