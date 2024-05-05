package domain

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

type Player struct {
	Conn        *websocket.Conn
	ID          ID
	IsSubmitted bool
	Name        *string
}

type Game struct {
	ID       ID
	HostID   ID
	IsActive bool
	lock     sync.RWMutex
	conns    map[ID]Player
}

func NewGame(hostID ID) *Game {
	game := Game{
		ID:       NewID(),
		HostID:   hostID,
		IsActive: true,
		lock:     sync.RWMutex{},
		conns:    make(map[ID]Player),
	}

	game.conns[hostID] = Player{
		ID:          hostID,
		IsSubmitted: false,
	}

	return &game
}

func (g *Game) Get(playerID ID) (Player, bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	player, ok := g.conns[playerID]
	return player, ok
}

func (g *Game) Set(playerID ID, player Player) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.conns[playerID] = player
}

func (g *Game) IsHost(playerID ID) bool {
	return g.HostID == playerID
}

func (g *Game) HandleSubmission(playerID ID, name string) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.IsActive {
		log.Infof("game %s is not active", g.ID)
		return nil
	}

	player, ok := g.conns[playerID]
	if !ok {
		return fmt.Errorf("player with ID %s not found", playerID)
	}

	if player.IsSubmitted {
		log.Infof("player %s has already submitted a name", playerID)
		return nil
	}

	player.Name = &name
	player.IsSubmitted = true

	g.conns[playerID] = player

	return nil
}

func (g *Game) ProcessGameInactive(playerID ID) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.IsActive {
		return nil
	}

	if playerID != g.HostID {
		return fmt.Errorf("player %s is not the host", playerID)
	}

	g.IsActive = false

	return nil
}
