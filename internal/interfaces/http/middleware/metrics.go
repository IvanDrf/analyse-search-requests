package middleware

import (
	"net/http"
	"time"

	"github.com/IvanDrf/analyse-search-requests/internal/infrastructure/adapters"
)

func PrometheusMiddleware(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler(w, r)

		adapters.HttpRequestsAmount.WithLabelValues(r.Method, r.URL.Path).Inc()
		adapters.HttpRequestsDuration.WithLabelValues(r.Method, r.URL.Path).Observe(float64(time.Since(start).Seconds()))
	}
}
