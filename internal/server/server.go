package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2/middleware/compress"
	fsession "github.com/gofiber/fiber/v2/middleware/session"

	"github.com/garrettladley/the_name_game/internal/server/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}:${port} ${pid} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	routes(app)

	utility(app)

	return app
}

func routes(app *fiber.App) {
	store := fsession.New()

	app.Get("/", handlers.Home)

	app.Post("/game/new", intoSessionedHandler(handlers.NewGame, store))
	app.Get("/game/join", handlers.JoinGameForm)
	app.Post("/game/join", intoSessionedHandler(handlers.JoinGame, store))
	app.Get("/game/:game_id", intoSessionedHandler(handlers.Game, store))
	app.Get("/game/:game_id/qr", handlers.JoinGameQrCode)
	app.Get("/game/:game_id/join", intoSessionedHandler(handlers.JoinGameFromQrCode, store))
	app.Get("/game/:game_id/post", intoSessionedHandler(handlers.PostGame, store))
	app.Post("/game/:game_id/submit", intoSessionedHandler(handlers.Submit, store))
	app.Post("/game/:game_id/end", intoSessionedHandler(handlers.EndGame, store))
	app.Get("/game/:game_id/post", intoSessionedHandler(handlers.PostGame, store))
}

func utility(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})
}

func intoSessionedHandler(handler func(c *fiber.Ctx, store *fsession.Store) error, store *fsession.Store) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return handler(c, store)
	}
}
