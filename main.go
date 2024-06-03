package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server"
	"github.com/garrettladley/the_name_game/internal/server/background"
	"github.com/garrettladley/the_name_game/internal/server/background/job"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := server.Setup()
	static(app)

	ctx := context.Background()

	background.Go(job.New(domain.GAMES).CleanGames(ctx))

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

func static(app *fiber.App) {
	app.Static("/public", filepath.Join(".", "public"), fiber.Static{
		Compress: true,
	})
	app.Static("/htmx", filepath.Join(".", "htmx"), fiber.Static{
		Compress: true,
	})
}
