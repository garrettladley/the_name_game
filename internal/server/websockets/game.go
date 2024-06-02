package websockets

import (
	"context"
	"fmt"

	"github.com/fasthttp/websocket"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/protocol"
	go_json "github.com/goccy/go-json"
	fiberws "github.com/gofiber/contrib/websocket"
)

type GameMessageHandler struct {
	conn   *fiberws.Conn
	game   *domain.Game
	player *domain.Player
}

func NewGameMessageHandler(conn *fiberws.Conn, game *domain.Game, player *domain.Player) *GameMessageHandler {
	return &GameMessageHandler{
		conn:   conn,
		game:   game,
		player: player,
	}
}

func (mh *GameMessageHandler) HandleIncomingMessage(ctx context.Context, done chan struct{}, errCh chan error) {
	defer close(done)

	for {
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		default:
			if err := mh.handleIncomingMessage(); err != nil {
				errCh <- err
				return
			}
		}
	}
}

func (mh *GameMessageHandler) handleIncomingMessage() error {
	msgType, buff, err := mh.conn.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			// connection closed, return gracefully
			return nil
		}
		return fmt.Errorf("error reading message: %w", err)
	}

	if msgType == websocket.TextMessage {
		var msg protocol.Message
		if err := go_json.Unmarshal(buff, &msg); err != nil {
			return fmt.Errorf("error unmarshalling message: %w", err)
		}

		switch msg.Type {
		case protocol.MessageTypeSubmitName:
			if err := mh.handleNameSubmission(msg.Data); err != nil {
				return err
			}
		case protocol.MessageTypeEndGame:
			if err := mh.handleEndGame(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown message type: %d", msg.Type)
		}
	}
	return nil
}

func (mh *GameMessageHandler) handleNameSubmission(data []byte) error {
	var submitName protocol.SubmitName
	if err := go_json.Unmarshal(data, &submitName); err != nil {
		return fmt.Errorf("error unmarshalling submit name: %w", err)
	}

	if err := mh.game.HandleSubmission(mh.player.ID, submitName.Name); err != nil {
		return fmt.Errorf("error handling submission: %w", err)
	}
	return nil
}

func (mh *GameMessageHandler) handleEndGame() error {
	if err := mh.game.ProcessGameInactive(mh.player.ID); err != nil {
		return fmt.Errorf("error processing game inactive: %w", err)
	}
	return nil
}
