package thread

import (
	"context"
	"errors"

	mongoClient "github.com/good-threads/backend/internal/client/mongo"
	e "github.com/good-threads/backend/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client interface {
	FetchAll(username string, ids []string) ([]Thread, []string, error)
	Create(transaction mongoClient.Transaction, username string, id string, name string) error
	EditName(transaction mongoClient.Transaction, username string, id string, name string) error
	AddKnot(transaction mongoClient.Transaction, username string, threadID string, knotID string, knotBody string) error
	EditKnotBody(transaction mongoClient.Transaction, username string, threadID string, knotID string, knotBody string) error
	DeleteKnot(transaction mongoClient.Transaction, username string, threadID string, knotID string) error
	FetchOne(username string, id string) (*Thread, error)
}

type client struct {
	mongoCollection *mongo.Collection
}

func Setup(mongoClient mongoClient.Client) Client {
	return &client{
		mongoCollection: mongoClient.MongoClient().Database("goodthreads").Collection("threads"),
	}
}

func (c *client) FetchAll(username string, ids []string) ([]Thread, []string, error) {

	cursor, err := c.mongoCollection.Aggregate(context.TODO(), mongo.Pipeline{
		{
			{"$match", bson.M{"id": bson.M{"$in": ids}}},
		},
		{
			{"$addFields", bson.M{"order": bson.M{"$indexOfArray": bson.A{ids, "$id"}}}},
		},
		{
			{"$sort", bson.M{"order": 1}},
		},
	})
	if err != nil && err != mongo.ErrNoDocuments {
		return []Thread{}, nil, err
	}

	activeThreads := []Thread{}
	err = cursor.All(context.TODO(), &activeThreads)
	if err != nil {
		return []Thread{}, nil, err
	}

	unassertedHiddenThreadIDs, err := c.mongoCollection.Distinct(context.TODO(), "id", bson.M{
		"username": username,
		"id":       bson.M{"$nin": ids},
	})
	if err != nil && err != mongo.ErrNoDocuments {
		return activeThreads, nil, err
	}

	hiddenThreadIDs := make([]string, 0)
	for _, uncastedID := range unassertedHiddenThreadIDs {
		id, ok := uncastedID.(string)
		if !ok {
			return activeThreads, nil, &e.ThreadIDIsNotString{}
		}
		hiddenThreadIDs = append(hiddenThreadIDs, id)
	}

	return activeThreads, hiddenThreadIDs, err
}

func (c *client) Create(transaction mongoClient.Transaction, username string, id string, name string) error {
	_, err := c.mongoCollection.InsertOne(transaction, Thread{
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

func (c *client) EditName(transaction mongoClient.Transaction, username string, id string, name string) error {
	result := c.mongoCollection.FindOneAndUpdate(transaction,
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

func (c *client) AddKnot(transaction mongoClient.Transaction, username string, threadID string, knotID string, knotBody string) error {
	filter := bson.M{
		"username": username,
		"id":       threadID,
	}
	result := c.mongoCollection.FindOneAndUpdate(transaction,
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

func (c *client) EditKnotBody(transaction mongoClient.Transaction, username string, threadID string, knotID string, knotBody string) error {
	result := c.mongoCollection.FindOneAndUpdate(transaction,
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

func (c *client) DeleteKnot(transaction mongoClient.Transaction, username string, threadID string, knotID string) error {
	result := c.mongoCollection.FindOneAndUpdate(transaction,
		bson.M{
			"username": username,
			"id":       threadID,
		},
		bson.M{
			"$pull": bson.M{
				"knots": bson.M{
					"id": knotID,
				},
			},
		},
	)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return &e.ThreadNotFound{}
	}
	// TODO(thomasmarlow): knot not found (no docs updated)
	return result.Err()
}

func (c *client) FetchOne(username string, id string) (*Thread, error) {
	var thread Thread
	err := c.mongoCollection.FindOne(context.TODO(), bson.M{
		"username": username,
		"id":       id,
	}).Decode(&thread)
	if err == mongo.ErrNoDocuments {
		return nil, &e.ThreadNotFound{}
	}
	return &thread, err
}
