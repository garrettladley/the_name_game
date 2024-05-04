package handlers

import (
	"fmt"
	"net/http"

	"github.com/garrettladley/the_name_game/internal/domain"
)

type CreateGameResponse struct {
	ID string `json:"id"`
}

func HandleNewGame(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid method %s", r.Method)
	}

	game := domain.NewGame()
	domain.GAMES.AddGame(game)

	http.Redirect(w, r, fmt.Sprintf("/game/%s", game.ID), http.StatusSeeOther)

	return nil
}
