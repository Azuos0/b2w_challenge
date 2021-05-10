package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Azuos0/b2w_challenge/app/controller"
	"github.com/Azuos0/b2w_challenge/app/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeRouter(router *mux.Router, db *mongo.Database) {
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(map[string]string{
			"message": "Welcome to the Star Wars Planet App API 😉",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write(response)
	}).Methods("GET")

	inititializePlanetRoutes(router, db)
}

func inititializePlanetRoutes(router *mux.Router, db *mongo.Database) {
	planetService := services.NewPlanetService(db)
	router.HandleFunc("/api/planets", controller.Search(planetService)).Methods("GET")
	router.HandleFunc("/api/planet", controller.CreatePlanet(planetService)).Methods("POST")
	router.HandleFunc("/api/planet/{id}", controller.GetPlanet(planetService)).Methods("GET")
	router.HandleFunc("/api/planet/{id}", controller.DeletePlanet(planetService)).Methods("DELETE")
}
