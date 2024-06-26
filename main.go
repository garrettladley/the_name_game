package main

import (
	"context"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server"
	"github.com/garrettladley/the_name_game/internal/server/background"
	"github.com/garrettladley/the_name_game/internal/server/background/jobs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func main() {
	ctx := context.Background()

	app := server.Setup()

	static(app)

	backgroundJobs(ctx)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	slog.Info("Shutting down server")
	if err := app.Shutdown(); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}

	slog.Info("Server shutdown")
}

//go:embed public
var PublicFS embed.FS

//go:embed htmx
var HtmxFS embed.FS

//go:embed assets
var AssetsFS embed.FS

func static(app *fiber.App) {
	app.Use("/public", filesystem.New(filesystem.Config{
		Root:       http.FS(PublicFS),
		PathPrefix: "public",
		Browse:     true,
	}))
	app.Use("/htmx", filesystem.New(filesystem.Config{
		Root:       http.FS(HtmxFS),
		PathPrefix: "htmx",
		Browse:     true,
	}))
	app.Use("/", filesystem.New(filesystem.Config{
		Root:       http.FS(AssetsFS),
		PathPrefix: "assets",
		Browse:     true,
	}))
}

func backgroundJobs(ctx context.Context) {
	j := jobs.New(domain.GAMES)

	background.Go(j.GamesInfo(ctx))
}
