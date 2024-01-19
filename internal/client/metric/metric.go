package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Client interface {
	RegisterBoardRead()
}

type client struct {
	boardReads prometheus.Counter
}

func Setup(registry *prometheus.Registry) Client {
	client := &client{
		boardReads: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "backend",
			Name:      "board_reads",
			Help:      "Number of times a board's state was read.",
		}),
	}
	registry.MustRegister(client.boardReads)
	return client
}

func (c *client) RegisterBoardRead() {
	c.boardReads.Inc()
}
