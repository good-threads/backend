package thread

import (
	"github.com/good-threads/backend/internal/client/thread"
)

type Logic interface {
	Get(username string, id string) (*thread.Thread, error)
}

type logic struct {
	client thread.Client
}

func Setup(client thread.Client) Logic {
	return &logic{client: client}
}

func (l *logic) Get(username string, id string) (*thread.Thread, error) {
	return l.client.FetchOne(username, id)
}
