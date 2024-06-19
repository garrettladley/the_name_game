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

func JoinGameFromLink(c *fiber.Ctx, store *fsession.Store) error {
	gameID, err := gameIDFromParams(c)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	game, ok := domain.GAMES.Get(*gameID)
	if !ok {
		return c.SendStatus(http.StatusNotFound)
	}

	playerID := domain.NewID()

	game.Join(playerID)

	if err := session.SetIDInSession(c, store, playerID, session.SetExpiry(constants.EXPIRE_AFTER)); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Redirect(fmt.Sprintf("/game/%s", game.ID))
}
