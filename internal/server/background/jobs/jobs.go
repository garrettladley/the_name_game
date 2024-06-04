package jobs

import (
	"github.com/garrettladley/the_name_game/internal/domain"
)

type Jobs struct {
	Games *domain.Games
}

func New(games *domain.Games) *Jobs {
	return &Jobs{Games: games}
}
