package server

import (
	"log"
	"net/http"

	"github.com/Azuos0/b2w_challenge/app/controller"
	"github.com/Azuos0/b2w_challenge/app/database"
	"github.com/Azuos0/b2w_challenge/app/routes"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

func (app *App) InitializeApp(uri string) {
	var err error

	app.DB, err = database.Connect(uri)

	if err != nil {
		log.Println(err)
	}

	planetController := controller.PlanetController{}
	planetController.SetService(app.DB)

	app.Router = mux.NewRouter()
	routes.InitializeMainRouter(app.Router)
	routes.InititializePlanetRoutes(app.Router, &planetController)
}

func (app *App) Run(port string) {
	log.Printf("Server listening at port %v \n", port)
	log.Fatal(http.ListenAndServe(port, app.Router))
}
