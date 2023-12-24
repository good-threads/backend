package session

import (
	"github.com/good-threads/backend/internal/client/session"
	"github.com/good-threads/backend/internal/client/user"
	e "github.com/good-threads/backend/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type Logic interface {
	Create(username string, password string) (string, error)
}

type logic struct {
	sessionClient session.Client
	userClient    user.Client
}

func Setup(sessionClient session.Client, userClient user.Client) Logic {
	return &logic{sessionClient: sessionClient, userClient: userClient}
}

func (l *logic) Create(username string, password string) (string, error) {
	if username == "" {
		return "", &e.BadUsername{}
	}
	if password == "" {
		return "", &e.BadPassword{}
	}
	user, err := l.userClient.Fetch(username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return "", &e.WrongCredentials{}
	}
	id := "TODO"
	err = l.sessionClient.Create(id, username)
	return id, err
}
