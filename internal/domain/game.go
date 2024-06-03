package domain

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/garrettladley/the_name_game/internal/constants"
)

type Player struct {
	ID          ID
	IsSubmitted bool
	Name        *string
}

type Game struct {
	ID             ID
	HostID         ID
	IsActive       bool
	ExpiresAt      time.Time
	submittedCount int
	lock           sync.RWMutex
	players        map[ID]Player
}

func NewGame(hostID ID) *Game {
	game := Game{
		ID:        NewID(),
		HostID:    hostID,
		IsActive:  true,
		ExpiresAt: time.Now().Add(constants.EXPIRE_AFTER),
		lock:      sync.RWMutex{},
		players:   make(map[ID]Player),
	}

	game.players[hostID] = Player{
		ID:          hostID,
		IsSubmitted: false,
	}

	return &game
}

func (g *Game) Get(playerID ID) (Player, bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	player, ok := g.players[playerID]
	return player, ok
}

func (g *Game) Join(playerID ID) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.players[playerID] = Player{
		ID:          playerID,
		IsSubmitted: false,
	}
}

func (g *Game) Set(playerID ID, player Player) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.players[playerID] = player
}

func (g *Game) IsHost(playerID ID) bool {
	return g.HostID == playerID
}

func (g *Game) HandleSubmission(playerID ID, name string) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.IsActive {
		return ErrGameOver
	}

	player, ok := g.players[playerID]
	if !ok {
		return fmt.Errorf("player with ID %s not found", playerID)
	}

	if player.IsSubmitted {
		return ErrUserAlreadySubmitted
	}

	player.Name = &name
	player.IsSubmitted = true

	g.players[playerID] = player

	g.submittedCount++

	return nil
}

func (g *Game) End() error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.IsActive {
		return nil
	}

	g.IsActive = false

	return nil
}

func (g *Game) Len() int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return len(g.players)
}

func (g *Game) SubmittedCount() int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.submittedCount
}

func (g *Game) Next() (*string, bool) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.IsActive {
		return nil, false
	}

	if g.submittedCount == 0 {
		return nil, false
	}

	var selectedPlayer Player
	var selectedID ID

	for {
		keys := make([]ID, 0, len(g.players))
		for id := range g.players {
			keys = append(keys, id)
		}

		randomIndex := rand.Intn(len(keys))
		selectedID = keys[randomIndex]
		selectedPlayer = g.players[selectedID]

		if selectedPlayer.IsSubmitted {
			delete(g.players, selectedID)
			g.submittedCount--
			return selectedPlayer.Name, true
		}
	}
}
