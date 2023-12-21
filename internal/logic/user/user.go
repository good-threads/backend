package user

import (
	e "github.com/good-threads/backend/internal/errors"
)

type Logic interface {
	Create(username string, password string) error
}

type logic struct{}

func New() Logic {
	return &logic{}
}

func (l *logic) Create(username string, password string) error {
	if username == "" {
		return &e.BadUsername{}
	}
	if password == "" {
		return &e.BadPassword{}
	}
	if username == "tom" {
		return &e.UsernameAlreadyTaken{}
	}
	return nil
}
