package http

import "net/http"

type searchServer struct {
	host string
	port int

	mux      http.ServeMux
	handlers *handlers
}

func NewSearchServer(host string, port int, handlers *handlers) *searchServer {
	return &searchServer{
		host: host,
		port: port,

		mux:      *http.NewServeMux(),
		handlers: handlers,
	}
}

func (s *searchServer) RegisterRoutes() {
	s.mux.HandleFunc("POST /api/v1/bad", s.handlers.saveBadWord)
	s.mux.HandleFunc("DELETE /api/v1/bad", s.handlers.deleteBadWord)

	s.mux.HandleFunc("GET /api/v1/searches", s.handlers.findMostPopularSearches)
}
