package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Client interface {
	GetHandler() http.Handler
	RegisterBoardRead()
}

type client struct {
	registry   *prometheus.Registry
	boardReads prometheus.Counter
}

func Setup() Client {
	client := &client{
		boardReads: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "backend",
			Name:      "board_reads",
			Help:      "Number of times a board's state was read.",
		}),
	}
	registry := prometheus.NewRegistry()
	registry.MustRegister(client.boardReads)
	return client
}

func (c *client) GetHandler() http.Handler {
	return promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{})
}

func (c *client) RegisterBoardRead() {
	c.boardReads.Inc()
}
