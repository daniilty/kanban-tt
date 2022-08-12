package core

import (
	"context"

	"github.com/daniilty/kanban-tt/schema"
	"github.com/daniilty/kanban-tt/tasks/internal/pg"
)

type Service interface {
	AddTask(context.Context, *Task) (int, error, Code)
	GetUserTasks(context.Context, string) ([]*Task, error)
	UpdateTask(context.Context, *Task) (error, Code)
	DeleteTask(context.Context, int) error
	DeleteExpiredTasks(context.Context) error

	AddStatus(context.Context, *Status) (int, error)
	GetStatuses(context.Context, string) ([]*Status, error)
	UpdateStatus(context.Context, *Status) (error, Code)
	DeleteStatus(context.Context, int) error
}

type service struct {
	db          pg.DB
	userService schema.UsersClient
}

func NewService(db pg.DB, userService schema.UsersClient) Service {
	return &service{
		db:          db,
		userService: userService,
	}
}
