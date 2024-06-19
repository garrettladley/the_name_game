package handlers

import (
	"encoding/base64"
	"fmt"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/views/game"
	"github.com/gofiber/fiber/v2"
	qrcode "github.com/skip2/go-qrcode"
)

func JoinGameQRCode(c *fiber.Ctx) error {
	gameID, err := gameIDFromParams(c)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if ok := domain.GAMES.Exists(*gameID); !ok {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var png []byte
	png, err = qrcode.Encode(fmt.Sprintf("https://%s/game/%s/join", c.Hostname(), gameID), qrcode.Medium, 256)
	if err != nil {
		return err
	}

	return into(c, game.JoinGameQR(base64.StdEncoding.EncodeToString(png)))
}
