package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/garrettladley/the_name_game/views/components"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func EndGame(c *fiber.Ctx, store *fsession.Store) error {
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

	if !g.IsHost(*playerID) {
		return c.SendStatus(http.StatusForbidden)
	}

	if err := g.End(); err != nil {
		slog.Error("error ending game", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	name, ok := g.Next()
	if !ok {
		if err := session.DeleteIDFromSession(c, store); err != nil {
			slog.Error("failed to delete player_id from session", "error", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		return hxRedirect(c, "/")
	}

	return into(c, components.NameInfo(*name, fmt.Sprintf("/game/%s/post", gameID)))
}
