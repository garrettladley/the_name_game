package handlers

import (
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type CreateGameResponse struct {
	HostID domain.ID `json:"host_id"`
	GameID domain.ID `json:"game_id"`
}

func NewGame(c *fiber.Ctx) error {
	hostID := domain.NewID()
	game := domain.NewGame(hostID)
	domain.GAMES.New(game)

	return c.Status(fiber.StatusCreated).JSON(
		CreateGameResponse{
			HostID: hostID,
			GameID: game.ID,
		},
	)
}
