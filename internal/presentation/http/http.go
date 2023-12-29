package http

import (
	"encoding/json"
	"net/http"

	e "github.com/good-threads/backend/internal/errors"
	"github.com/good-threads/backend/internal/logic/board"
	"github.com/good-threads/backend/internal/logic/common"
	"github.com/good-threads/backend/internal/logic/session"
	"github.com/good-threads/backend/internal/logic/user"
)

type Presentation interface {
	Ping(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	CreateSession(w http.ResponseWriter, r *http.Request)
	GetBoard(w http.ResponseWriter, r *http.Request)
}

type presentation struct {
	common  common.Logic
	user    user.Logic
	session session.Logic
	board   board.Logic
}

func Setup(common common.Logic, user user.Logic, session session.Logic, board board.Logic) Presentation {
	return &presentation{common: common, user: user, session: session, board: board}
}

func (p *presentation) Ping(w http.ResponseWriter, r *http.Request) {
	result := p.common.Ping()
	w.Write([]byte(result))
}

func (p *presentation) CreateUser(w http.ResponseWriter, r *http.Request) {

	var requestBody Credentials
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		respondMessage(w, http.StatusBadRequest, "Incorrect JSON format")
		return
	}

	err := p.user.Create(requestBody.Username, requestBody.Password)
	if err == nil {
		respondMessage(w, http.StatusCreated, "User created")
		return
	}
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
}

func (p *presentation) CreateSession(w http.ResponseWriter, r *http.Request) {

	var requestBody Credentials
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		respondMessage(w, http.StatusBadRequest, "Incorrect JSON format")
		return
	}

	id, err := p.session.Create(requestBody.Username, requestBody.Password)
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    id,
			HttpOnly: true,
		})
		respondMessage(w, http.StatusCreated, "Session created")
		return
	}
	switch err.(type) {
	case *e.BadPassword:
		respondMessage(w, http.StatusBadRequest, "Password must be ...")
	case *e.BadUsername:
		respondMessage(w, http.StatusBadRequest, "Username must be ...")
	case *e.WrongCredentials:
		respondMessage(w, http.StatusUnauthorized, "Wrong credentials")
	default:
		respondMessage(w, http.StatusInternalServerError, "Your request couldn't be processed")
	}
}

func (p *presentation) GetBoard(w http.ResponseWriter, r *http.Request) {

	username := r.Context().Value("username").(string)

	board, err := p.board.Get(username)
	if err == nil {
		respondMessage(w, http.StatusOK, board)
		return
	}
	switch err.(type) {
	default:
		respondMessage(w, http.StatusInternalServerError, "Your request couldn't be processed")
	}
}
