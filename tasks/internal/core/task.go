package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/daniilty/kanban-tt/tasks/internal/pg"
)

type Task struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Priority uint32 `json:"priority"`
	OwnerID  string `json:"owner_id,omitempty"`
	StatusID uint32 `json:"status_id"`
}

func (t *Task) toDB() *pg.Task {
	now := time.Now()

	return &pg.Task{
		Content:   t.Content,
		Priority:  int(t.Priority),
		OwnerID:   t.OwnerID,
		StatusID:  int(t.StatusID),
		CreatedAt: &now,
	}
}

func (s *service) AddTask(ctx context.Context, t *Task) (error, bool) {
	const maxPriority = 4

	if t.Priority > maxPriority {
		return fmt.Errorf("invalid priority: cannot be bigger than %d", maxPriority), true
	}

	status, err := s.db.GetStatusWithLowestPriority(ctx, t.OwnerID)
	if err != nil {
		return err, errors.Is(err, pg.ErrNoStatuses)
	}
	t.StatusID = uint32(status.ID)

	err = s.db.AddTask(ctx, t.toDB())
	if err != nil {
		return err, false
	}

	return nil, true
}

func (s *service) GetTasks(ctx context.Context, uid string) ([]*Task, error) {
	tasks, err := s.db.GetTasks(ctx, uid)
	if err != nil {
		return nil, err
	}

	return dbTasksToView(tasks), nil
}

func (s *service) UpdateTask(ctx context.Context, t *Task) error {
	return s.db.UpdateTask(ctx, t.toDB())
}

func (s *service) DeleteTask(ctx context.Context, id string) error {
	return s.db.DeleteTask(ctx, id)
}

func dbTasksToView(tt []*pg.Task) []*Task {
	res := make([]*Task, 0, len(tt))

	for i := range tt {
		res = append(res, dbTaskToView(tt[i]))
	}

	return res
}

func dbTaskToView(t *pg.Task) *Task {
	return &Task{
		ID:       t.ID,
		Content:  t.Content,
		Priority: uint32(t.Priority),
		StatusID: uint32(t.StatusID),
	}
}
