package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Azuos0/b2w_challenge/app/models"
	"github.com/Azuos0/b2w_challenge/app/services"
	"github.com/Azuos0/b2w_challenge/app/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlanetController struct {
	PlanetService *services.PlanetService
}

func (c *PlanetController) SetService(db *mongo.Database) {
	c.PlanetService = services.NewPlanetService(db)
}

func (controller *PlanetController) CreatePlanet() http.HandlerFunc {
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

		res, err := controller.PlanetService.Create(planet)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusCreated, res)
	}
}

func (controller *PlanetController) GetPlanet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		res, err := controller.PlanetService.Get(id)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				utils.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	}
}

func (controller *PlanetController) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		page := r.URL.Query().Get("page")

		var searchPage int64 = 1
		var err error
		
		if page != "" {
			searchPage, err = strconv.ParseInt(page, 10, 32)
			if err != nil {
				searchPage = 1
			}
		}

		res, err := controller.PlanetService.Search(searchPage, name)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	}
}

func (controller *PlanetController) DeletePlanet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]

		res, err := controller.PlanetService.Delete(id)
		if err != nil {
			if err.Error() == "no planet with this id was found in this so far far away galaxy" {
				utils.RespondWithError(w, http.StatusNotFound, err.Error())
				return
			}
			utils.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, res)
	}
}
