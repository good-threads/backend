package metric

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Client interface {
	GetHandler() http.Handler
	Middleware(next http.Handler) http.Handler
	RegisterBoardRead()
}

type client struct {
	registry   *prometheus.Registry
	collectors map[string]prometheus.Collector
}

func Setup() Client {
	collectors := map[string]prometheus.Collector{
		"requests": prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "backend",
				Name:      "http_requests",
				Help:      "Number of HTTP requests.",
			},
			[]string{"method", "route", "status"},
		),
		"boardReads": prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "backend",
			Name:      "board_reads",
			Help:      "Number of times a board's state was read.",
		}),
	}
	client := &client{
		registry:   prometheus.NewRegistry(),
		collectors: collectors,
	}
	collectorList := []prometheus.Collector{}
	for _, collector := range collectors {
		collectorList = append(collectorList, collector)
	}
	client.registry.MustRegister(collectorList...)
	return client
}

func (c *client) GetHandler() http.Handler {
	return promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{})
}

func (c *client) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &statusCapturingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(wrapped, r)
		c.collectors["requests"].(*prometheus.CounterVec).WithLabelValues(r.Method, r.URL.Path, fmt.Sprint(wrapped.status)).Inc()
	})
}

type statusCapturingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusCapturingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (c *client) RegisterBoardRead() {
	c.collectors["boardReads"].(prometheus.Counter).Inc()
}
