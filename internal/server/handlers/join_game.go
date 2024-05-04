package handlers

import (
	"fmt"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/protocol"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

func Join(c *websocket.Conn) {
	defer c.Close()

	gameID := c.Params("id")
	hostID := c.Query("host_id", "")

	if gameID == "" && hostID == "" {
		return
	}

	game, ok := domain.GAMES.Get(domain.ID(gameID))
	if !ok {
		return
	}
	var player domain.Player
	if hostID == "" { // new player
		playerID := domain.NewID()
		player = domain.Player{
			Conn:        c,
			PlayedID:    playerID,
			IsSubmitted: false,
		}

		game.Set(playerID, player)
	} else { // host
		hostID := domain.ID(hostID)
		player, ok := game.Get(hostID)
		if !ok {
			log.Infof("host %s not found in game %s", hostID, gameID)
			return
		}

		player.Conn = c

		game.Set(hostID, player)
	}

	for {
		var submitName protocol.SubmitName
		err := c.ReadJSON(&submitName)
		if err != nil {
			log.Error(err)
			continue
		}

		err = game.HandleSubmission(submitName.PlayerID, submitName.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// TODO: check if the game is over?
		// TODO: allow host to end game

	}

	fmt.Println(game)
}
