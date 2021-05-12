package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Azuos0/b2w_challenge/app/models"
	"github.com/Azuos0/b2w_challenge/app/server"
	"github.com/Azuos0/b2w_challenge/app/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var app server.App

func init() {
	if os.Getenv("PORT") == "" {
		err := godotenv.Load("../../.env")

		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	uri := os.Getenv("MONGODB_TEST_DATABASE")
	app.InitializeApp(uri)
}

func clearDatabase() {
	planetService := services.NewPlanetService(app.DB)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	planetService.Collection.Drop(ctx)
	defer cancel()
}

func addMockPlanet(planet models.Planet) string {
	planetService := services.NewPlanetService(app.DB)
	mockedPlanet, _ := planetService.Create(planet)

	return mockedPlanet.ID.Hex()
}

func addMockPlanet2(planet models.Planet) *models.Planet {
	planetService := services.NewPlanetService(app.DB)
	mockedPlanet, _ := planetService.Create(planet)

	return mockedPlanet
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
}

func TestCreateValidPlanet(t *testing.T) {
	var jsonStr = []byte(`{
		"name": "Tatooine",
		"climate": "arid",
		"terrain": "desert"
	}`)

	req, _ := http.NewRequest("POST", "/api/planet", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	require.Equal(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, "Tatooine", m["name"])
	require.Equal(t, "arid", m["climate"])
	require.Equal(t, "desert", m["terrain"])
	require.NotNil(t, m["_id"])
	require.NotNil(t, m["appearances"])

	clearDatabase()
}

func TestCreateInvalidPlanet(t *testing.T) {
	var jsonStr = []byte(`{
		"name": "Tatooine",
		"climate": "arid"
	}`)

	req, _ := http.NewRequest("POST", "/api/planet", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)

	require.Equal(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, "terrain: Missing required field", m["error"])

	clearDatabase()
}

func TestGetExistentPlanet(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Terrain: "Desert",
		Climate: "Arid",
	}

	id := addMockPlanet(planet)
	url := fmt.Sprintf("/api/planet/%v", id)

	req, _ := http.NewRequest("GET", url, nil)
	response := executeRequest(req)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, id, m["_id"])

	clearDatabase()
}

func TestDeleteExistentPlanet(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Terrain: "Desert",
		Climate: "Arid",
	}

	id := addMockPlanet(planet)
	url := fmt.Sprintf("/api/planet/%v", id)

	req, _ := http.NewRequest("DELETE", url, nil)
	response := executeRequest(req)

	var m string
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "Planet was deleted successfully!", m)

	clearDatabase()
}

func TestGetNonExistentPlanet(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	url := fmt.Sprintf("/api/planet/%v", id)

	req, _ := http.NewRequest("GET", url, nil)
	response := executeRequest(req)

	require.Equal(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, "mongo: no documents in result", m["error"])
}

func TestDeleteNonExistentPlanet(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	url := fmt.Sprintf("/api/planet/%v", id)

	req, _ := http.NewRequest("DELETE", url, nil)
	response := executeRequest(req)

	require.Equal(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, "no planet with this id was found in this so far far away galaxy", m["error"])
}

func TestListPlanets(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Terrain: "Desert",
		Climate: "Arid",
	}

	planet2 := models.Planet{
		Name:    "Tund",
		Climate: "unknown",
		Terrain: "barren, ash",
	}

	mockedPlanet1 := addMockPlanet2(planet)
	mockedPlanet2 := addMockPlanet2(planet2)

	req, _ := http.NewRequest("GET", "/api/planets", nil)
	response := executeRequest(req)

	var m services.SearchResponse
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, int64(2), m.Total)
	require.Equal(t, int64(1), m.TotalPage)
	require.Contains(t, m.Result, *mockedPlanet1)
	require.Contains(t, m.Result, *mockedPlanet2)

	clearDatabase()
}

func TestSearchNonExistentPlanet(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Terrain: "Desert",
		Climate: "Arid",
	}
	addMockPlanet(planet)

	urlString := "/api/planets?name=Tholoth"

	req, _ := http.NewRequest("GET", urlString, nil)
	response := executeRequest(req)

	var m services.SearchResponse
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, http.StatusOK, response.Code)
	require.Empty(t, m.Result)

	clearDatabase()
}
func TestSearchExistentPlanet(t *testing.T) {
	planet1 := models.Planet{
		Name:    "Tatooine",
		Terrain: "Desert",
		Climate: "Arid",
	}
	planet2 := models.Planet{
		Name:    "Tholoth",
		Terrain: "Unknown",
		Climate: "Unknown",
	}

	addMockPlanet(planet1)
	mockedPlanet := addMockPlanet2(planet2)

	urlString := "/api/planets?name=Tholoth"

	req, _ := http.NewRequest("GET", urlString, nil)
	response := executeRequest(req)

	var m services.SearchResponse
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, http.StatusOK, response.Code)
	require.NotEmpty(t, m.Result)
	require.Contains(t, m.Result, *mockedPlanet)

	clearDatabase()
}

func TestListPagination(t *testing.T) {
	planet := models.Planet{
		Name:    "Tatooine",
		Terrain: "Desert",
		Climate: "Arid",
	}
	addMockPlanet(planet)

	urlString := "/api/planets?page=2"

	req, _ := http.NewRequest("GET", urlString, nil)
	response := executeRequest(req)

	var m services.SearchResponse
	json.Unmarshal(response.Body.Bytes(), &m)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, int64(2), m.Page)
	require.Empty(t, m.Result)

	clearDatabase()
}
