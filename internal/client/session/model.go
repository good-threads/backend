package session

import "time"

type Session struct {
	ID             string    `bson:"id"`
	Username       string    `bson:"username"`
	LastUpdateDate time.Time `bson:"lastUpdateDate"`
}
