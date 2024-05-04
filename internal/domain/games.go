package domain

import "sync"

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
