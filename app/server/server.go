package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Azuos0/b2w_challenge/app/database"
	"github.com/Azuos0/b2w_challenge/app/routes"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Database
	Ctx    context.Context
}

func (app *App) InitializeApp() {
	var err error

	app.DB, app.Ctx, err = database.Connect(os.Getenv("MONGODB_DATABASE"))

	if err != nil {
		log.Println(err)
	}

	app.Router = mux.NewRouter()
	routes.InitializeRouter(app.Router, app.DB, app.Ctx)
}

func (app *App) Run(port string) {
	log.Fatal(http.ListenAndServe(port, app.Router))
}