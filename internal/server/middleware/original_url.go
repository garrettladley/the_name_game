package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type contextKey byte

const (
	OriginalURLKey contextKey = 0
)

func SetOriginalURL(c *fiber.Ctx) error {
	c.SetUserContext(context.WithValue(c.Context(), OriginalURLKey, c.OriginalURL()))

	return c.Next()
}
