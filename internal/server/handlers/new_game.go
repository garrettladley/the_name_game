package handlers

import (
	"fmt"
	"log/slog"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/server/session"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func NewGame(c *fiber.Ctx, store *fsession.Store) error {
	hostID := domain.NewID()
	game := domain.NewGame(hostID)
	domain.GAMES.New(game)

	err := session.SetInSession(c, store, "player_id", hostID.String(), session.SetExpiry(constants.EXPIRE_AFTER))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	slog.Info("new game created", "game_id", game.ID.String(), "host_id", hostID.String())

	return c.Redirect(fmt.Sprintf("/game/%s", game.ID), fiber.StatusSeeOther)
}
