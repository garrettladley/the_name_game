package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/garrettladley/the_name_game/internal/constants"
	"github.com/garrettladley/the_name_game/internal/server/background"
)

func (j *Jobs) CleanGames(ctx context.Context) background.JobFunc {
	return func() {
		t := time.NewTicker(constants.CLEAN_INTERVAL)

		for range t.C {
			slog.Info("cleaning expired games")
			j.Games.DeleteExpired()
		}
	}
}
