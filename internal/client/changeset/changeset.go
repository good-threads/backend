package changeset

import (
	"context"

	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client interface {
	FetchLastID(username string) (*string, error)
}

type client struct {
	mongoCollection *mongo.Collection
}

func Setup(mongoClient *mongo.Client) Client {
	return &client{
		mongoCollection: mongoClient.Database("goodthreads").Collection("processed_changesets"),
	}
}

func (c *client) FetchLastID(username string) (*string, error) {
	var changeset Changeset
	if err := c.mongoCollection.FindOne(
		context.TODO(),
		bson.M{"username": username},
		options.FindOne().SetSort(bson.M{"datetime": -1}),
	).Decode(&changeset); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &e.NoChangesetFound{}
		}
		return nil, err
	}
	return &changeset.ID, nil
}
