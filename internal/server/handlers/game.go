package handlers

import (
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func Game(c *fiber.Ctx, store *fsession.Store) error {
	gameID, err := gameIDFromParams(c)
	if err != nil {
		slog.Error("bad game_id", "error", err)
		return c.SendStatus(http.StatusBadRequest)
	}

	playerID, err := session.GetIDFromSession(c, store)
	if err != nil {
		slog.Error("failed to get player_id from session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	g, ok := domain.GAMES.Get(*gameID)
	if !ok {
		slog.Error("game not found", "game_id", gameID)
		return c.SendStatus(http.StatusNotFound)
	}

	return into(c, game.Index(*gameID, g.HostID == *playerID))
}
