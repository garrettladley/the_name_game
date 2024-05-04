package handlers

import (
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type CreateGameResponse struct {
	GameID domain.ID `json:"game_id"`
}

func NewGame(c *fiber.Ctx) error {
	game := domain.NewGame()
	domain.GAMES.AddGame(game)

	return c.Status(fiber.StatusCreated).JSON(CreateGameResponse{GameID: game.ID})
}
