package session

import (
	"context"
	"time"

	mongoClient "github.com/good-threads/backend/internal/client/mongo"

	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Create(id string, username string) error
	Fetch(sessionID string) (*Session, error)
}

type client struct {
	mongoCollection *mongo.Collection
}

func Setup(mongoClient *mongo.Client) Client {
	return &client{
		mongoCollection: mongoClient.Database("goodthreads").Collection("sessions"),
	}
}

func (c *client) Create(id string, username string) error {

	c.mongoCollection.DeleteMany(context.TODO(), bson.M{"username": username})

	_, err := c.mongoCollection.InsertOne(
		context.TODO(),
		Session{
			ID:             id,
			Username:       username,
			LastUpdateDate: mongoClient.NanoTime{time.Now()},
		},
	)

	return err
}

func (c *client) Fetch(sessionID string) (*Session, error) {
	var session Session
	if err := c.mongoCollection.FindOne(context.TODO(), bson.M{"id": sessionID}).Decode(&session); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &e.SessionNotFound{}
		}
		return nil, err
	}
	return &session, nil
}
