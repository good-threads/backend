package http

import (
	"github.com/good-threads/backend/internal/client/thread"
	"github.com/good-threads/backend/internal/logic/board"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Board struct {
	ActiveThreads          []thread.Thread `json:"activeThreads"`
	HiddenThreads          []string        `json:"hiddenThreads"`
	LastProcessedCommandID *string         `json:"lastProcessedCommandID"`
}

type BoardUpdates struct {
	LastProcessedCommandID *string         `json:"lastProcessedCommandID"`
	Commands               []board.Command `json:"commands"`
}

type BoardUpdateOKResponse struct {
	LastProcessedCommandID *string `json:"lastProcessedCommandID"`
}
