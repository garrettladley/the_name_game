package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func JoinGame(c *fiber.Ctx, store *fsession.Store) error {
	gameID := c.Params("game_id")

	if gameID == "" {
		slog.Error("game_id empty")
		return c.SendStatus(http.StatusBadRequest)
	}

	params := game.JoinParams{
		GameID: c.FormValue("game_id"),
	}

	var errs game.JoinErrors
	if len(params.GameID) != domain.IDLength {
		errs.GameID = "Game Code must be of the form XXX-XXX"
		return into(c, game.JoinForm(params, errs))

	}

	g, ok := domain.GAMES.Get(domain.ID(gameID))
	if !ok {
		errs.GameID = "Game not found"
		return into(c, game.JoinForm(params, errs))
	}

	playerID := domain.NewID()

	g.Join(playerID)

	err := session.SetInSession(c, store, "player_id", playerID.String(), session.SetExpiry(constants.EXPIRE_AFTER))
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return hxRedirect(c, fmt.Sprintf("/game/%s", g.ID))
}
