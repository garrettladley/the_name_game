package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func JoinGame(c *fiber.Ctx, store *fsession.Store) error {
	params := game.JoinParams{
		GameID: strings.TrimSpace(c.FormValue("game_id")),
	}

	var errs game.JoinErrors
	gameID, err := domain.ParseID(params.GameID)
	if err != nil {
		errs.GameID = "Invalid Game Code"
		return into(c, game.JoinForm(params, errs))
	}

	g, err := domain.GAMES.Get(*gameID)
	if err != nil {
		errs.GameID = "Error finding game"
		return into(c, game.JoinForm(params, errs))
	}

	playerID := domain.NewID()

	g.Join(playerID)

	if err := session.SetID(c, store, playerID, session.SetExpiry(constants.EXPIRE_AFTER)); err != nil {
		slog.Error("failed to set player_id in session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	// fixme overlap of join and submit a name
	return into(c, game.Index(false, *gameID))
}
