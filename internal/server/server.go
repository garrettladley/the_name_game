package server

import (
	"net/http"
	"time"

	"github.com/garrettladley/the_name_game/internal/server/handlers"
	"github.com/garrettladley/the_name_game/internal/server/middleware"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fsession "github.com/gofiber/fiber/v2/middleware/session"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func Setup() *fiber.App {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}:${port} ${pid} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	staticPaths := map[string]struct{}{
		"/":                  {},
		"/htmx/htmx.min.js":  {},
		"/public/styles.css": {},
		"/site.webmanifest":  {},
		"/favicon-32x32.png": {},
		"/favicon-16x16.png": {},
		"/favicon.ico":       {},
	}

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			_, found := staticPaths[c.Path()]
			return !found
		},
		Expiration:   time.Hour * 24 * 365, // 1 year
		CacheControl: true,
	}))

	app.Use(middleware.SetBaseURL)

	routes(app)

	utility(app)

	return app
}

func routes(app *fiber.App) {
	store := fsession.New()

	app.Get("/", handlers.Home)

	app.Route("/game", func(r fiber.Router) {
		r.Post("new", intoSessionedHandler(handlers.NewGame, store))
		r.Get("/join", handlers.JoinGameForm)
		r.Post("/join", intoSessionedHandler(handlers.JoinGame, store))
		r.Route("/:game_id", func(r fiber.Router) {
			r.Get("/", intoSessionedHandler(handlers.Game, store))
			r.Get("/qr", handlers.JoinGameQRCode)
			r.Get("/join", intoSessionedHandler(handlers.JoinGameFromLink, store))
			r.Get("/post", intoSessionedHandler(handlers.PostGame, store))
			r.Post("/submit", intoSessionedHandler(handlers.Submit, store))
			r.Post("/end", intoSessionedHandler(handlers.EndGame, store))
		})
	})
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
