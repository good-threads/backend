package http

import (
	"encoding/json"
	"fmt"
	"log"
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
	UpdateBoard(w http.ResponseWriter, r *http.Request)
	GetUsernameFromSession(wrappedHandler http.Handler) http.Handler
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
		log.Println("ERROR:", fmt.Sprintf("%T", err), err)
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
		log.Println("ERROR:", fmt.Sprintf("%T", err), err)
		respondMessage(w, http.StatusInternalServerError, "Your request couldn't be processed")
	}
}

func (p *presentation) GetBoard(w http.ResponseWriter, r *http.Request) {

	username, ok := r.Context().Value("username").(string)
	if !ok {
		respondMessage(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	threads, lastProcessedCommandID, err := p.board.Get(username)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Board{
			Threads:                threads,
			LastProcessedCommandID: lastProcessedCommandID,
		})
		return
	}
	switch err.(type) {
	default:
		log.Println("ERROR:", fmt.Sprintf("%T", err), err)
		respondMessage(w, http.StatusInternalServerError, "Your request couldn't be processed")
	}
}

func (p *presentation) UpdateBoard(w http.ResponseWriter, r *http.Request) {

	username, ok := r.Context().Value("username").(string) // TODO(thomasmarlow): de-duplicate these 5 lines
	if !ok {
		respondMessage(w, http.StatusUnauthorized, "Invalid session")
		return
	}

	var requestBody BoardUpdates
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("ERROR:", fmt.Sprintf("%T", err), err)
		respondMessage(w, http.StatusBadRequest, "Incorrect JSON format")
		return
	}

	lastProcessedCommandID, err := p.board.Update(username, requestBody.LastProcessedCommandID, requestBody.Commands)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(BoardUpdateOKResponse{
			LastProcessedCommandID: lastProcessedCommandID,
		})
		return
	}
	switch err.(type) {
	case *e.ReceivedCommandsWouldRewriteHistory:
		respondMessage(w, http.StatusPreconditionFailed, "Requested commands would rewrite history - maybe another client is updating the same board; please discard commands and pull the latest board state")
	default:
		log.Println("ERROR:", fmt.Sprintf("%T", err), err)
		respondMessage(w, http.StatusInternalServerError, "Your request couldn't be processed")
	}
}
