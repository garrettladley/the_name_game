package server

import (
	go_json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fsession "github.com/gofiber/fiber/v2/middleware/session"

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

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	routes(app)

	return app
}

func routes(app *fiber.App) {
	store := fsession.New()

	app.Get("/", handlers.Home)

	app.Post("/new_game", intoSessionedHandler(handlers.NewGame, store))
	app.Get("/join", handlers.Join)
	app.Get("/game/:game_id/qr", handlers.JoinGameQrCode)
	app.Get("/game/:game_id/join", intoSessionedHandler(handlers.JoinGameFromQrCode, store))
	app.Post("/game/:game_id", intoSessionedHandler(handlers.JoinGame, store))
	app.Get("/game/:game_id", intoSessionedHandler(handlers.Game, store))
	app.Get("/game/:game_id/post", intoSessionedHandler(handlers.PostGame, store))
	app.Post("/game/:game_id/submit", intoSessionedHandler(handlers.Submit, store))
	app.Post("/game/:game_id/end", intoSessionedHandler(handlers.EndGame, store))
	app.Get("/game/:game_id/post", intoSessionedHandler(handlers.PostGame, store))
}

func intoSessionedHandler(handler func(c *fiber.Ctx, store *fsession.Store) error, store *fsession.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return handler(c, store)
	}
}
