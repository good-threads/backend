package handlers

import (
	"net/http"

	"github.com/good-threads/backend/internal/logic/common"
)

type Handlers interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	common common.Logic
}

func New(common common.Logic) Handlers {
	return &handlers{common: common}
}

func (h *handlers) Ping(w http.ResponseWriter, r *http.Request) {
	result := h.common.Ping()
	w.Write([]byte(result))
}
