package main

import (
	"log"
	"os"

	"github.com/Azuos0/b2w_challenge/app/server"
	"github.com/joho/godotenv"
)

func main() {
	app := server.App{}

	if os.Getenv("PORT") == "" {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	app.InitializeApp()

	port := os.Getenv("PORT")
	app.Run(port)
}
