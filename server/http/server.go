package http

import (
	"net/http"

	"../cache"
)

// Server ...
type Server struct {
	cache.Cache
}

// Listen ...
func (s *Server) Listen() {
	http.Handle("/cache/", s.getCacheHandler())
	http.Handle("/status", s.getStatusHandler())
	http.ListenAndServe("127.0.0.1:8000", nil)
}

// New ...
func New(c cache.Cache) *Server {
	return &Server{c}
}
