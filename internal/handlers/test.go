package handlers

import (
	"net/http"

	"github.com/good-threads/backend/internal/logic/welcome"
)

type Handlers interface {
	Welcome(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	logic welcome.Logic
}

func New(logic welcome.Logic) Handlers {
	return &handlers{logic: logic}
}

func (h *handlers) Welcome(w http.ResponseWriter, r *http.Request) {
	result := h.logic.Behavior()
	w.Write([]byte(result))
}
