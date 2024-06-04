package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func PostGame(c *fiber.Ctx, store *fsession.Store) error {
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

	var view templ.Component
	if g.IsHost(*playerID) {
		if g.SubmittedCount() == 0 {
			if err := session.DeleteIDFromSession(c, store); err != nil {
				slog.Error("failed to delete player_id from session", "error", err)
				return c.SendStatus(http.StatusInternalServerError)
			}
			return hxRedirect(c, "/")
		}
		name, _ := g.Next() // ignore error as we know there is a next name
		view = game.NameInfo(*name, fmt.Sprintf("/game/%s/post", gameID))
	} else {
		view = game.Post()
	}

	return into(c, view)
}
