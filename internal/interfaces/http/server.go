package http

import (
	"fmt"
	"log"
	"net/http"
)

type SearchServer struct {
	host string
	port int

	mux      http.ServeMux
	server   http.Server
	handlers *handlers
}

func NewSearchServer(host string, port int, handlers *handlers) *SearchServer {
	return &SearchServer{
		host: host,
		port: port,

		mux:      *http.NewServeMux(),
		server:   http.Server{},
		handlers: handlers,
	}
}

func (s *SearchServer) registerRoutes() {
	s.mux.HandleFunc("POST /api/v1/bad", s.handlers.saveBadWord)
	s.mux.HandleFunc("DELETE /api/v1/bad", s.handlers.deleteBadWord)

	s.mux.HandleFunc("GET /api/v1/searches", s.handlers.findMostPopularSearches)
}

func (s *SearchServer) Start() {
	s.registerRoutes()

	s.server.Addr = fmt.Sprintf("%s:%d", s.host, s.port)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("can't start http server, error=%s", err)
	}
}

func (s *SearchServer) Stop() {
	s.server.Close()
	s.handlers.close()
}
