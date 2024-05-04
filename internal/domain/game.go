package domain

import (
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

const IDLength int = 7 // 3 runes 1 dash 3 runes

type ID string

type Player struct {
	Conn        *websocket.Conn
	PlayedID    ID
	IsSubmitted bool
	Name        *string
}

type Game struct {
	lock     sync.RWMutex
	ID       ID
	HostID   ID
	IsActive bool
	Conns    map[ID]Player
}

func NewGame(hostID ID) *Game {
	game := Game{
		lock:     sync.RWMutex{},
		ID:       NewID(),
		HostID:   hostID,
		IsActive: true,
		Conns:    make(map[ID]Player),
	}

	game.Conns[hostID] = Player{
		PlayedID:    hostID,
		IsSubmitted: false,
	}

	return &game
}

func (g *Game) Get(playerID ID) (Player, bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	player, ok := g.Conns[playerID]
	return player, ok
}

func (g *Game) Set(playerID ID, player Player) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.Conns[playerID] = player
}

func (g *Game) HandleSubmission(playerID ID, name string) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	if !g.IsActive {
		log.Infof("game %s is not active", g.ID)
		return nil
	}

	player, ok := g.Conns[playerID]
	if !ok {
		return fmt.Errorf("player with ID %s not found", playerID)
	}

	if player.IsSubmitted {
		log.Infof("player %s has already submitted a name", playerID)
		return nil
	}

	player.Name = &name
	player.IsSubmitted = true

	g.Conns[playerID] = player

	return nil
}

func (g *Game) ProcessGameOver(playerID ID) error {
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

var GAMES = NewGames()

type Games struct {
	lock  sync.RWMutex
	games map[ID]*Game
}

func NewGames() *Games {
	return &Games{
		lock:  sync.RWMutex{},
		games: make(map[ID]*Game),
	}
}

func (g *Games) New(game *Game) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.games[game.ID] = game
}

func (g *Games) Get(id ID) (*Game, bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	game, ok := g.games[id]
	return game, ok
}

func (g *Games) Delete(id ID) {
	g.lock.Lock()
	defer g.lock.Unlock()

	delete(g.games, id)
}
