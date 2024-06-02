package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/views/components"
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func PostGame(c *fiber.Ctx, store *fsession.Store) error {
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

	var view templ.Component
	if g.IsHost(domain.ID(playerID)) {
		if g.SubmittedCount() == 0 {
			return c.Redirect("/", http.StatusSeeOther)
		}
		var next string
		if g.SubmittedCount() > 1 {
			next = fmt.Sprintf("/game/%s/post", gameID)
		} else {
			next = "/"
		}
		name, _ := g.Next() // can safely call due to check above
		view = components.NameInfo(*name, next)
	} else {
		view = game.Post()
	}

	return into(c, view)
}
