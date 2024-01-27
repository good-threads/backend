package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Client interface {
	MongoClient() *mongo.Client
	Transactionally(operations func(Transaction) error) error
}

type client struct {
	mongoClient *mongo.Client
}

type Transaction mongo.SessionContext

func Setup(mongoDBURI string) Client {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		log.Fatalf("unable connect to mongo: %e", err)
	}
	return &client{mongoClient: mongoClient}
}

func (c *client) MongoClient() *mongo.Client {
	return c.mongoClient
}

func (c *client) Transactionally(operations func(Transaction) error) error {

	session, err := c.mongoClient.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		return nil, operations(ctx)
	}, options.Transaction().SetWriteConcern(writeconcern.Majority()))

	return err
}
