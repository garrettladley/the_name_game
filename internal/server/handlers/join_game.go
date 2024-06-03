package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func JoinGame(c *fiber.Ctx, store *fsession.Store) error {
	params := game.JoinParams{
		GameID: c.FormValue("game_id"),
	}

	var errs game.JoinErrors
	gameID, err := domain.ParseID(params.GameID)
	if err != nil {
		errs.GameID = "Invalid Game Code"
		return into(c, game.JoinForm(params, errs))
	}

	g, ok := domain.GAMES.Get(*gameID)
	if !ok {
		errs.GameID = "Game not found"
		return into(c, game.JoinForm(params, errs))
	}

	playerID := domain.NewID()

	g.Join(playerID)

	if err := session.SetIDInSession(c, store, playerID); err != nil {
		slog.Error("failed to set player_id in session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return hxRedirect(c, fmt.Sprintf("/game/%s", g.ID))
}
