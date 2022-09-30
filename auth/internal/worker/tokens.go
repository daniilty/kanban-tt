package worker

import (
	"context"
	"time"

	"github.com/daniilty/kanban-tt/auth/internal/pg"
	"go.uber.org/zap"
)

type TokensCleaner struct {
	timeout time.Duration
	db      pg.DB
	logger  *zap.SugaredLogger
}

func NewTokensCleaner(timeout time.Duration, db pg.DB, logger *zap.SugaredLogger) *TokensCleaner {
	return &TokensCleaner{
		timeout: timeout,
		db:      db,
		logger:  logger,
	}
}

func (t *TokensCleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(t.timeout)

	t.logger.Info("Starting tokens cleaner.")
	for {
		err := t.db.DeleteExpiredTokens(ctx)
		if err != nil {
			t.logger.Infow("DeleteExpiredTasks", "err", err)
		}

		select {
		case <-ctx.Done():
			t.logger.Info("Stopping task cleaner.")
			ticker.Stop()
			return
		case <-ticker.C:
		}
	}
}
