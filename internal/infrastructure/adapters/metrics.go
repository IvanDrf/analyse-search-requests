package adapters

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestsAmount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "HTTP_REQUESTS_AMOUNT",
			Help: "Amount number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	HttpRequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "HTTP_REQUESTS_DURATION",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(HttpRequestsAmount)
	prometheus.MustRegister(HttpRequestsDuration)
}
