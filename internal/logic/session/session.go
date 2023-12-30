package session

import (
	"github.com/good-threads/backend/internal/client/session"
	"github.com/good-threads/backend/internal/client/user"
	e "github.com/good-threads/backend/internal/errors"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type Logic interface {
	Create(username string, password string) (string, error)
	GetUsername(sessionID string) (string, error)
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
		switch err.(type) {
		case *e.UserNotFound:
			return "", &e.WrongCredentials{}
		default:
			return "", err
		}
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return "", &e.WrongCredentials{}
	}
	id, err := ksuid.NewRandom()
	if err != nil {
		return "", err
	}
	err = l.sessionClient.Create(id.String(), username)
	return id.String(), err
}

func (l *logic) GetUsername(sessionID string) (string, error) {
	_, err := ksuid.Parse(sessionID)
	if err != nil {
		return "", &e.InvalidSession{}
	}
	session, err := l.sessionClient.Fetch(sessionID)
	if err != nil {
		return "", err
	}
	return session.Username, nil
}
