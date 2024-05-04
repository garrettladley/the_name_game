package server

import (
	go_json "github.com/goccy/go-json"

	"github.com/garrettladley/the_name_game/internal/server/handlers"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Setup() *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder: go_json.Marshal,
		JSONDecoder: go_json.Unmarshal,
	})

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}:${port} ${pid} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Static("/public", "./public")
	app.Static("/htmx", "./htmx")

	routes(app)

	return app
}

func routes(app *fiber.App) {
	app.Get("/", handlers.Home)

	app.Post("/new_game", handlers.NewGame)

	app.Get("/ws/:id", websocket.New(handlers.Join))

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendFile("./views/404.html")
	})
}
