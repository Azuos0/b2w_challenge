package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Azuos0/b2w_challenge/app/controller"
	"github.com/gorilla/mux"
)

func InitializeMainRouter(router *mux.Router) {
	router.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		response, _ := json.Marshal(map[string]string{
			"message": "Welcome to the Star Wars Planet App API 😉",
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		w.Write(response)
	}).Methods("GET")
}

func InititializePlanetRoutes(router *mux.Router, controller *controller.PlanetController) {
	router.HandleFunc("/api/planets", controller.Search()).Methods("GET")
	router.HandleFunc("/api/planet", controller.CreatePlanet()).Methods("POST")
	router.HandleFunc("/api/planet/{id}", controller.GetPlanet()).Methods("GET")
	router.HandleFunc("/api/planet/{id}", controller.DeletePlanet()).Methods("DELETE")
}
