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

	ctx, cancel := context.WithTimeoutCause(context.Background(), constants.EXPIRE_AFTER, errors.New("timeout"))
	defer cancel()

	done := make(chan struct{})
	errorCh := make(chan error, 1)
	go websockets.NewGameMessageHandler(conn, game, &player).HandleIncomingMessage(ctx, done, errorCh)

	select {
	case <-ctx.Done():
		slog.Error("context done", "error", ctx.Err())
	case err := <-errorCh:
		if err != nil {
			slog.Error("error handling incoming message", "error", err)
		}
	case <-done:
		slog.Info("done handling incoming message", "game_id", gameID, "player_id", playerID)
	}

	slog.Info("event loop ended", "game_id", gameID, "player_id", playerID)
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
