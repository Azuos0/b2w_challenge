package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func dbInstance() (*mongo.Client, context.Context, error) {
	mongoURI := os.Getenv("MONGODB_URL")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, nil, err
	}

	ctx := context.TODO()

	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}

func Connect(database string) (*mongo.Database, context.Context, error) {
	client, ctx, err := dbInstance()
	if err != nil {
		return nil, nil, err
	}

	return client.Database(database), ctx, nil
}

func GetCollection(database *mongo.Database, collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}
