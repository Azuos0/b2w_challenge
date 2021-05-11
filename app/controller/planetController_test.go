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
