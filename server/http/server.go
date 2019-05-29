package http

import (
	"log"
	"net/http"

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
	http.Handle("/cache/", s.getCacheHandler())
	http.Handle("/status", s.getStatusHandler())
	http.Handle("/cluster", s.getClusterHandler())
	http.ListenAndServe(s.Addr()+":6800", nil)
	log.Println("http listen address", s.Addr()+":6800")
}

// New ...
func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
