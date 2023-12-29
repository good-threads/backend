package session

import "time"

type Session struct {
	ID             string    `bson:"id"`
	Username       string    `bson:"username"`
	LastUpdateDate time.Time `bson:"last_update_date"`
}

type SessionSearchFilter struct {
	ID       string `bson:"id,omitempty"`
	Username string `bson:"username,omitempty"`
}
