package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Azuos0/b2w_challenge/app/models"
	"github.com/Azuos0/b2w_challenge/app/services"
	"github.com/Azuos0/b2w_challenge/app/utils"
	"github.com/gorilla/mux"
)

func CreatePlanet(planetService *services.PlanetClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		planet := models.Planet{}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &planet)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = planet.Validate()
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := planetService.Create(planet)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	}
}

func GetPlanet(planetService *services.PlanetClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		res, err := planetService.Get(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	}
}

func Search(planetService *services.PlanetClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("name")

		res, err := planetService.Search(query)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	}
}

func DeletePlanet(planetService *services.PlanetClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		res, err := planetService.Delete(id)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	}
}
