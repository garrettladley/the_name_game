package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

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
		return c.SendStatus(http.StatusInternalServerError)
	}

	slog.Info("new game created", "game_id", game.ID.String(), "host_id", hostID.String())
	slog.Info("join info", "see", fmt.Sprintf("/%s/%s", game.ID, hostID))

	return c.Redirect(fmt.Sprintf("/game/%s", game.ID), http.StatusSeeOther)
}
