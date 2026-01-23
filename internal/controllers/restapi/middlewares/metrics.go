package middlewares

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ReqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "chat_service",
			Name:      "http_requests_total",
			Help:      "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
	ReqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "chat_service",
			Name:      "http_request_duration_seconds",
			Help:      "Request latency in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
	MessageSize = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "chat_service",
			Name:      "message_size_bytes",
			Help:      "Message size in bytes",
			Buckets:   prometheus.ExponentialBuckets(100, 2, 8),
		},
	)
)

func RegisterMetrics() {
	register := func(c prometheus.Collector) {
		if err := prometheus.Register(c); err != nil {
			var are prometheus.AlreadyRegisteredError
			if errors.As(err, &are) {
				return
			}
			panic(err)
		}
	}

	register(ReqCount)
	register(ReqDuration)
	register(MessageSize)
	register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	register(collectors.NewGoCollector())
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func NewMetricsMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &responseWriterWrapper{ResponseWriter: w}
			start := time.Now()
			next.ServeHTTP(rec, r)
			lat := time.Since(start).Seconds()

			ReqDuration.WithLabelValues(r.Method, r.RequestURI).Observe(lat)
			ReqCount.WithLabelValues(r.Method, r.RequestURI, strconv.Itoa(rec.status)).Inc()
		})
	}
}
