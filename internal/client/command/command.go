package command

import (
	"context"
	"time"

	mongoClient "github.com/good-threads/backend/internal/client/mongo"
	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	FetchLastID(username string) (*string, error)
	RegisterProcessed(username string, id string) error
}

type client struct {
	mongoCollection *mongo.Collection
}

func Setup(mongoClient *mongo.Client) Client {
	return &client{
		mongoCollection: mongoClient.Database("goodthreads").Collection("processed_commands"),
	}
}

func (c *client) FetchLastID(username string) (*string, error) {
	var command Command
	if err := c.mongoCollection.FindOne(
		context.TODO(),
		bson.M{"username": username},
		options.FindOne().SetSort(bson.M{"datetime": -1}),
	).Decode(&command); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &e.NoCommandFound{}
		}
		return nil, err
	}
	return &command.ID, nil
}

func (c *client) RegisterProcessed(username string, id string) error {
	_, err := c.mongoCollection.InsertOne(
		context.TODO(),
		Command{
			ID:       id, // TODO(thomasmarlow): unique index
			Username: username,
			Datetime: mongoClient.NanoTime{time.Now()}, // TODO(thomasmarlow): persist creationDatetime and processingDatetime
		},
	)
	return err
}
