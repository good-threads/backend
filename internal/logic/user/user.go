package user

import (
	"github.com/good-threads/backend/internal/client/user"
	e "github.com/good-threads/backend/internal/errors"
	"golang.org/x/crypto/bcrypt"
)

type Logic interface {
	Create(username string, password string) error
}

type logic struct {
	client user.Client
}

func Setup(client user.Client) Logic {
	return &logic{client: client}
}

func (l *logic) Create(username string, password string) error {
	if username == "" {
		return &e.BadUsername{}
	}
	if password == "" {
		return &e.BadPassword{}
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	return l.client.Persist(username, bytes)
}
