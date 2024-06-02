package domain

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/gofiber/contrib/websocket"
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

func (g *Game) Join(playerID ID) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.conns[playerID] = Player{
		ID:          playerID,
		IsSubmitted: false,
	}
}

func (g *Game) SetPlayerConn(playerID ID, conn *websocket.Conn) error {
	g.lock.Lock()
	defer g.lock.Unlock()

	player, ok := g.conns[playerID]
	if !ok {
		return fmt.Errorf("player with ID %s not found", playerID)
	}

	player.Conn = conn
	g.conns[playerID] = player

	return nil
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
		return ErrGameOver
	}

	player, ok := g.conns[playerID]
	if !ok {
		return fmt.Errorf("player with ID %s not found", playerID)
	}

	if player.IsSubmitted {
		return ErrUserAlreadySubmitted
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

	for _, player := range g.conns {
		if player.ID != g.HostID && player.Conn != nil {
			if err := player.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "game ended")); err != nil {
				slog.Error("error writing close message", "error", err)
			}

			if err := player.Conn.Close(); err != nil {
				slog.Error("error closing connection", "error", err)
			}
		}
	}

	return nil
}
