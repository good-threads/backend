package session

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	Create(id string, username string) error
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

	c.mongoCollection.DeleteMany(context.TODO(), SessionSearchFilter{Username: username})

	result, err := c.mongoCollection.InsertOne(
		context.TODO(),
		Session{
			ID:             id,
			Username:       username,
			LastUpdateDate: time.Now(),
		},
	)
	if err != nil {
		log.Println("mongo error while inserting session:", err)
		return err
	}

	log.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return nil
}
