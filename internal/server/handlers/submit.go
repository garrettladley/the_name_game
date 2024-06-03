package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/views/game"

	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func Submit(c *fiber.Ctx, store *fsession.Store) error {
	gameID := c.Params("game_id")

	session, err := store.Get(c)
	if err != nil {
		slog.Error("error getting session", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	playerID, ok := session.Get("player_id").(string)
	if !ok {
		slog.Error("player_id not found in session")
		return c.SendStatus(http.StatusInternalServerError)
	}

	if gameID == "" || playerID == "" {
		slog.Error("game_id or player_id empty", "game_id", gameID, "player_id", playerID)
		return c.SendStatus(http.StatusBadRequest)
	}

	g, ok := domain.GAMES.Get(domain.ID(gameID))
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
		return into(c, game.SubmitForm(domain.ID(gameID), g.IsHost(domain.ID(playerID)), params, errs))
	}

	if err := g.HandleSubmission(domain.ID(playerID), params.Name); err != nil {
		if errors.Is(err, domain.ErrUserAlreadySubmitted) {
			return c.SendStatus(http.StatusConflict)
		}
		slog.Error("error handling submission", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	if g.IsHost(domain.ID(playerID)) {
		return into(c, game.SubmitSuccess())
	}

	return hxRedirect(c, fmt.Sprintf("/game/%s/post", gameID))
}
