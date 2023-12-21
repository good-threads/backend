package handlers

import (
	"encoding/json"
	"net/http"

	e "github.com/good-threads/backend/internal/errors"
	"github.com/good-threads/backend/internal/logic/common"
	"github.com/good-threads/backend/internal/logic/user"
)

type Handlers interface {
	Ping(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type handlers struct { // TODO(thomasmarlow): "HTTPPresentation" is a better name
	common common.Logic
	user   user.Logic
}

func New(common common.Logic, user user.Logic) Handlers {
	return &handlers{common: common, user: user}
}

func (h *handlers) Ping(w http.ResponseWriter, r *http.Request) {
	result := h.common.Ping()
	w.Write([]byte(result))
}

func (h *handlers) CreateUser(w http.ResponseWriter, r *http.Request) {

	var requestBody CreateUserRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		respondMessage(w, http.StatusBadRequest, "Incorrect JSON format")
		return
	}

	if err := h.user.Create(requestBody.Username, requestBody.Password); err == nil {
		respondMessage(w, http.StatusCreated, "User created")
	} else {
		switch err.(type) {
		case *e.UsernameAlreadyTaken:
			respondMessage(w, http.StatusConflict, "Username already taken")
		case *e.BadPassword:
			respondMessage(w, http.StatusBadRequest, "Password must be ...")
		case *e.BadUsername:
			respondMessage(w, http.StatusBadRequest, "Username must be ...")
		default:
			respondMessage(w, http.StatusInternalServerError, "Your request couldn't be processed")
		}
		return
	}
}
