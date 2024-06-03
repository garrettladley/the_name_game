package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/server/session"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func NewGame(c *fiber.Ctx, store *fsession.Store) error {
	hostID := domain.NewID()
	game := domain.NewGame(hostID)
	domain.GAMES.New(game)

	if err := session.SetIDInSession(c, store, hostID); err != nil {
		slog.Error("failed to set player_id in session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return hxRedirect(c, fmt.Sprintf("/game/%s", game.ID))
}
