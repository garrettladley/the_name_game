package handlers

import (
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
)

func Join(c *fiber.Ctx) error {
	return into(c, game.Join())
}
