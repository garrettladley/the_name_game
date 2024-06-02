package main

import (
	"log"
	"path/filepath"

	"github.com/garrettladley/the_name_game/internal/server"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := server.Setup()
	static(app)

	log.Fatal(app.Listen(":3000"))
}

func static(app *fiber.App) {
	app.Static("/public", filepath.Join(".", "public"), fiber.Static{
		Compress: true,
	})
	app.Static("/htmx", filepath.Join(".", "htmx"), fiber.Static{
		Compress: true,
	})
}
