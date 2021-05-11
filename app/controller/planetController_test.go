package controller_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Azuos0/b2w_challenge/app/server"
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

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)

	return rr
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
