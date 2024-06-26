package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

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

	playerID, err := session.GetID(c, store)
	if err != nil {
		slog.Error("failed to get player_id from session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	g, err := domain.GAMES.Get(*gameID)
	if err != nil {
		return err
	}

	if !g.IsHost(*playerID) {
		return c.SendStatus(http.StatusForbidden)
	}

	if g.Unseen() == 0 {
		if err := session.DeleteID(c, store); err != nil {
			slog.Error("failed to delete player_id from session", "error", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		if err := session.Destroy(c, store); err != nil {
			slog.Error("failed to destroy session", "error", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		return hxRedirect(c, "/")
	}
	name, _ := g.Next() // ignore error as we know there is a next name
	return into(c, game.NameInfo(*name, fmt.Sprintf("/game/%s/post", gameID)))
}
