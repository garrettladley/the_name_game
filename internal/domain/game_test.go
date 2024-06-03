package domain

import "testing"

func TestNext(t *testing.T) {
	t.Parallel()

	hostID := NewID()
	game := NewGame(hostID)
	GAMES.New(game)

	game.HandleSubmission(hostID, "host")

	if game.SubmittedCount() != 1 {
		t.Errorf("expected 1 submission, got %d", game.submittedCount)
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

	if game.SubmittedCount() != 0 {
		t.Errorf("expected 0 submissions, got %d", game.submittedCount)
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
