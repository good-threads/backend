package http

import "github.com/good-threads/backend/internal/client/thread"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Board struct {
	Threads                  []thread.Thread `json:"threads"`
	LastProcessedChangesetID *string         `json:"lastProcessedChangesetID"`
}
