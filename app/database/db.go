package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func dbInstance() (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGODB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func Connect(database string) (*mongo.Database, error) {
	client, err := dbInstance()
	if err != nil {
		return nil, err
	}

	return client.Database(database), nil
}

func GetCollection(database *mongo.Database, collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}
