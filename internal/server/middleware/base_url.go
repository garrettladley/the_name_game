package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	BaseURLKey contextKey = 0
)

func SetBaseURL(c *fiber.Ctx) error {
	var scheme string
	if c.Secure() {
		scheme = "https"
	} else {
		scheme = "http"
	}

	c.Locals(BaseURLKey, fmt.Sprintf("%s://%s", scheme, c.Hostname()))
	return c.Next()
}
