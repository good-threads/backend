package user

import (
	"log"

	e "github.com/good-threads/backend/internal/errors"
)

type Client interface {
	Persist(username string, passwordHash []byte) error
}

type client struct {
	takenUsername string
}

func Setup(takenUsername string) Client {
	return &client{takenUsername: takenUsername}
}

func (c *client) Persist(username string, passwordHash []byte) error {
	if username == c.takenUsername {
		return &e.UsernameAlreadyTaken{}
	}
	log.Println(username, passwordHash)
	return nil
}
