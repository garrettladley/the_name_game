package websockets

import (
	"log/slog"

	"github.com/fasthttp/websocket"
	go_json "github.com/goccy/go-json"
	fiberws "github.com/gofiber/contrib/websocket"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/protocol"
)

type GameMessageHandler struct {
	conn   *fiberws.Conn
	done   chan struct{}
	game   *domain.Game
	player *domain.Player
}

func NewGameMessageHandler(conn *fiberws.Conn, done chan struct{}, game *domain.Game, player *domain.Player) *GameMessageHandler {
	return &GameMessageHandler{
		conn:   conn,
		done:   done,
		game:   game,
		player: player,
	}
}

func (mh *GameMessageHandler) HandleIncomingMessage() {
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

func (mh *GameMessageHandler) handleIncomingMessage() error {
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

func (mh *GameMessageHandler) handleReadMessageError(err error) error {
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

func (mh *GameMessageHandler) handleNameSubmission(data []byte) error {
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

func (mh *GameMessageHandler) handleEndGame() error {
	if err := mh.game.ProcessGameInactive(mh.player.ID); err != nil {
		slog.Error("error processing game inactive", "error", err)
		return err
	} else {
		close(mh.done)
	}
	return nil
}
