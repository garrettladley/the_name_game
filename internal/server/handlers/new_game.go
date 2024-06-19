package handlers

import (
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/garrettladley/the_name_game/views/game"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func NewGame(c *fiber.Ctx, store *fsession.Store) error {
	hostID := domain.NewID()
	g := domain.NewGame(hostID)
	domain.GAMES.New(g)

	if err := session.SetIDInSession(c, store, hostID, session.SetExpiry(constants.EXPIRE_AFTER)); err != nil {
		slog.Error("failed to set player_id in session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return into(c, game.Index(g.ID))
}
