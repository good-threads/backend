package user

import (
	"context"
	"errors"

	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Persist(username string, passwordHash []byte) error
	Fetch(username string) (*User, error)
	AddThread(username string, id string) error
	RemoveThread(username string, id string) error
	RelocateThread(username string, id string, newIndex uint) error
}

type client struct {
	mongoCollection *mongo.Collection
}

func Setup(mongoClient *mongo.Client) Client {
	return &client{
		mongoCollection: mongoClient.Database("goodthreads").Collection("users"),
	}
}

func (c *client) Persist(username string, passwordHash []byte) error {
	_, err := c.mongoCollection.InsertOne(
		context.TODO(),
		User{
			Name:         username,
			PasswordHash: passwordHash,
			Threads:      []string{},
		},
	)
	if mongo.IsDuplicateKeyError(err) {
		return &e.UsernameAlreadyTaken{}
	}
	return err
}

func (c *client) Fetch(username string) (*User, error) {
	var user User
	if err := c.mongoCollection.FindOne(context.TODO(), bson.M{"name": username}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &e.UserNotFound{}
		}
		return nil, err
	}
	return &user, nil
}

func (c *client) AddThread(username string, id string) error {
	result := c.mongoCollection.FindOneAndUpdate(context.TODO(),
		bson.M{
			"name": username,
		},
		bson.M{
			"$push": bson.M{
				"threads": id,
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.UserNotFound{}
	}
	return result.Err()
}

func (c *client) RemoveThread(username string, id string) error {

	result := c.mongoCollection.FindOneAndUpdate(context.TODO(),
		bson.M{
			"name": username,
		},
		bson.M{
			"$pull": bson.M{
				"threads": id,
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.UserNotFound{}
	}
	return result.Err()
}

func (c *client) RelocateThread(username string, id string, newIndex uint) error {

	result := c.mongoCollection.FindOneAndUpdate(context.TODO(),
		bson.M{
			"name": username,
		},
		bson.M{
			"$push": bson.M{
				"threads": bson.M{
					"$each":     bson.A{id},
					"$position": newIndex,
				},
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.UserNotFound{}
	}
	return result.Err()
}
