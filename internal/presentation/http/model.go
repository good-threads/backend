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
	Threads                  []thread.Thread `json:"threads"`
	LastProcessedChangesetID *string         `json:"lastProcessedChangesetID"`
}

type BoardUpdates struct {
	LastProcessedChangesetID *string           `json:"lastProcessedChangesetID"`
	NewChangesets            []board.Changeset `json:"newChangesets"`
}

type BoardUpdateOKResponse struct {
	LastProcessedChangesetID string
}
