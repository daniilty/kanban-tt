package core

import (
	"context"

	"github.com/daniilty/kanban-tt/tasks/internal/pg"
)

type Service interface {
	AddTask(context.Context, *Task) (error, bool)
	GetTasks(context.Context, string) ([]*Task, error)
	UpdateTask(context.Context, *Task) (error, bool)
	DeleteTask(context.Context, string) error

	AddStatus(context.Context, *Status) error
	GetStatuses(context.Context, string) ([]*Status, error)
	UpdateStatus(context.Context, *Status) (error, bool)
	DeleteStatus(context.Context, string) error
}

type service struct {
	db pg.DB
}

func NewService(db pg.DB) Service {
	return &service{
		db: db,
	}
}
