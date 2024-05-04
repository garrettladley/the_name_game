package domain

import (
	"sync"

	"nhooyr.io/websocket"
)

const IDLength int = 7 // 3 runes 1 dash 3 runes

type ID string

type Player struct {
	Conn     *websocket.Conn
	PlayedID ID
	Name     string
}

type Game struct {
	ID    ID
	Host  Player
	Conns map[ID]Player
}

func NewGame() *Game {
	return &Game{
		ID:    NewID(),
		Conns: make(map[ID]Player),
	}
}

var GAMES = NewGames()

type Games struct {
	Games sync.Map // K: ID, V: *Game
}

func NewGames() *Games {
	return &Games{}
}

func (g *Games) AddGame(game *Game) {
	g.Games.Store(game.ID, game)
}

func (g *Games) GetGame(id ID) (*Game, bool) {
	game, ok := g.Games.Load(id)
	if !ok {
		return nil, false
	}
	return game.(*Game), true
}

func (g *Games) DeleteGame(id ID) {
	g.Games.Delete(id)
}
