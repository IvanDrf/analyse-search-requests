package http

import (
	"fmt"
	"log"
	"net/http"
)

type searchServer struct {
	host string
	port int

	mux      http.ServeMux
	server   http.Server
	handlers *handlers
}

func NewSearchServer(host string, port int, handlers *handlers) *searchServer {
	return &searchServer{
		host: host,
		port: port,

		mux:      *http.NewServeMux(),
		server:   http.Server{},
		handlers: handlers,
	}
}

func (s *searchServer) RegisterRoutes() {
	s.mux.HandleFunc("POST /api/v1/bad", s.handlers.saveBadWord)
	s.mux.HandleFunc("DELETE /api/v1/bad", s.handlers.deleteBadWord)

	s.mux.HandleFunc("GET /api/v1/searches", s.handlers.findMostPopularSearches)
}

func (s *searchServer) Start() {
	s.server.Addr = fmt.Sprintf("%s:%d", s.host, s.port)

	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("can't start http server, error=%s", err)
	}
}
