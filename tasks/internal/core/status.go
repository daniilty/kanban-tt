package core

import (
	"context"
	"errors"
	"strconv"

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
	exists, err := s.db.IsStatusWithNameExists(ctx, status.Name, status.OwnerID)
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

	return buildStatusList(statuses)
}

func (s *service) UpdateStatusName(ctx context.Context, status *Status) (Code, error) {
	dbStatus, code, err := s.getDBStatusFor(ctx, status.ID, status.OwnerID)
	if err != nil {
		return code, err
	}

	if dbStatus.Name == status.Name {
		return CodeOK, nil
	}

	err = s.db.UpdateStatusName(ctx, status.ID, status.Name)
	if err != nil {
		return CodeDBFail, err
	}

	return CodeOK, nil
}

func (s *service) UpdateStatusParent(ctx context.Context, status *Status) (Code, error) {
	dbStatus, code, err := s.getDBStatusFor(ctx, status.ID, status.OwnerID)
	if err != nil {
		return code, err
	}

	if dbStatus.ParentID == status.ParentID {
		return CodeOK, nil
	}

	if status.ParentID != 0 {
		_, code, err = s.getDBStatusFor(ctx, status.ParentID, status.OwnerID)
		if err != nil {
			return code, err
		}
	}

	err = s.db.UpdateStatusParent(ctx, dbStatus, status.ParentID)
	if err != nil {
		return CodeDBFail, err
	}

	return CodeOK, nil
}

func (s *service) DeleteStatus(ctx context.Context, uid int, id int) (Code, error) {
	status, code, err := s.getDBStatusFor(ctx, id, strconv.Itoa(uid))
	if err != nil {
		return code, err
	}

	err = s.db.DeleteStatus(ctx, status)
	if err != nil {
		return CodeDBFail, err
	}

	return CodeOK, nil
}

func (s *service) getDBStatusFor(ctx context.Context, statusID int, ownerID string) (*pg.Status, Code, error) {
	status, err := s.db.GetStatus(ctx, statusID)
	if err != nil {
		if errors.Is(err, pg.ErrNoStatuses) {
			return nil, CodeNoSuchStatus, err
		}

		return nil, CodeDBFail, err
	}

	if status.OwnerID != ownerID {
		return nil, CodeNotPermitted, ErrNotPermitted
	}

	return status, CodeOK, nil
}

func buildStatusList(ss []*pg.Status) ([]*Status, error) {
	if len(ss) == 0 {
		return []*Status{}, nil
	}

	// root node
	var next *pg.Status = nil

	res := make([]*Status, 0, len(ss))
	mapped := make(map[int]*pg.Status, len(ss))

	for _, s := range ss {
		if s.ParentID == 0 {
			next = s
			continue
		}

		mapped[s.ID] = s
	}

	// cannot build statuses if no root(wtf?)
	if next == nil {
		return nil, errors.New("no root in linked list")
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

	return res, nil
}

func dbStatusToView(status *pg.Status) *Status {
	return &Status{
		ID:   status.ID,
		Name: status.Name,
	}
}
