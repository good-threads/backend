package board

import (
	mongoClient "github.com/good-threads/backend/internal/client/mongo"
)

type Command struct {
	ID       string               `json:"id"`
	Type     string               `json:"type"`
	Datetime mongoClient.NanoTime `json:"datetime"` // TODO(thomasmarlow): unused
	Payload  any                  `json:"payload"`
}

type PayloadCreateThread struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadEditThreadName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PayloadHideThread struct {
	ID string `json:"id"`
}

type PayloadRelocateThread struct {
	ID       string `json:"id"`
	NewIndex uint   `json:"newIndex"`
}

type PayloadCreateKnot struct {
	ThreadID string `json:"threadID"`
	KnotID   string `json:"knotID"`
	KnotBody string `json:"knotBody"`
}

type PayloadEditKnotBody struct {
	ThreadID string `json:"threadID"`
	KnotID   string `json:"knotID"`
	KnotBody string `json:"knotBody"`
}

type PayloadDeleteKnot struct {
	ThreadID string `json:"threadID"`
	KnotID   string `json:"knotID"`
}
