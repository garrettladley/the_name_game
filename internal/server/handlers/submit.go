package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

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

	playerID, err := session.GetID(c, store)
	if err != nil {
		slog.Error("failed to get player_id from session", "error", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	g, err := domain.GAMES.Get(*gameID)
	if err != nil {
		return err
	}

	if !g.IsActive {
		return c.SendStatus(http.StatusForbidden)
	}

	params := game.SubmitParams{
		Name: strings.TrimSpace(c.FormValue("name")),
	}

	var errs game.SubmitErrors
	ok := true
	if len(params.Name) < 2 {
		ok = false
		errs.Name = "Name must be at least 2 characters"
	}
	if len(params.Name) > 50 {
		ok = false
		errs.Name = "Name must be less than 50 characters"
	}

	if !ok {
		return into(c, game.SubmitForm(*gameID, params, errs))
	}

	slog.Info("handling submission", "game_id", gameID, "player_id", playerID, "name", params.Name)
	if err := g.HandleSubmission(*playerID, params.Name); err != nil {
		if errors.Is(err, domain.ErrUserAlreadySubmitted) {
			return c.SendStatus(http.StatusConflict)
		}
		slog.Error("error handling submission", "err", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	isHost := g.IsHost(*playerID)

	if !isHost {
		if err := session.Destroy(c, store); err != nil {
			slog.Error("failed to destroy session", "error", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
	}

	return into(c, game.SubmitSuccess(isHost, *gameID))
}
