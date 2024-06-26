package domain

import (
	"bytes"
	"compress/gzip"
	"io"
	"log/slog"
	"strings"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/gofiber/storage/memory/v2"
)

var GAMES = NewGames()

type Games struct {
	store *memory.Storage
}

func NewGames() *Games {
	return &Games{
		store: memory.New(),
	}
}

func (g *Games) New(game *Game) error {
	mGame, err := game.Marshal()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(*mGame)); err != nil {
		return err
	}

	if err := gz.Flush(); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	return g.store.Set(game.ID.String(), b.Bytes(), constants.EXPIRE_AFTER)
}

func (g *Games) Get(id ID) (*Game, error) {
	game, err := g.store.Get(id.String())
	if err != nil {
		return nil, err
	}

	r, err := gzip.NewReader(strings.NewReader(string(game)))
	if err != nil {
		return nil, err
	}
	s, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	umGame, err := UnmarshalGame(string(s))
	if err != nil {
		return nil, err
	}

	return umGame, nil
}

func (g *Games) Set(game *Game) error {
	mGame, err := game.Marshal()
	if err != nil {
		return err
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(*mGame)); err != nil {
		return err
	}

	if err := gz.Flush(); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	return g.store.Set(game.ID.String(), b.Bytes(), constants.EXPIRE_AFTER)
}

func (g *Games) Exists(id ID) bool {
	_, err := g.store.Get(id.String())
	return err == nil
}

func (g *Games) Slog() func() {
	return func() {
		keys, err := g.store.Keys()
		if err != nil {
			slog.Error("failed to get keys", "error", err)
			return
		}
		for _, key := range keys {
			game, err := g.Get(ID(key))
			if err != nil {
				slog.Error("failed to get game", "error", err)
				continue
			}
			game.Slog()()
		}
	}
}
