package user

import (
	e "github.com/good-threads/backend/internal/errors"
)

type Logic interface {
	Create(username string, password string) error
}

type logic struct {
	takenUsername string
}

func Setup(takenUsername string) Logic {
	return &logic{takenUsername: takenUsername}
}

func (l *logic) Create(username string, password string) error {
	if username == "" {
		return &e.BadUsername{}
	}
	if password == "" {
		return &e.BadPassword{}
	}
	if username == l.takenUsername {
		return &e.UsernameAlreadyTaken{}
	}
	return nil
}
