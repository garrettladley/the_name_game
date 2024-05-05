package handlers

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func NewGame(c *fiber.Ctx, store *session.Store) error {
	hostID := domain.NewID()
	game := domain.NewGame(hostID)
	domain.GAMES.New(game)

	session, err := store.Get(c)
	if err != nil {
		slog.Error("error getting session", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	session.Set("player_id", hostID.String())

	session.SetExpiry(2 * time.Second)

	err = session.Save()
	if err != nil {
		slog.Error("error saving session", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Redirect(fmt.Sprintf("/game/%s", game.ID), fiber.StatusSeeOther)
}
