package jobs

import (
	"context"
	"time"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/server/background"
)

func (j *Jobs) GamesInfo(ctx context.Context) background.JobFunc {
	return func() {
		t := time.NewTicker(constants.GAMES_INFO_INTERVAL)

		for range t.C {
			j.Games.Slog()()
		}
	}
}
