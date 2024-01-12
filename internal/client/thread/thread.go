package thread

import (
	"context"
	"errors"

	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Fetch(ids []string) ([]Thread, error)
	Create(username string, id string, name string) error
	EditName(username string, id string, name string) error
	AddKnot(username string, threadID string, knotID string, knotBody string) error
	EditKnot(username string, threadID string, knotID string, knotBody string) error
	DeleteKnot(username string, threadID string, knotID string) error
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
	err = cursor.All(context.TODO(), &threads)

	return threads, err
}

func (c *client) Create(username string, id string, name string) error {
	_, err := c.mongoCollection.InsertOne(context.TODO(), Thread{
		ID:       id,
		Name:     name,
		Username: username,
		Knots:    []Knot{},
	})
	if mongo.IsDuplicateKeyError(err) {
		return &e.GeneratedIDClashed{}
	}
	return err
}

func (c *client) EditName(username string, id string, name string) error {
	result := c.mongoCollection.FindOneAndUpdate(context.TODO(),
		bson.M{
			"username": username,
			"id":       id,
		},
		bson.M{
			"$set": bson.M{
				"name": name,
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.ThreadNotFound{}
	}
	return result.Err()
}

func (c *client) AddKnot(username string, threadID string, knotID string, knotBody string) error {
	filter := bson.M{
		"username": username,
		"id":       threadID,
	}
	result := c.mongoCollection.FindOneAndUpdate(context.TODO(),
		filter,
		bson.M{
			"$push": bson.M{
				"knots": Knot{
					ID:   knotID,
					Body: knotBody,
				},
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.ThreadNotFound{}
	}
	return result.Err()
}

func (c *client) EditKnot(username string, threadID string, knotID string, knotBody string) error {
	result := c.mongoCollection.FindOneAndUpdate(context.TODO(),
		bson.M{
			"username": username,
			"id":       threadID,
			"knots.id": knotID,
		},
		bson.M{
			"$set": bson.M{
				"knots.$.body": knotBody,
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.KnotNotFound{}
	}
	return result.Err()
}

func (c *client) DeleteKnot(username string, threadID string, knotID string) error {
	result := c.mongoCollection.FindOneAndUpdate(context.TODO(),
		bson.M{
			"username": username,
			"id":       threadID,
		},
		bson.M{
			"$pull": bson.M{
				"knots": knotID,
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.ThreadNotFound{}
	}
	// TODO(thomasmarlow): knot not found (no docs updated)
	return result.Err()
}
