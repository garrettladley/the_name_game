package server

import (
	"net/http"

	go_json "github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/middleware/session"

	"github.com/garrettladley/the_name_game/internal/server/handlers"
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
	store := session.New()

	app.Get("/", handlers.Home)

	app.Post("/new_game", intoSessionedHandler(handlers.NewGame, store))
	app.Post("/game/:game_id", intoSessionedHandler(handlers.JoinGame, store))
	app.Get("/game/:game_id", intoSessionedHandler(handlers.Game, store))

	app.Get("/ws/:game_id/:player_id", websocket.New(handlers.WSJoin))

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(http.StatusNotFound).SendFile("./views/404.html")
	})
}

func intoSessionedHandler(handler func(c *fiber.Ctx, store *session.Store) error, store *session.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return handler(c, store)
	}
}
