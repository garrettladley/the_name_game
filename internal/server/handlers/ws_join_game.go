package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/protocol"
	go_json "github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
)

func WSJoin(conn *websocket.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("error closing connection", "error", err)
		}
	}()

	gameID, playerID, err := extractIDs(conn)
	if err != nil {
		slog.Error("error extracting IDs", "error", err)
		return
	}

	game, ok := domain.GAMES.Get(gameID)
	if !ok {
		slog.Error("game not found", "game_id", gameID)
		return
	}

	if !game.IsActive {
		slog.Error("game is not active", "game_id", gameID)
		return
	}

	player, ok := game.Get(playerID)
	if !ok {
		slog.Error("player not found", "player_id", playerID)
		return
	}

	if player.IsSubmitted {
		slog.Error("player has already submitted name", "player_id", playerID)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), constants.EXPIRE_AFTER)
	defer cancel()

	slog.Info("validation succeeded, conn established, beginning event loop", "game_id", gameID, "player_id", playerID)

	done := make(chan struct{})
	go eventLoop(conn, done, game, &player)

	select {
	case <-done:
		slog.Info("game ended or player submitted name", "game_id", gameID, "player_id", playerID)
	case <-timeoutCtx.Done():
		slog.Error("game timed out", "game_id", gameID)
	}
}

func extractIDs(conn *websocket.Conn) (domain.ID, domain.ID, error) {
	gameID := domain.ID(conn.Params("game_id"))
	if gameID == "" {
		return "", "", errors.New("game id not provided")
	}

	playerID := domain.ID(conn.Params("player_id"))
	if playerID == "" {
		return "", "", errors.New("player id not provided")
	}

	return gameID, playerID, nil
}

func eventLoop(conn *websocket.Conn, done chan struct{}, game *domain.Game, player *domain.Player) {
	for {
		select {
		case <-done:
			return
		default:
			if err := handleIncomingMessage(conn, done, game, player); err != nil {
				return
			}
		}
	}
}

func handleIncomingMessage(conn *websocket.Conn, done chan struct{}, game *domain.Game, player *domain.Player) error {
	msgType, buff, err := conn.ReadMessage()
	if err != nil {
		return handleReadMessageError(err, game, player)
	}

	if msgType != websocket.TextMessage {
		slog.Error("unexpected message type", "message_type", msgType)
		return nil // continue processing
	}

	var msg protocol.Message
	if err := go_json.Unmarshal(buff, &msg); err != nil {
		slog.Error("error unmarshalling message", "error", err)
		return nil // continue processing
	}

	switch msg.Type {
	case protocol.MessageTypeSubmitName:
		return handleNameSubmission(done, msg.Data, game, player)
	case protocol.MessageTypeEndGame:
		return handleEndGame(done, game, player)
	default:
		slog.Error("unknown message type", "message_type", msg.Type)
		return nil // continue processing
	}
}

func handleReadMessageError(err error, game *domain.Game, player *domain.Player) error {
	switch err {
	case nil:
		return nil
	case websocket.ErrCloseSent:
		slog.Info("connection closed", "game_id", game.ID, "player_id", player.ID)
	default:
		slog.Error("error reading message", "error", err)
	}
	return err
}

func handleNameSubmission(done chan struct{}, data []byte, game *domain.Game, player *domain.Player) error {
	defer close(done)

	var submitName protocol.SubmitName
	if err := go_json.Unmarshal(data, &submitName); err != nil {
		slog.Error("error unmarshalling submit name", "error", err)
		return nil // continue processing
	}

	if err := game.HandleSubmission(player.ID, submitName.Name); err != nil {
		slog.Error("error handling submission", "error", err)
		return err
	}

	slog.Info("submission accepted", "game_id", game.ID, "player_id", player.ID, "name", submitName.Name)

	return nil
}

func handleEndGame(done chan struct{}, game *domain.Game, player *domain.Player) error {
	defer close(done)

	if err := game.ProcessGameInactive(player.ID); err != nil {
		slog.Error("error processing game inactive", "error", err)
		return err
	}

	return nil
}
