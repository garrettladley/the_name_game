package domain

import (
	"fmt"

	go_json "github.com/goccy/go-json"
)

type Game struct {
	ID       ID            `json:"id"`
	HostID   ID            `json:"host_id"`
	IsActive bool          `json:"is_active"`
	Players  map[ID]Player `json:"players"`
}

type Player struct {
	ID            ID      `json:"id"`
	SubmittedName *string `json:"submitted_name"`
	Seen          bool    `json:"seen"`
}

func NewGame(hostID ID) (*Game, error) {
	host := Player{
		ID:            hostID,
		SubmittedName: nil,
	}

	game := &Game{
		ID:       NewID(),
		HostID:   hostID,
		IsActive: true,
		Players:  map[ID]Player{hostID: host},
	}

	if err := GAMES.Set(game); err != nil {
		return nil, err
	}

	return game, nil
}

func (g *Game) Get(playerID ID) (*Player, error) {
	player, ok := g.Players[playerID]
	if !ok {
		return nil, fmt.Errorf("player with id '%s' not found in game '%s'", playerID, g.ID)
	}

	return &player, nil
}

func (g *Game) Join(playerID ID) error {
	game, err := GAMES.Get(g.ID)
	if err != nil {
		return err
	}

	if _, ok := game.Players[playerID]; ok {
		return fmt.Errorf("player with id '%s' already in game '%s'", playerID, g.ID)
	}

	player := Player{
		ID:            playerID,
		SubmittedName: nil,
	}

	game.Players[playerID] = player

	if err := GAMES.Set(game); err != nil {
		return err
	}

	return nil
}

func (g *Game) IsHost(playerID ID) bool {
	return g.HostID == playerID
}

func (g *Game) Unseen() int {
	var unseen int
	for _, player := range g.Players {
		if !player.Seen {
			unseen++
		}
	}
	return unseen
}

func (g *Game) HandleSubmission(playerID ID, name string) error {
	player, err := g.Get(playerID)
	if err != nil {
		return err
	}

	if player.SubmittedName != nil {
		return ErrUserAlreadySubmitted
	}

	player.SubmittedName = &name
	g.Players[playerID] = *player

	if err := GAMES.Set(g); err != nil {
		return err
	}

	return nil
}

func (g *Game) End() error {
	g.IsActive = false

	if err := GAMES.Set(g); err != nil {
		return err
	}

	return nil
}

func (g *Game) Next() (*string, bool) {
	if g.Unseen() == 0 {
		return nil, false
	}

	var selectedPlayerID ID

	for playerID, player := range g.Players {
		if player.SubmittedName != nil {
			selectedPlayerID = playerID
			break
		}
	}

	player := g.Players[selectedPlayerID]
	name := *player.SubmittedName
	player.Seen = true
	g.Players[selectedPlayerID] = player

	if err := GAMES.Set(g); err != nil {
		return nil, false
	}

	return &name, true
}

func (g *Game) Slog() func() {
	return func() {
		// TODO
	}
}

func (g *Game) Marshal() (*string, error) {
	game, err := go_json.Marshal(g)
	if err != nil {
		return nil, err
	}
	strGame := string(game)
	return &strGame, nil
}

func UnmarshalGame(strGame string) (*Game, error) {
	var game Game
	err := go_json.Unmarshal([]byte(strGame), &game)
	if err != nil {
		return nil, err
	}
	return &game, nil
}
