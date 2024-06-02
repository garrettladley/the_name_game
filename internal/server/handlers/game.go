package handlers

import (
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func Game(c *fiber.Ctx, store *session.Store) error {
	gameID := c.Params("game_id")

	session, err := store.Get(c)
	if err != nil {
		slog.Error("error getting session", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	playerID := session.Get("player_id")
	if playerID == nil {
		slog.Error("player_id not found in session")
		return c.SendStatus(http.StatusInternalServerError)
	}

	if gameID == "" || playerID == "" {
		slog.Error("game_id or player_id empty", "game_id", gameID, "player_id", playerID)
		return c.SendStatus(http.StatusBadRequest)
	}

	return into(c, game.Game(domain.ID(gameID), domain.ID(playerID.(string))))
}
