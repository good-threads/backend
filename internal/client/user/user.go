package user

import (
	"context"
	"errors"

	mongoClient "github.com/good-threads/backend/internal/client/mongo"
	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Persist(username string, passwordHash []byte) error
	Fetch(username string) (*User, error)
	AddThread(transaction mongoClient.Transaction, username string, id string) error
	RemoveThread(transaction mongoClient.Transaction, username string, id string) error
	RelocateThread(transaction mongoClient.Transaction, username string, id string, newIndex uint) error
}

type client struct {
	mongoCollection *mongo.Collection
}

func Setup(mongoClient mongoClient.Client) Client {
	return &client{
		mongoCollection: mongoClient.MongoClient().Database("goodthreads").Collection("users"),
	}
}

func (c *client) Persist(username string, passwordHash []byte) error {
	_, err := c.mongoCollection.InsertOne(
		context.TODO(),
		User{
			Name:          username,
			PasswordHash:  passwordHash,
			ActiveThreads: []string{},
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

func (c *client) AddThread(transaction mongoClient.Transaction, username string, id string) error {
	result := c.mongoCollection.FindOneAndUpdate(transaction,
		bson.M{
			"name": username,
		},
		bson.M{
			"$push": bson.M{
				"activeThreads": id,
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.UserNotFound{}
	}
	return result.Err()
}

func (c *client) RemoveThread(transaction mongoClient.Transaction, username string, id string) error {

	result := c.mongoCollection.FindOneAndUpdate(transaction,
		bson.M{
			"name": username,
		},
		bson.M{
			"$pull": bson.M{
				"activeThreads": id,
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.UserNotFound{}
	}
	return result.Err()
}

func (c *client) RelocateThread(transaction mongoClient.Transaction, username string, id string, newIndex uint) error {

	result := c.mongoCollection.FindOneAndUpdate(transaction,
		bson.M{
			"name": username,
		},
		bson.M{
			"$pull": bson.M{
				"activeThreads": id,
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.UserNotFound{}
	}

	result = c.mongoCollection.FindOneAndUpdate(transaction,
		bson.M{
			"name": username,
		},
		bson.M{
			"$push": bson.M{
				"activeThreads": bson.M{
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
