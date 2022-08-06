package worker

import (
	"context"
	"time"

	"github.com/daniilty/kanban-tt/tasks/internal/pg"
	"go.uber.org/zap"
)

type TasksCleaner struct {
	timeout time.Duration
	db      pg.DB
	logger  *zap.SugaredLogger
}

func NewTasksCleaner(timeout time.Duration, db pg.DB, logger *zap.SugaredLogger) *TasksCleaner {
	return &TasksCleaner{
		timeout: timeout,
		db:      db,
		logger:  logger,
	}
}

func (t *TasksCleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(t.timeout)

	t.logger.Info("Starting task cleaner.")
	for {
		err := t.db.DeleteExpiredTasks(ctx)
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
