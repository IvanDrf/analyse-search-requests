package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/IvanDrf/analyse-search-requests/internal/infrastructure/adapters"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	server http.Server
	mux    *http.ServeMux
}

func NewMetricsServer(host string, port int) *MetricsServer {
	return &MetricsServer{
		mux:    http.NewServeMux(),
		server: http.Server{Addr: fmt.Sprintf("%s:%d", host, port)},
	}
}

func (s *MetricsServer) registerRoutes() {
	s.mux.Handle("GET /metrics", promhttp.Handler())

	s.server.Handler = s.mux
}

func (s *MetricsServer) Start() {
	adapters.RegisterMetrics()
	s.registerRoutes()

	slog.Info("MetricsServer:Start", slog.String("addr", s.server.Addr))
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("can't start metrics http server, error%s", err)
	}
}

func (s *MetricsServer) Stop(ctx context.Context) {
	s.server.Shutdown(ctx)

	slog.Info("MetricsServer:Stop", slog.String("status", "successfully stopped search server"))
}
