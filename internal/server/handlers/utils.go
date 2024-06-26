package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func into(c *fiber.Ctx, component templ.Component) error {
	return adaptor.HTTPHandler(templ.Handler(component))(c)
}

func hxRedirect(c *fiber.Ctx, to string) error {
	if len(c.Get("HX-Request")) > 0 {
		c.Set("HX-Redirect", to)
		return c.SendStatus(http.StatusSeeOther)
	}
	return c.SendStatus(http.StatusSeeOther)
}

func gameIDFromParams(c *fiber.Ctx) (*domain.ID, error) {
	return domain.ParseID(c.Params("game_id"))
}
