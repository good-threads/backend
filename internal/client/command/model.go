package command

import "time"

type Command struct {
	ID       string    `bson:"id"`
	Username string    `bson:"username"`
	Datetime time.Time `bson:"datetime"`
}
