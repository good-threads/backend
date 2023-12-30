package user

import (
	"context"

	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Persist(username string, passwordHash []byte) error
	Fetch(username string) (*User, error)
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
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &e.UsernameAlreadyTaken{}
		} else {
			return err
		}
	}

	return nil
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
