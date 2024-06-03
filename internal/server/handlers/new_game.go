package handlers

import (
	"fmt"
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

	for i := 0; i < 5; i++ {
		playerID := domain.NewID()
		game.Join(playerID)
		game.HandleSubmission(playerID, fmt.Sprintf("player_%d", i))
	}

	if err := session.SetInSession(c, store, "player_id", hostID.String(), session.SetExpiry(constants.EXPIRE_AFTER)); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return hxRedirect(c, fmt.Sprintf("/game/%s", game.ID))
}
