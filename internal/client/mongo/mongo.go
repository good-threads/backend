package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Setup(mongoDBURI string) *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		log.Fatalf("unable connect to mongo: %e", err)
	}
	return client
}
