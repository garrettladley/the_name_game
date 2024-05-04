package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
	"nhooyr.io/websocket"
)

func HandleJoinGame(w http.ResponseWriter, r *http.Request) error {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return fmt.Errorf("failed to accept websocket: %w", err)
	}
	defer c.CloseNow()

	timer := time.NewTimer(10 * time.Minute)
	l := rate.NewLimiter(rate.Every(time.Millisecond*100), 10)

	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timed out waiting for game to start")
		default:
			ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
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
	}
}
