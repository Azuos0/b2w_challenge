package services

import (
	"context"
	"fmt"

	"github.com/Azuos0/b2w_challenge/app/database"
	"github.com/Azuos0/b2w_challenge/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlanetClient struct {
	Ctx        context.Context
	Collection *mongo.Collection
}

func NewPlanetClient(ctx context.Context, db *mongo.Database) *PlanetClient {
	client := &PlanetClient{
		Ctx:        ctx,
		Collection: database.GetCollection(db, "planets"),
	}

	return client
}

func (client *PlanetClient) Create(p models.Planet) (models.Planet, error) {
	planet := models.Planet{}

	res, err := client.Collection.InsertOne(client.Ctx, p)
	if err != nil {
		return planet, err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	return client.Get(id)
}

func (client *PlanetClient) Get(id string) (models.Planet, error) {
	planet := models.Planet{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return planet, err
	}

	err = client.Collection.FindOne(client.Ctx, bson.M{"_id": _id}).Decode(&planet)
	if err != nil {
		return planet, err
	}

	return planet, nil
}

func (client *PlanetClient) Search(filter interface{}) ([]models.Planet, error) {
	planets := []models.Planet{}

	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := client.Collection.Find(client.Ctx, filter)
	if err != nil {
		return planets, err
	}

	for cursor.Next(client.Ctx) {
		row := models.Planet{}
		cursor.Decode(&row)
		planets = append(planets, row)
	}

	return planets, nil
}

func (client *PlanetClient) Delete(id string) (string, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	res, err := client.Collection.DeleteOne(client.Ctx, bson.M{"_id": _id})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v planet deleted successfully!", res.DeletedCount), nil
}
