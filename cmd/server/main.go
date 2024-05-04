package main

import (
	"log"

	"github.com/garrettladley/the_name_game/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	app := server.Setup()

	log.Fatal(app.Listen(":3000"))
}
