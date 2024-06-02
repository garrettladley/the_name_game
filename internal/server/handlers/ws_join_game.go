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

	if player.Conn != nil {
		slog.Error("player already has connection", "player_id", playerID)
		return
	}

	if err := game.SetPlayerConn(playerID, conn); err != nil {
		slog.Error("error setting player connection", "error", err)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), constants.EXPIRE_AFTER)
	defer cancel()

	slog.Info("validation succeeded, conn established, beginning event loop", "game_id", gameID, "player_id", playerID)

	done := make(chan struct{})
	mh := NewMessageHandler(conn, done, game, &player)
	go mh.HandleIncomingMessage()

	select {
	case <-done:
		slog.Info("game ended or player submitted name", "game_id", gameID, "player_id", playerID)
	case <-timeoutCtx.Done():
		slog.Error("game timed out", "game_id", gameID)
	}
}

type MessageHandler struct {
	conn   *websocket.Conn
	done   chan struct{}
	game   *domain.Game
	player *domain.Player
}

func NewMessageHandler(conn *websocket.Conn, done chan struct{}, game *domain.Game, player *domain.Player) *MessageHandler {
	return &MessageHandler{
		conn:   conn,
		done:   done,
		game:   game,
		player: player,
	}
}

func (mh *MessageHandler) HandleIncomingMessage() {
	for {
		select {
		case <-mh.done:
			return
		default:
			if err := mh.handleIncomingMessage(); err != nil {
				return
			}
		}
	}
}

func (mh *MessageHandler) handleIncomingMessage() error {
	msgType, buff, err := mh.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			slog.Info("connection closed", "game_id", mh.game.ID, "player_id", mh.player.ID, "message", err.Error())
			close(mh.done)
			return nil
		}
		return mh.handleReadMessageError(err)
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
		return mh.handleNameSubmission(msg.Data)
	case protocol.MessageTypeEndGame:
		return mh.handleEndGame()

	default:
		slog.Error("unknown message type", "message_type", msg.Type)
		return nil // continue processing
	}
}

func (mh *MessageHandler) handleReadMessageError(err error) error {
	switch err {
	case nil:
		return nil
	case websocket.ErrCloseSent:
		slog.Info("connection closed", "game_id", mh.game.ID, "player_id", mh.player.ID)
	default:
		slog.Error("error reading message", "error", err)
	}
	return err
}

func (mh *MessageHandler) handleNameSubmission(data []byte) error {
	var submitName protocol.SubmitName
	if err := go_json.Unmarshal(data, &submitName); err != nil {
		slog.Error("error unmarshalling submit name", "error", err)
		return nil // continue processing
	}

	if err := mh.game.HandleSubmission(mh.player.ID, submitName.Name); err != nil {
		slog.Error("error handling submission", "error", err)
		return err
	} else {
		slog.Info("submission accepted", "game_id", mh.game.ID, "player_id", mh.player.ID, "name", submitName.Name)
		close(mh.done)
	}
	return nil
}

func (mh *MessageHandler) handleEndGame() error {
	if err := mh.game.ProcessGameInactive(mh.player.ID); err != nil {
		slog.Error("error processing game inactive", "error", err)
		return err
	} else {
		close(mh.done)
	}
	return nil
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
