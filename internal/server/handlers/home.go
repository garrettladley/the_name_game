package handlers

import (
	"github.com/a-h/templ"
	"github.com/garrettladley/the_name_game/views/home"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Home(c *fiber.Ctx) error {
	return adaptor.HTTPHandler(templ.Handler(home.Index()))(c)
}
