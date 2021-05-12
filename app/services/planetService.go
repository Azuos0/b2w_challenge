package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Azuos0/b2w_challenge/app/database"
	"github.com/Azuos0/b2w_challenge/app/models"
	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlanetService struct {
	Collection *mongo.Collection
}

type SearchResponse struct {
	Page      int64            `json:"page"`
	PerPage   int64            `json:"perPage"`
	Prev      int64            `json:"prev"`
	Next      int64            `json:"next"`
	Total     int64            `json:"total"`
	TotalPage int64            `json:"totalPage"`
	Result    *[]models.Planet `json:"result"`
}

type swapiPlanetResponse struct {
	Count   int           `json:"count"`
	Results []swapiPlanet `json:"results"`
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

func NewPlanetService(db *mongo.Database) *PlanetService {
	client := &PlanetService{
		Collection: database.GetCollection(db, "planets"),
	}

	return client
}

func (client *PlanetService) Create(planet models.Planet) (*models.Planet, error) {
	planet.ID = primitive.NewObjectID()
	planet.Appearances, _ = getPlanetNumberOfApperances(planet.Name)
	planet.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	res, err := client.Collection.InsertOne(ctx, planet)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID.(primitive.ObjectID).Hex()
	return client.Get(id)
}

func (client *PlanetService) Get(id string) (*models.Planet, error) {
	planet := models.Planet{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	err = client.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&planet)
	if err != nil {
		return nil, err
	}

	return &planet, nil
}

func (client *PlanetService) Delete(id string) (string, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	res, err := client.Collection.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		return "", err
	}

	if res.DeletedCount > 0 {
		return "Planet was deleted successfully!", nil
	} else {
		err = errors.New("no planet with this id was found in this so far far away galaxy")
		return "", err
	}

}

func getPlanetNumberOfApperances(name string) (int, error) {
	swapiRes := swapiPlanetResponse{}
	name = strings.ReplaceAll(name, " ", "%20") //replace whitespaces for equivalent %20 for url

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

	planetName := swapiRes.Results[0].Name

	//if planets have different names
	if !strings.EqualFold(name, planetName) {
		return 0, nil
	}

	return len(swapiRes.Results[0].Films), nil
}

func (client *PlanetService) Search(page int64, name string) (*SearchResponse, error) {
	planets := []models.Planet{}
	var filter bson.M

	if name == "" {
		filter = bson.M{}
	} else {
		filter = bson.M{"name": bson.M{"$regex": name, "$options": "im"}}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	paginatedData, err := mongopagination.New(client.Collection).Context(ctx).Limit(30).Page(page).Filter(filter).Decode(&planets).Find()
	if err != nil {
		return nil, err
	}

	result := SearchResponse{
		Page:      paginatedData.Pagination.Page,
		Next:      paginatedData.Pagination.Next,
		Prev:      paginatedData.Pagination.Prev,
		PerPage:   paginatedData.Pagination.PerPage,
		Total:     paginatedData.Pagination.Total,
		TotalPage: paginatedData.Pagination.TotalPage,
		Result:    &planets,
	}

	return &result, nil
}
