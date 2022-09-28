package core

import (
	"context"

	"github.com/daniilty/kanban-tt/users/internal/pg"
)

var _ Service = (*ServiceImpl)(nil)

type Service interface {
	AddUser(context.Context, *User) (int, error)
	GetUser(context.Context, string) (*User, bool, error)
	GetUserByEmail(context.Context, string) (*User, bool, error)
	GetUserTaskTTL(context.Context, string) (int, error)
	GetDefaultTTL() int64
	GetTTLs() []int64
	IsUserWithEmailExists(context.Context, string) (bool, error)
	IsValidUserCredentials(context.Context, string, string) (bool, error)
	UpdateUser(context.Context, *User) error
	UnconfirmEmail(context.Context, string) error
}

type ServiceImpl struct {
	db pg.DB
}

func NewServiceImpl(db pg.DB) *ServiceImpl {
	return &ServiceImpl{
		db: db,
	}
}
