package worker

import (
	"context"
	"time"

	"github.com/daniilty/kanban-tt/tasks/internal/core"
	"go.uber.org/zap"
)

type TasksCleaner struct {
	timeout time.Duration
	service core.Service
	logger  *zap.SugaredLogger
}

func NewTasksCleaner(timeout time.Duration, service core.Service, logger *zap.SugaredLogger) *TasksCleaner {
	return &TasksCleaner{
		timeout: timeout,
		service: service,
		logger:  logger,
	}
}

func (t *TasksCleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(t.timeout)

	t.logger.Info("Starting task cleaner.")
	for {
		err := t.service.DeleteExpiredTasks(ctx)
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
