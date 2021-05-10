package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

type swapiPlanetResponse struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Created  time.Time     `json:"created"`
	Edited   time.Time     `json:"edited"`
	Url      string        `json:"url"`
	Results  []swapiPlanet `json:"results"`
}

type swapiPlanet struct {
	Name            string   `json:"name"`
	Rotation_period string   `json:"rotation_period"`
	Orbital_period  string   `json:"orbital_period"`
	Diameter        string   `json:"diameter"`
	Climate         string   `json:"climate"`
	Gravity         string   `json:"gravity"`
	Terrain         string   `json:"terrain"`
	Surface_water   string   `json:"surface_water"`
	Population      string   `json:"population"`
	Residents       []string `json:"residents"`
	Films           []string `json:"films"`
}

func NewPlanetClient(ctx context.Context, db *mongo.Database) *PlanetClient {
	client := &PlanetClient{
		Ctx:        ctx,
		Collection: database.GetCollection(db, "planets"),
	}

	return client
}

func (client *PlanetClient) Create(planet models.Planet) (models.Planet, error) {
	p := models.Planet{}

	planet.ID = primitive.NewObjectID()
	planet.Appearances, _ = getPlanetNumberOfApperances(planet.Name)
	planet.CreatedAt = time.Now()

	res, err := client.Collection.InsertOne(client.Ctx, planet)
	if err != nil {
		return p, err
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

func (client *PlanetClient) Search(name string) ([]models.Planet, error) {
	planets := []models.Planet{}
	var filter bson.M

	if name == "" {
		// if filter == nil {
		filter = bson.M{}
	} else {
		filter = bson.M{"name": bson.M{"$regex": name, "$options": "im"}}
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

func getPlanetNumberOfApperances(name string) (int, error) {
	swapiRes := swapiPlanetResponse{}

	res, err := http.Get(fmt.Sprintf("https://swapi.dev/api/planets/?search=%v", name))
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &swapiRes)
	if err != nil {
		return 0, err
	}

	if swapiRes.Count > 1 || swapiRes.Count == 0 {
		return 0, nil
	}

	return len(swapiRes.Results[0].Films), nil
}
