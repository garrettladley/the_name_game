package websockets

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	return &GameMessageHandler{
		conn:   conn,
		game:   game,
		player: player,
	}
}

func (h *GameMessageHandler) Serve(ctx context.Context, done chan struct{}, errCh chan error) {
	go h.writePump()
	go h.readPump(ctx, done, errCh)
}

func (h *GameMessageHandler) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		h.conn.Close()
	}()

	for range ticker.C {
		h.conn.SetWriteDeadline(time.Now().Add(writeWait))
		if err := h.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			return
		}
	}
}

func (h *GameMessageHandler) readPump(ctx context.Context, done chan struct{}, errCh chan error) {
	defer close(done)

	for {
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		default:
			msgType, buff, err := h.conn.ReadMessage()

			slog.Info("message received", "msgType", msgType, "buff", string(buff), "err", err)

			if err != nil {
				if !fiberws.IsCloseError(err, websocket.CloseGoingAway) {
					errCh <- err
				}
				return
			}
			if err := h.handleIncomingMessage(msgType, buff); err != nil {
				errCh <- err
				return
			}
		}
	}
}

func (h *GameMessageHandler) handleIncomingMessage(msgType int, buff []byte) error {
	if msgType == websocket.TextMessage {
		var msg protocol.Message
		if err := go_json.Unmarshal(buff, &msg); err != nil {
			return fmt.Errorf("error unmarshalling message: %w", err)
		}

		switch msg.Type {
		case protocol.MessageTypeSubmitName:
			return h.handleNameSubmission(msg.Data)
		case protocol.MessageTypeEndGame:
			return h.handleEndGame()
		default:
			return fmt.Errorf("unknown message type: %d", msg.Type)
		}
	}
	return nil
}

func (h *GameMessageHandler) handleNameSubmission(data []byte) error {
	var submitName protocol.SubmitName
	if err := go_json.Unmarshal(data, &submitName); err != nil {
		return fmt.Errorf("error unmarshalling submit name: %w", err)
	}

	if err := h.game.HandleSubmission(h.player.ID, submitName.Name); err != nil {
		return fmt.Errorf("error handling submission: %w", err)
	}
	return nil
}

func (h *GameMessageHandler) handleEndGame() error {
	if err := h.game.ProcessGameInactive(h.player.ID); err != nil {
		return fmt.Errorf("error processing game inactive: %w", err)
	}
	return nil
}
