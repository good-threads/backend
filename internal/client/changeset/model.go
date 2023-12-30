package changeset

import "time"

type Changeset struct {
	ID       string    `bson:"id"`
	Username string    `bson:"username"`
	Datetime time.Time `bson:"datetime"`
}
