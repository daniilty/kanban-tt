package core

import (
	"context"

	"github.com/daniilty/kanban-tt/tasks/internal/pg"
)

type Status struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority uint32 `json:"priority"`
	OwnerID  string `json:"owner_id,omitempty"`
}

func (s *Status) toDB() *pg.Status {
	return &pg.Status{
		Name:     s.Name,
		Priority: int(s.Priority),
		OwnerID:  s.OwnerID,
	}
}

func (s *service) AddStatus(ctx context.Context, status *Status) error {
	exists, err := s.db.IsStatusWithNameExists(ctx, status.Name)
	if err != nil {
		return err
	}

	if exists {
		return ErrStatusWithNameExists
	}

	return s.db.AddStatus(ctx, status.toDB())
}

func (s *service) GetStatuses(ctx context.Context, uid string) ([]*Status, error) {
	statuses, err := s.db.GetStatuses(ctx, uid)
	if err != nil {
		return nil, err
	}

	return dbStatusesToView(statuses), nil
}

func (s *service) UpdateStatus(ctx context.Context, status *Status) error {
	return s.db.UpdateStatus(ctx, status.toDB())
}

func (s *service) DeleteStatus(ctx context.Context, id string) error {
	return s.db.DeleteStatus(ctx, id)
}

func dbStatusesToView(ss []*pg.Status) []*Status {
	res := make([]*Status, 0, len(ss))

	for i := range ss {
		res = append(res, dbStatusToView(ss[i]))
	}

	return res
}

func dbStatusToView(status *pg.Status) *Status {
	return &Status{
		ID:       status.ID,
		Name:     status.Name,
		Priority: uint32(status.Priority),
	}
}
