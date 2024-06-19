package domain

import (
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"github.com/garrettladley/the_name_game/internal/constants"
)

type Player struct {
	ID           ID
	Name         *string
	beenSelected bool
}

func (p *Player) IsSubmitted() bool {
	return p.Name != nil
}

type Game struct {
	ID                ID
	HostID            ID
	IsActive          bool
	ExpiresAt         time.Time
	lock              sync.RWMutex
	players           map[ID]Player
	remainingToSelect int
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
		ID:           hostID,
		Name:         nil,
		beenSelected: false,
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
		ID: playerID,
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
	slog.Info("handling submission", "player_id", playerID, "name", name, "is_host", g.IsHost(playerID))
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.IsActive {
		return ErrGameOver
	}

	player, ok := g.players[playerID]
	if !ok {
		return fmt.Errorf("player with ID %s not found", playerID)
	}

	if player.IsSubmitted() {
		return ErrUserAlreadySubmitted
	}

	player.Name = &name
	g.players[playerID] = player

	slog.Info("after submission", "player", *g.players[playerID].Name)
	g.remainingToSelect++

	return nil
}

func (g *Game) End() error {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.IsActive = false

	return nil
}

func (g *Game) Len() int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return len(g.players)
}

func (g *Game) RemainingToSelect() int {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return g.remainingToSelect
}

func (g *Game) Next() (*string, bool) {
	g.lock.Lock()
	defer g.lock.Unlock()

	var all []string
	for _, player := range g.players {
		if player.IsSubmitted() {
			all = append(all, *player.Name)
		}
	}
	slog.Info("all names", "names", all)
	if g.IsActive {
		return nil, false
	}

	if g.remainingToSelect == 0 {
		return nil, false
	}

	var (
		selectedPlayer Player
		selectedID     ID
	)

	var available []ID
	for id, player := range g.players {
		if player.IsSubmitted() && !player.beenSelected {
			slog.Info("player available", "id", id, "name", *player.Name)
			available = append(available, id)
		}
	}

	if len(available) == 0 {
		return nil, false
	}

	randomIndex := rand.Intn(len(available))
	selectedID = available[randomIndex]
	selectedPlayer = g.players[selectedID]

	selectedPlayer.beenSelected = true
	g.players[selectedID] = selectedPlayer
	g.remainingToSelect--
	slog.Info("returning name", "name", *selectedPlayer.Name)
	return selectedPlayer.Name, true
}

func (g *Game) Slog() func() {
	g.lock.RLock()
	defer g.lock.RUnlock()

	return func() {
		slog.Info("game", "id", g.ID, "host", g.HostID, "active", g.IsActive, "expires_at", g.ExpiresAt, "players", g.players)
	}
}
