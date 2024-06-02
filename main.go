package main

import (
	"log"

	"github.com/garrettladley/the_name_game/internal/server"
)

func main() {
	app := server.Setup()

	log.Fatal(app.Listen(":3000"))
}
