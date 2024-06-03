package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/session"

	"github.com/garrettladley/the_name_game/views/game"

	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func Submit(c *fiber.Ctx, store *fsession.Store) error {
	gameID, err := gameIDFromParams(c)
	if err != nil {
		slog.Error("bad game_id", "error", err)
		return c.SendStatus(http.StatusBadRequest)
	}

	playerID, err := session.GetIDFromSession(c, store)
	if err != nil {
		slog.Error("failed to get player_id from session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	g, ok := domain.GAMES.Get(*gameID)
	if !ok {
		slog.Error("game not found", "game_id", gameID)
		return c.SendStatus(http.StatusNotFound)
	}

	if !g.IsActive {
		return c.SendStatus(http.StatusForbidden)
	}

	params := game.SubmitParams{
		Name: c.FormValue("name"),
	}

	var errs game.SubmitErrors
	ok = true
	if len(params.Name) < 2 {
		ok = false
		errs.Name = "Name must be at least 2 characters"
	}
	if len(params.Name) > 50 {
		ok = false
		errs.Name = "Name must be less than 50 characters"
	}

	if !ok {
		return into(c, game.SubmitForm(*gameID, g.IsHost(*playerID), params, errs))
	}

	if err := g.HandleSubmission(*playerID, params.Name); err != nil {
		if errors.Is(err, domain.ErrUserAlreadySubmitted) {
			return c.SendStatus(http.StatusConflict)
		}
		slog.Error("error handling submission", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	if g.IsHost(*playerID) {
		return into(c, game.SubmitSuccess())
	}

	if err := session.DeleteIDFromSession(c, store); err != nil {
		slog.Error("failed to delete player_id from session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return hxRedirect(c, fmt.Sprintf("/game/%s/post", gameID))
}
