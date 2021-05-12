package services_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/Azuos0/b2w_challenge/app/database"
	"github.com/Azuos0/b2w_challenge/app/models"
	"github.com/Azuos0/b2w_challenge/app/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func loadDatabase() *mongo.Database {
	dbName := os.Getenv("MONGODB_TEST_DATABASE")

	//load env variables if they are not loaded yet,
	if dbName == "" {
		godotenv.Load("../../.env")
		_ = os.Getenv("MONGODB_TEST_DATABASE")
		dbName = os.Getenv("MONGODB_TEST_DATABASE")
	}

	db, _ := database.Connect(dbName)

	return db
}

func clearDatabase(collection *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	collection.Drop(ctx)
	defer cancel()
}

func mockPlanet(mockPlanet models.Planet) (*models.Planet, *services.PlanetService) {
	db := loadDatabase()
	service := services.NewPlanetService(db)

	newPlanet, _ := service.Create(mockPlanet)

	return newPlanet, service
}

func TestNewPlanetService(t *testing.T) {
	db := loadDatabase()
	service := services.NewPlanetService(db)

	require.NotNil(t, service)
	require.NotNil(t, service.Collection)
}

func TestCreateNewPlanet(t *testing.T) {
	db := loadDatabase()
	service := services.NewPlanetService(db)

	mockPlanet := models.Planet{
		Name:    "Tatooine",
		Climate: "Arid",
		Terrain: "Desert",
	}

	newPlanet, err := service.Create(mockPlanet)

	require.NotNil(t, newPlanet.ID)
	require.NotNil(t, newPlanet.Appearances)
	require.Nil(t, err)

	clearDatabase(service.Collection)
}

func TestGetPlanetWithValidId(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Climate: "Arid",
		Terrain: "Desert",
	}

	mockedPlanet, service := mockPlanet(planet)
	id := mockedPlanet.ID.Hex()

	insertedPlanet, err := service.Get(id)

	require.Equal(t, mockedPlanet, insertedPlanet)
	require.Nil(t, err)
	clearDatabase(service.Collection)
}

func TestGetNonExistentPlanet(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Climate: "Arid",
		Terrain: "Desert",
	}

	_, service := mockPlanet(planet)
	id := primitive.NewObjectID().Hex()

	p, err := service.Get(id)

	require.Nil(t, p)
	require.Equal(t, mongo.ErrNoDocuments, err)

	clearDatabase(service.Collection)
}

func TestGetInvalidId(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Climate: "Arid",
		Terrain: "Desert",
	}

	_, service := mockPlanet(planet)
	id := "234567"

	p, err := service.Get(id)

	require.Nil(t, p)
	require.Equal(t, errors.New("the provided hex string is not a valid ObjectID"), err)

	clearDatabase(service.Collection)
}

func TestListPlanets(t *testing.T) {
	planet1 := models.Planet{
		Name:    "Tatooine",
		Climate: "Arid",
		Terrain: "Desert",
	}

	planet2 := models.Planet{
		Name:    "Tund",
		Climate: "unknown",
		Terrain: "barren, ash",
	}

	mock1, service := mockPlanet(planet1)
	mock2, _ := mockPlanet(planet2)

	planets := &[]models.Planet{*mock1, *mock2}

	res, err := service.Search(1, "")

	require.Nil(t, err)
	require.Equal(t, planets, res.Result)
	require.Equal(t, int64(2), res.Total)

	clearDatabase(service.Collection)
}

// func TestSearchPlanet(t *testing.T) {
// 	planet1 := models.Planet{
// 		Name:    "Tatooine",
// 		Climate: "Arid",
// 		Terrain: "Desert",
// 	}

// 	planet2 := models.Planet{
// 		Name:    "Tund",
// 		Climate: "unknown",
// 		Terrain: "barren, ash",
// 	}

// 	mockedPlanet, service := mockPlanet(planet1)
// 	mockPlanet(planet2)

// 	res, err := service.Search(1, mockedPlanet.Name)

// 	require.Nil(t, err)
// 	require.Equal(t, int64(1), res.Total)

// 	for _, p range res.Result {

// 	}

// 	require.Equal(t, int64(1), res.Total)

// 	clearDatabase(service.Collection)
// }

func TestDeletePlanet(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Climate: "Arid",
		Terrain: "Desert",
	}

	mockedPlanet, service := mockPlanet(planet)
	id := mockedPlanet.ID.Hex()

	res, err := service.Delete(id)

	require.Nil(t, err)
	require.Equal(t, "Planet was deleted successfully!", res)
}

func TestDeletePlanetNotFound(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Climate: "Arid",
		Terrain: "Desert",
	}

	_, service := mockPlanet(planet)
	id := primitive.NewObjectID().Hex()

	res, err := service.Delete(id)

	require.Equal(t, "", res)
	require.Equal(t, "no planet with this id was found in this so far far away galaxy", err.Error())

	clearDatabase(service.Collection)
}
