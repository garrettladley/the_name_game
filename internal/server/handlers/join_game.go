package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func JoinGame(c *fiber.Ctx, store *fsession.Store) error {
	gameID := c.Params("game_id")

	if gameID == "" {
		slog.Error("game_id empty")
		return c.SendStatus(http.StatusBadRequest)
	}

	game, ok := domain.GAMES.Get(domain.ID(gameID))
	if !ok {
		slog.Error("game not found", "game_id", gameID)
		return c.SendStatus(http.StatusNotFound)
	}

	playerID := domain.NewID()

	game.Join(playerID)

	err := session.SetInSession(c, store, "player_id", playerID.String(), session.SetExpiry(constants.EXPIRE_AFTER))
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	slog.Info("player joined game", "game_id", gameID, "player_id", playerID.String())

	return c.Redirect(fmt.Sprintf("/game/%s", game.ID), http.StatusSeeOther)
}
