package handlers

import (
	"net/http"

	"github.com/good-threads/backend/internal/logic"
)

type test interface {
	Test(w http.ResponseWriter, r *http.Request)
}

type Test struct {
	l *logic.Test
}

func NewTest(l *logic.Test) *Test {
	return &Test{l: l}
}

func (h *Test) Handler(w http.ResponseWriter, r *http.Request) {
	result := h.l.Behavior()
	w.Write([]byte(result))
}
