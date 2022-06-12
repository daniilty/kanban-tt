package core

import (
	"context"
	"errors"
)

type User struct {
	ID             string
	Name           string
	Email          string
	EmailConfirmed bool
	PasswordHash   string
}

func (s *ServiceImpl) AddUser(ctx context.Context, user *User) error {
	return s.db.AddUser(ctx, user.toDB())
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
	return s.db.UpdateUser(ctx, user.toDB())
}
