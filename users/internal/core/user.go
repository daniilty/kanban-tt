package core

import (
	"context"
	"errors"
	"time"

	"github.com/daniilty/kanban-tt/users/internal/pg"
	"github.com/daniilty/kanban-tt/users/internal/slice"
)

type User struct {
	ID             string
	Name           string
	Email          string
	EmailConfirmed bool
	PasswordHash   string
	TaskTTL        int
}

func (s *ServiceImpl) AddUser(ctx context.Context, user *User) (int, error) {
	u := user.toDB()
	now := time.Now()
	u.TaskTTL = int(s.GetDefaultTTL())
	u.CreatedAt = &now

	return s.db.AddUser(ctx, u)
}

func (s *ServiceImpl) GetUser(ctx context.Context, id string) (*User, bool, error) {
	exists, err := s.db.IsUserWithIDExists(ctx, id)
	if err != nil {
		return nil, false, err
	}

	if !exists {
		return nil, true, errors.New("user with such id does not exist")
	}

	user, err := s.db.GetUser(ctx, id)
	if err != nil {
		return nil, false, err
	}

	return convertDBUserToService(user), true, nil
}

func (s *ServiceImpl) GetUserTaskTTL(ctx context.Context, id string) (int, error) {
	ttl, err := s.db.GetUserTaskTTL(ctx, id)
	if err != nil && errors.Is(err, pg.ErrNoRows) {
		return 0, ErrNoSuchUser
	}

	return ttl, err
}

func (s *ServiceImpl) GetUserByEmail(ctx context.Context, email string) (*User, bool, error) {
	exists, err := s.db.IsUserWithEmailExists(ctx, email)
	if err != nil {
		return nil, false, err
	}

	if !exists {
		return nil, true, errors.New("user with such email does not exist")
	}

	user, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, false, err
	}

	return convertDBUserToService(user), true, nil
}

func (s *ServiceImpl) IsUserWithEmailExists(ctx context.Context, email string) (bool, error) {
	return s.db.IsUserWithEmailExists(ctx, email)
}

func (s *ServiceImpl) IsValidUserCredentials(ctx context.Context, email string, passwordHash string) (bool, error) {
	return s.db.IsUserWithEmailPasswordExists(ctx, email, passwordHash)
}

func (s *ServiceImpl) UpdateUser(ctx context.Context, user *User) error {
	if user.TaskTTL != 0 && !slice.Contains(s.GetTTLs(), int64(user.TaskTTL)) {
		return ErrNoSuchTTL
	}

	return s.db.UpdateUser(ctx, user.toDB())
}
