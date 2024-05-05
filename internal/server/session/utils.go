package session

import (
	"encoding/gob"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func GetFromSession[T any](c *fiber.Ctx, store *session.Store, key string) (*T, error) {
	fmt.Println("fctx", c)
	session, err := store.Get(c)
	if err != nil {
		return nil, err
	}

	value := session.Get(key)
	if value == nil {
		return nil, fmt.Errorf("failed to get value from session with key %s", key)
	}

	valueAsT, ok := value.(*T)
	if !ok {
		return nil, fmt.Errorf("failed to convert value to type %T", value)
	}

	return valueAsT, nil
}

type Gobbale interface {
	gob.GobEncoder
	gob.GobDecoder
}
type Gobbable interface {
	gob.GobEncoder
	gob.GobDecoder
}

// value must be built-in Go types due to the limitations of session.Set(string, any)
func SetInSession[T any](c *fiber.Ctx, store *session.Store, key string, value T, opts ...SessionSetterOpts) error {
	session, err := store.Get(c)
	if err != nil {
		return err
	}

	session.Set(key, value)

	for _, opt := range opts {
		opt(session)
	}

	return session.Save()
}

type SessionSetterOpts func(*session.Session)

func SetExpiry(d time.Duration) SessionSetterOpts {
	return func(s *session.Session) {
		s.SetExpiry(d)
	}
}
