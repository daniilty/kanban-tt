package core

import (
	"context"
	"errors"

	"github.com/daniilty/kanban-tt/tasks/internal/pg"
)

type Status struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	OwnerID  string `json:"owner_id,omitempty"`
	ParentID int    `json:"parent_id,omitempty"`
	ChildID  int    `json:"child_id,omitempty"`
}

func (s *Status) toDB() *pg.Status {
	return &pg.Status{
		ID:       s.ID,
		Name:     s.Name,
		OwnerID:  s.OwnerID,
		ParentID: s.ParentID,
		ChildID:  s.ChildID,
	}
}

func (s *service) AddStatus(ctx context.Context, status *Status) (int, error) {
	exists, err := s.db.IsStatusWithNameExists(ctx, status.Name)
	if err != nil {
		return 0, err
	}

	if exists {
		return 0, ErrStatusWithNameExists
	}

	if status.ParentID == 0 {
		return s.db.AddParent(ctx, status.toDB())
	}

	id, err := s.db.AddChild(ctx, status.toDB())
	if err != nil && errors.Is(err, pg.ErrNoStatuses) {
		return 0, ErrNoSuchParent
	}

	return id, err
}

func (s *service) GetStatuses(ctx context.Context, uid string) ([]*Status, error) {
	statuses, err := s.db.GetStatuses(ctx, uid)
	if err != nil {
		return nil, err
	}

	return dbStatusesToView(statuses), nil
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
	if len(ss) == 0 {
		return []*Status{}
	}

	// root node
	var next *pg.Status = nil

	res := make([]*Status, 0, len(ss))
	mapped := make(map[int]*pg.Status, len(ss))

	for _, s := range ss {
		if s.ParentID == 0 {
			next = s
			break
		}

		mapped[s.ID] = s
	}

	// cannot build statuses if no root(wtf?)
	if next == nil {
		return []*Status{}
	}

	res = append(res, dbStatusToView(next))

	for next.ChildID != 0 {
		child, ok := mapped[next.ChildID]
		if !ok {
			break
		}

		res = append(res, dbStatusToView(child))
		next = child
	}

	return res
}

func dbStatusToView(status *pg.Status) *Status {
	return &Status{
		ID:   status.ID,
		Name: status.Name,
	}
}
