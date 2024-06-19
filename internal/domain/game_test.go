package domain

import (
	"slices"
	"testing"
)

func TestNextOnlyHost(t *testing.T) {
	t.Parallel()

	games := NewGames()

	hostID := NewID()
	game := NewGame(hostID)

	games.New(game)

	if err := game.HandleSubmission(hostID, "host"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if game.RemainingToSelect() != 1 {
		t.Errorf("expected 1 submission, got %d", game.remainingToSelect)
	}

	if game.Len() != 1 {
		t.Errorf("expected 1 player, got %d", game.Len())
	}

	if err := game.End(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	name, ok := game.Next()

	if !ok {
		t.Error("expected a name")
	}

	if *name != "host" {
		t.Errorf("expected host, got %s", *name)
	}

	if game.RemainingToSelect() != 0 {
		t.Errorf("expected 0 submissions, got %d", game.remainingToSelect)
	}

	if game.Len() != 0 {
		t.Errorf("expected 0 players, got %d", game.Len())
	}

	name, ok = game.Next()

	if ok {
		t.Error("expected no name")
	}

	if name != nil {
		t.Errorf("expected nil, got %s", *name)
	}
}

func TestNextHostAndOne(t *testing.T) {
	t.Parallel()

	games := NewGames()

	hostID := NewID()
	game := NewGame(hostID)

	games.New(game)

	if err := game.HandleSubmission(hostID, "host"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	playerID := NewID()
	game.Join(playerID)
	if err := game.HandleSubmission(playerID, "player"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	submissions := []string{"host", "player"}

	if game.RemainingToSelect() != 2 {
		t.Errorf("expected 2 submissions, got %d", game.remainingToSelect)
	}

	if game.Len() != 2 {
		t.Errorf("expected 2 players, got %d", game.Len())
	}

	if err := game.End(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	name, ok := game.Next()

	if !ok {
		t.Error("expected a name")
	}

	if name == nil {
		t.Fatal("expected a name")
	}

	if !slices.Contains(submissions, *name) {
		t.Errorf("expected host or player, got %s", *name)
	}

	for index, submission := range submissions {
		if submission == *name {
			submissions = append(submissions[:index], submissions[index+1:]...)
		}
	}

	if game.RemainingToSelect() != 1 {
		t.Errorf("expected 1 submissions, got %d", game.remainingToSelect)
	}

	if game.Len() != 1 {
		t.Errorf("expected 1 players, got %d", game.Len())
	}

	name, ok = game.Next()

	if !ok {
		t.Error("expected a name")
	}

	if name == nil {
		t.Fatal("expected a name")
	}

	remainingSubmission := submissions[0]

	if remainingSubmission != *name {
		t.Errorf("expected %s, got %s", remainingSubmission, *name)
	}

	if game.RemainingToSelect() != 0 {
		t.Errorf("expected 0 submissions, got %d", game.remainingToSelect)
	}

	if game.Len() != 0 {
		t.Errorf("expected 0 players, got %d", game.Len())
	}

	name, ok = game.Next()

	if ok {
		t.Error("expected no name")
	}

	if name != nil {
		t.Errorf("expected nil, got %s", *name)
	}
}
