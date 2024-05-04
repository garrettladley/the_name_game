package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/garrettladley/the_name_game/internal/domain"
	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

type Server struct {
	Logf func(f string, v ...interface{})
}

var SUBPROTOCOLS = []string{
	"echo",
	"new_game",
	"join_game",
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: SUBPROTOCOLS,
	})
	if err != nil {
		s.Logf("%v", err)
		return
	}
	defer c.CloseNow()

	if slices.Contains(SUBPROTOCOLS, c.Subprotocol()) {
		c.Close(websocket.StatusPolicyViolation, fmt.Sprintf("invalid subprotocol %q, client must support one of %v", c.Subprotocol(), SUBPROTOCOLS))
		return
	}

	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)
	for {
		err := handle(c.Subprotocol())(r.Context(), c, l)
		if err != nil {
			s.Logf("%v", err)
			return
		}
	}
}

func handle(protocol string) func(context.Context, *websocket.Conn, *rate.Limiter) error {
	switch protocol {
	case "new_game":
		return newGame
	case "join_game":
		return joinGame
	default:
		return echo
	}
}

func newGame(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	game := domain.NewGame()

	domain.GAMES.AddGame(game)

	return err
}

func joinGame(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	return nil
}

func echo(ctx context.Context, c *websocket.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("failed to io.Copy: %w", err)
	}

	err = w.Close()
	return err
}
