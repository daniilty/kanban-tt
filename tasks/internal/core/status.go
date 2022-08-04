package core

import (
	"context"
	"errors"
	"sort"

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
		ID:       s.ID,
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

	res := dbStatusesToView(statuses)
	sortStatuses(res)

	return res, nil
}

func (s *service) UpdateStatus(ctx context.Context, status *Status) (error, Code) {
	err := s.db.UpdateStatus(ctx, status.toDB())
	if err != nil {
		var code Code = CodeDBFail
		if errors.Is(err, pg.ErrEmptyModel) {
			code = CodeEmptyModel
		}

		return err, code
	}

	return nil, CodeOK
}

func (s *service) DeleteStatus(ctx context.Context, id int) error {
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

func sortStatuses(statuses []*Status) {
	sort.Slice(statuses, func(i, j int) bool {
		if statuses[i].Priority == statuses[j].Priority {
			return statuses[i].Name < statuses[j].Name
		}

		return statuses[i].Priority < statuses[j].Priority
	})
}
