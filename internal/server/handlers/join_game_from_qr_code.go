package handlers

import (
	"fmt"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func JoinGameFromQrCode(c *fiber.Ctx, store *fsession.Store) error {
	gameID := c.Params("game_id")
	if gameID == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	game, ok := domain.GAMES.Get(domain.ID(gameID))
	if !ok {
		return c.SendStatus(http.StatusNotFound)
	}

	playerID := domain.NewID()

	game.Join(playerID)

	if err := session.SetInSession(c, store, "player_id", playerID.String(), session.SetExpiry(constants.EXPIRE_AFTER)); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return hxRedirect(c, fmt.Sprintf("/game/%s", game.ID))
}
