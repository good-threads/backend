package session

import mongoClient "github.com/good-threads/backend/internal/client/mongo"

type Session struct {
	ID             string               `bson:"id"`
	Username       string               `bson:"username"`
	LastUpdateDate mongoClient.NanoTime `bson:"lastUpdateDate"`
}
