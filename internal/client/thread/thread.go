package thread

import (
	"context"

	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Fetch(ids []string) ([]Thread, error)
}

type client struct {
	mongoCollection *mongo.Collection
}

func Setup(mongoClient *mongo.Client) Client {
	return &client{
		mongoCollection: mongoClient.Database("goodthreads").Collection("threads"),
	}
}

func (c *client) Fetch(ids []string) ([]Thread, error) {

	cursor, err := c.mongoCollection.Find(
		context.TODO(),
		bson.M{
			"id": bson.M{
				"$in": ids,
			},
		},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &e.NoThreadsFound{}
		}
		return nil, err
	}

	var threads []Thread
	err = cursor.All(context.TODO(), threads)

	return threads, err
}
