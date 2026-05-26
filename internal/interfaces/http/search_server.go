package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/analyse-search-requests/internal/interfaces/http/middleware"
)

type SearchServer struct {
	mux      *http.ServeMux
	server   http.Server
	handlers *handlers
}

func NewSearchServer(host string, port int, handlers *handlers) *SearchServer {
	return &SearchServer{
		mux:      http.NewServeMux(),
		server:   http.Server{Addr: fmt.Sprintf("%s:%d", host, port)},
		handlers: handlers,
	}
}

func (s *SearchServer) registerRoutes() {
	s.mux.HandleFunc("POST /api/v1/bad", middleware.PrometheusMiddleware(s.handlers.saveBadWord))
	s.mux.HandleFunc("DELETE /api/v1/bad", middleware.PrometheusMiddleware(s.handlers.deleteBadWord))

	s.mux.HandleFunc("GET /api/v1/searches", middleware.PrometheusMiddleware(s.handlers.findMostPopularSearches))

	s.server.Handler = s.mux
}

func (s *SearchServer) Start() {
	s.registerRoutes()

	slog.Info("SearchServer:Start", slog.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("can't start http server, error=%s", err)
	}
}

func (s *SearchServer) Stop(ctx context.Context) {
	s.server.Shutdown(ctx)
	s.handlers.close()

	slog.Info("SearchServer:Stop", slog.String("status", "successfully stopped search server"))
}
