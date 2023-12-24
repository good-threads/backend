package user

import (
	"context"
	"log"

	e "github.com/good-threads/backend/internal/errors"
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

	result, err := c.mongoCollection.InsertOne(
		context.TODO(),
		User{
			Name:         username,
			PasswordHash: passwordHash,
		},
	)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &e.UsernameAlreadyTaken{}
		} else {
			log.Println("mongo error while inserting user:", err)
			return err
		}
	}

	log.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return nil
}

func (c *client) Fetch(username string) (*User, error) {
	result := c.mongoCollection.FindOne(context.TODO(), UserSearchFilter{Username: username})
	return nil, nil
}
