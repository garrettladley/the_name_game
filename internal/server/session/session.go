package session

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/garrettladley/the_name_game/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

const (
	player_id_key string = "0"
)

func GetIDFromSession(c *fiber.Ctx, store *session.Store) (*domain.ID, error) {
	session, err := store.Get(c)
	if err != nil {
		slog.Error("failed to get session", "error", err)
		return nil, err
	}

	intfID := session.Get(player_id_key)
	if intfID == nil {
		return nil, fmt.Errorf("failed to get value from session with key %s", player_id_key)
	}

	id, ok := intfID.(domain.ID)
	if !ok {
		return nil, fmt.Errorf("failed to convert value from session with key %s to domain.ID", player_id_key)
	}

	return &id, nil
}

func SetIDInSession(c *fiber.Ctx, store *session.Store, value domain.ID, opts ...SessionSetterOpts) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}

	session.Set(player_id_key, value)

	for _, opt := range opts {
		opt(session)
	}

	return session.Save()
}

func DeleteIDFromSession(c *fiber.Ctx, store *session.Store) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}

	session.Delete(player_id_key)

	return session.Save()
}

type SessionSetterOpts func(*session.Session)

func SetExpiry(d time.Duration) SessionSetterOpts {
	return func(s *session.Session) {
		s.SetExpiry(d)
	}
}
