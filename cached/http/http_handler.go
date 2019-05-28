package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type cacheHandler struct {
	*Server
}

func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPut:
		data, _ := ioutil.ReadAll(r.Body)
		if len(data) != 0 {
			err := h.Set(key, data)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	case http.MethodGet:
		val, err := h.Get(key)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		if len(val) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write(val)
		return
	case http.MethodDelete:
		e := h.Del(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

func (s *Server) getCacheHandler() http.Handler {
	return &cacheHandler{s}
}
