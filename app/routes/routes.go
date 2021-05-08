package routes

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azuos0/b2w_challenge/app/controller"
	"github.com/Azuos0/b2w_challenge/app/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeRouter(router *mux.Router, db *mongo.Database, ctx context.Context) {
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(map[string]string{
			"message": "Welcome to the Star Wars Planet App API 😉",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write(response)
	}).Methods("GET")

	planetClient := services.NewPlanetClient(ctx, db)
	inititializePlanetRoutes(router, planetClient)
}

func inititializePlanetRoutes(router *mux.Router, planetClient *services.PlanetClient) {
	router.HandleFunc("/api/planets", controller.Search(planetClient)).Methods("GET")
	router.HandleFunc("/api/planets", controller.CreatePlanet(planetClient)).Methods("POST")
	router.HandleFunc("/api/planets/{id}", controller.GetPlanet(planetClient)).Methods("GET")
	router.HandleFunc("/api/planets/{id}", controller.DeletePlanet(planetClient)).Methods("DELETE")
}
