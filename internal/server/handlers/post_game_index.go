package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/views/components"
	"github.com/gofiber/fiber/v2"
	fsession "github.com/gofiber/fiber/v2/middleware/session"
)

func PostGameIndex(c *fiber.Ctx, store *fsession.Store) error {
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

	if !g.IsHost(domain.ID(playerID)) {
		return c.SendStatus(http.StatusForbidden)
	}

	name, ok := g.Next()
	if !ok {
		slog.Error("no names left")
		return c.Redirect("/", http.StatusSeeOther)
	}

	var next string
	if g.Len() == 0 {
		// TODO fix html rendering here
		next = "/"
	} else {
		next = fmt.Sprintf("/game/%s/post", gameID)
	}

	return into(c, components.NameInfo(*name, next))
}
