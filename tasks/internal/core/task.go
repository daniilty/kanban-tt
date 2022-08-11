package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/daniilty/kanban-tt/schema"
	"github.com/daniilty/kanban-tt/tasks/internal/pg"
)

const (
	// неверное значения приоритета
	CodeInvalidPriority Code = "INVALID_PRIORITY"
	// у нас нет статусов в бд
	CodeNoStatuses Code = "NO_STATUSES"
	// у нас нет статуса с таким айди в бд
	CodeNoSuchStatus Code = "NO_SUCH_STATUS"
	// пришла пустая структура лол
	CodeEmptyModel Code = "EMPTY_MODEL"
)

type Task struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Priority uint32 `json:"priority"`
	OwnerID  string `json:"owner_id,omitempty"`
	StatusID uint32 `json:"status_id"`
}

func (t *Task) toDB() *pg.Task {
	return &pg.Task{
		ID:       t.ID,
		Content:  t.Content,
		Priority: int(t.Priority),
		OwnerID:  t.OwnerID,
		StatusID: int(t.StatusID),
	}
}

func (s *service) AddTask(ctx context.Context, t *Task) (error, Code) {
	const maxPriority = 4

	if t.Priority > maxPriority {
		return fmt.Errorf("invalid priority: cannot be bigger than %d", maxPriority), CodeInvalidPriority
	}

	status, err := s.db.GetStatusWithLowestPriority(ctx, t.OwnerID)
	if err != nil {
		if errors.Is(err, pg.ErrNoStatuses) {
			return err, CodeNoStatuses
		}

		return err, CodeDBFail
	}

	t.StatusID = uint32(status.ID)

	tDB := t.toDB()
	now := time.Now()
	tDB.CreatedAt = &now

	err = s.db.AddTask(ctx, tDB)
	if err != nil {
		return err, CodeDBFail
	}

	return nil, CodeOK
}

func (s *service) GetUserTasks(ctx context.Context, uid string) ([]*Task, error) {
	tasks, err := s.db.GetUserTasks(ctx, uid)
	if err != nil {
		return nil, err
	}

	return dbTasksToView(tasks), nil
}

func (s *service) UpdateTask(ctx context.Context, t *Task) (error, Code) {
	exists, err := s.db.IsStatusWithIDExists(ctx, int(t.StatusID))
	if err != nil {
		return err, CodeDBFail
	}

	if !exists {
		return fmt.Errorf("status with such id does not exist: %d", t.StatusID), CodeNoSuchStatus
	}

	err = s.db.UpdateTask(ctx, t.toDB())
	if err != nil {
		if errors.Is(err, pg.ErrEmptyModel) {
			return err, CodeEmptyModel
		}

		return err, CodeDBFail
	}

	return nil, CodeOK
}

func (s *service) DeleteTask(ctx context.Context, id int) error {
	return s.db.DeleteTask(ctx, id)
}

func (s *service) DeleteExpiredTasks(ctx context.Context) error {
	userChecked := map[string]struct{}{}

	tasks, err := s.db.GetTasks(ctx)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		_, ok := userChecked[t.OwnerID]
		if ok {
			continue
		}

		resp, err := s.userService.GetUserTaskTTL(ctx, &schema.GetUserTaskTTLRequest{
			Id: t.OwnerID,
		})
		if err != nil {
			return err
		}

		ttl := resp.GetTaskTtl()

		err = s.db.DeleteExpiredTasks(ctx, t.OwnerID, int(ttl))
		if err != nil {
			return err
		}
	}

	return nil
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
