package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/garrettladley/the_name_game/internal/server/websockets"
	fiberws "github.com/gofiber/contrib/websocket"
)

func WSJoin(conn *fiberws.Conn) {
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
	mh := websockets.NewGameMessageHandler(conn, done, game, &player)
	go mh.HandleIncomingMessage()

	select {
	case <-done:
		slog.Info("game ended or player submitted name", "game_id", gameID, "player_id", playerID)
	case <-timeoutCtx.Done():
		slog.Error("game timed out", "game_id", gameID)
	}
}

func extractIDs(conn *fiberws.Conn) (domain.ID, domain.ID, error) {
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
