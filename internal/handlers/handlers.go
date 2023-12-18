package handlers

import (
	"net/http"

	"github.com/good-threads/backend/internal/logic/welcome"
)

type Handlers interface {
	Welcome(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	welcome welcome.Logic
}

func New(welcome welcome.Logic) Handlers {
	return &handlers{welcome: welcome}
}

func (h *handlers) Welcome(w http.ResponseWriter, r *http.Request) {
	result := h.welcome.Behavior()
	w.Write([]byte(result))
}
