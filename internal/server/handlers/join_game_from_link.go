package handlers

import (
	"net/http"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/garrettladley/the_name_game/views/game"

	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func JoinGameFromLink(c *fiber.Ctx, store *fsession.Store) error {
	gameID, err := gameIDFromParams(c)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	g, err := domain.GAMES.Get(*gameID)
	if err != nil {
		return err
	}

	playerID := domain.NewID()

	g.Join(playerID)

	if err := session.SetID(c, store, playerID, session.SetExpiry(constants.EXPIRE_AFTER)); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return into(c, game.Index(true, *gameID))
}
