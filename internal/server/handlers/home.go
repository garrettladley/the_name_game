package handlers

import (
	"github.com/garrettladley/the_name_game/views/home"
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return into(c, home.Index())
}
