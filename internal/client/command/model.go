package command

import mongoClient "github.com/good-threads/backend/internal/client/mongo"

type Command struct {
	ID       string               `bson:"id"`
	Username string               `bson:"username"`
	Datetime mongoClient.NanoTime `bson:"datetime"`
}
