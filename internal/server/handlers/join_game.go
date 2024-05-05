package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/protocol"
	"github.com/gofiber/contrib/websocket"
)

func Join(conn *websocket.Conn) {
	defer conn.Close()

	rawGameID := conn.Params("game_id")

	if rawGameID == "" {
		slog.Error("game id not provided")
		return
	}

	gameID := domain.ID(rawGameID)

	rawPlayerID := conn.Params("plaer_id")

	if rawPlayerID == "" {
		slog.Error("player id not provided")
		return
	}

	playerID := domain.ID(rawPlayerID)

	game, ok := domain.GAMES.Get(gameID)
	if !ok {
		slog.Error("game not found", "game_id", gameID)
		return
	}

	player, ok := game.Get(playerID)
	if !ok {
		slog.Error("player not found", "player_id", playerID)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), constants.EXPIRE_AFTER)
	defer cancel()

	slog.Info("validation succeeded, conn established, beginning event loop", "game_id", gameID, "player_id", playerID)

	go eventLoop(timeoutCtx, conn, game, &player)

	select {
	case <-timeoutCtx.Done():
		slog.Error("game timed out", "game_id", gameID)
	}
}

func eventLoop(ctx context.Context, conn *websocket.Conn, game *domain.Game, player *domain.Player) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgType, buff, err := conn.ReadMessage()
			if err != nil {
				slog.Error("failed to read message", err)
				continue
			}

			if msgType != websocket.BinaryMessage {
				slog.Error("unexpected message type", "message_type", msgType)
				continue
			}

			var submitName protocol.SubmitName
			err = submitName.UnmarshallBinary(buff)

			err = game.HandleSubmission(player.ID, submitName.Name)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// TODO: check if the game is over?
			// TODO: allow host to end game
		}
	}
}
